package template

import (
	"path"
	"strings"
	"thop/pkg/log"
	"thop/pkg/platform"
	"thop/pkg/template"
)

type YamlStorage struct {
	log      log.Logger
	config   *TemplateConfig
	mkdirAll platform.MkdirAllFn
	readDir  platform.ReadDirFn
}

func NewYamlStorage(
	log log.Logger,
	config *TemplateConfig,
	mkdirAll platform.MkdirAllFn,
	readDir platform.ReadDirFn,
) template.FileStorage {
	return &YamlStorage{
		log:      log,
		config:   config,
		mkdirAll: mkdirAll,
		readDir:  readDir,
	}
}

func (s *YamlStorage) List() (results []*template.File, err error) {
	s.log.Debug("Listing template files")

	// ensure exists
	if err = s.mkdirAll(s.config.FileStoragePath); err != nil {
		return
	}

	// collect all yaml files from the tree
	return s.listTemplatesFromRoot(s.config.FileStoragePath)
}

func (s *YamlStorage) LoadTemplate(path template.FilePath) *template.Template {
	s.log.Debug("Loading template from path \"" + string(path) + "\"")
	return &template.Template{
		Name: "/path/file1.yaml",
	}
}

func (s *YamlStorage) listTemplatesFromRoot(root string) (results []*template.File, err error) {
	files, err := s.readDir(root)
	for _, file := range files {
		if file.IsDir() {
			var nestedResults []*template.File
			nestedResults, err = s.listTemplatesFromRoot(path.Join(root, file.Name()))
			if err != nil {
				return
			}

			results = append(results, nestedResults...)
			continue
		}

		if isYamlFile(file.Name()) {
			path := path.Join(root, file.Name())
			templateFile := template.NewFile(normalizeName(file.Name()), path)
			results = append(results, templateFile)
		}
	}

	return
}

func isYamlFile(name string) bool {
	return strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml")
}

func normalizeName(fileName string) string {
	name := strings.TrimSuffix(fileName, ".yaml")
	name = strings.TrimSuffix(name, ".yml")
	name = strings.ReplaceAll(name, "_", " ")
	return name
}
