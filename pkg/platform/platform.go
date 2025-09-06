package platform

import (
	"os"
	"os/exec"
)

type ExitFn func(int)
type GetwdFn func() (string, error)
type OpenFileFn func(string, int, os.FileMode) (*os.File, error)
type ExecFn func(*exec.Cmd) (string, int, error)
type MkdirAllFn func(string) error
type ReadDirFn func(string) ([]os.DirEntry, error)
type WriteFileFn func(string, []byte) error
