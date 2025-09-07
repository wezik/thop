package platform

import (
	"os"
	"os/exec"
)

// Platform abstraction for the system interactions
type Platform interface {
	Exit(int)
	OpenFile(string, int, os.FileMode) (*os.File, error)
	Exec(*exec.Cmd) (string, int, error)
	MkdirAll(string) error
	ReadDir(string) ([]os.DirEntry, error)
	WriteFile(string, []byte) error
	ReadFile(string) ([]byte, error)
}

type SystemPlatform struct {}

func NewSystemPlatform() Platform {
	return &SystemPlatform{}
}

func (s *SystemPlatform) Exit(code int) {
	os.Exit(code)
}

func (s *SystemPlatform) OpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(path, flag, perm)
}

func (s *SystemPlatform) Exec(cmd *exec.Cmd) (string, int, error) {
	res, err := cmd.Output()
	return string(res), cmd.ProcessState.ExitCode(), err
}

func (s *SystemPlatform) MkdirAll(path string) error {
	return os.MkdirAll(path, 0755)
}

func (s *SystemPlatform) ReadDir(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}

func (s *SystemPlatform) WriteFile(path string, bytes []byte) error {
	return os.WriteFile(path, bytes, 0644)
}

func (s *SystemPlatform) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
