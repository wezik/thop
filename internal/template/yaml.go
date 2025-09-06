package template

import (
	"path"
	"path/filepath"
	"strings"
	"thop/pkg/log"
	"thop/pkg/platform"
	"thop/pkg/template"

	"github.com/goccy/go-yaml"
)

type YamlStorage struct {
	log       log.Logger
	config    *TemplateConfig
	mkdirAll  platform.MkdirAllFn
	readDir   platform.ReadDirFn
	writeFile platform.WriteFileFn
}

func NewYamlStorage(
	log log.Logger,
	config *TemplateConfig,
	mkdirAll platform.MkdirAllFn,
	readDir platform.ReadDirFn,
	writeFile platform.WriteFileFn,
) template.FileStorage {
	return &YamlStorage{
		log:       log,
		config:    config,
		mkdirAll:  mkdirAll,
		readDir:   readDir,
		writeFile: writeFile,
	}
}

// TODO: this implementation for now is replacing the file if it already exists
// this is not ideal, but don't care for now, will maybe split to "new" and "update" later
func (s *YamlStorage) Save(template *template.Template) (err error) {
	s.log.Debug("Saving template \"" + template.Name() + "\"")

	// ensure exists
	if err = s.mkdirAll(s.config.FileStoragePath); err != nil {
		return
	}

	filePath := filepath.Join(s.config.FileStoragePath, fileSafeName(template.Name()))

	// write yaml
	yamlTemplate := mapToYamlTemplate(template)
	bytes, err := yaml.Marshal(yamlTemplate)
	if err != nil {
		return
	}

	return s.writeFile(filePath, bytes)
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
	return nil
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
			templateFile := mapToTemplateFile(file, path)
			results = append(results, templateFile)
		}
	}

	return
}

func isYamlFile(name string) bool {
	return strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml")
}
