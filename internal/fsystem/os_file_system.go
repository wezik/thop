package fsystem

import "os"

type FileSystem interface {
	MkdirAll(path string) error
	ReadDir(path string) ([]os.DirEntry, error)
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
	RemoveAll(path string) error
}

type OsFileSystem struct{}

func NewOsFileSystem() FileSystem {
	return &OsFileSystem{}
}

func (s *OsFileSystem) MkdirAll(path string) error {
	return os.MkdirAll(path, 0755)
}

func (s *OsFileSystem) ReadDir(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}

func (s *OsFileSystem) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (s *OsFileSystem) WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func (s *OsFileSystem) RemoveAll(path string) error {
	return os.RemoveAll(path)
}
