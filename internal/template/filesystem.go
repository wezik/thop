package template

import "thop/pkg/template"

type FileSystemStorage struct {
}

func NewFileSystemStorage() template.FileStorage {
	return &FileSystemStorage{}
}

func (f *FileSystemStorage) List() (files []*template.File, err error) {
	files = []*template.File{
		template.NewFile("file1", "/path/file1.yaml"),
		template.NewFile("file2", "/path/file2.txt"),
	}

	return
}
