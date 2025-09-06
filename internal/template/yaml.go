package template

import (
	"thop/pkg/log"
	"thop/pkg/template"
)

type YamlStorage struct {
	log log.Logger
}

func NewYamlStorage(
	log log.Logger,
) template.FileStorage {
	return &YamlStorage{
		log: log,
	}
}

func (y *YamlStorage) List() (files []*template.File, err error) {
	y.log.Debug("Listing template files")
	files = []*template.File{
		template.NewFile("file1", "/path/file1.yaml"),
		template.NewFile("file2", "/path/file2.txt"),
	}
	return
}

func (y *YamlStorage) LoadTemplate(path template.FilePath) *template.Template {
	y.log.Debug("Loading template from path \"" + string(path) + "\"")
	return &template.Template{
		Name: "/path/file1.yaml",
	}
}
