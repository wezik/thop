package storage

import (
	"fmt"
	"path/filepath"
	"thop/internal/config"
	"thop/internal/fsystem"
	"thop/internal/problem"
	"thop/internal/types/project"

	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
)

type Storage interface {
	List() ([]project.Project, error)
	Find(project.Name) (project.Project, error)
	Save(*project.Project) error
	Delete(uuid project.UUID) error
	PrepareTemplateFile(project.Project) (string, error)
}

type YamlStorage struct {
	Config     *config.Config
	FileSystem fsystem.FileSystem
}

const (
	ErrFailedToCreateTemplateDir problem.Key = "STORAGE_FAILED_TO_CREATE_TEMPLATE_DIR"
	ErrFailedToDeleteProject     problem.Key = "STORAGE_FAILED_TO_DELETE_PROJECT"
	ErrFailedToReadTemplateDir   problem.Key = "STORAGE_FAILED_TO_READ_TEMPLATE_DIR"
	ErrFailedToSaveProject       problem.Key = "STORAGE_FAILED_TO_SAVE_PROJECT"
	ErrFailedToSerializeProject  problem.Key = "STORAGE_FAILED_TO_SERIALIZE_PROJECT"
	ErrProjectNotFound           problem.Key = "STORAGE_PROJECT_NOT_FOUND"
)

const (
	templateFileName = "template.yaml"
	templatesDirName = "templates"
)

func (s *YamlStorage) List() ([]project.Project, error) {
	cfgDir := s.Config.GetConfigDir()

	templatesDir := filepath.Join(cfgDir, templatesDirName)

	// from what I understand, running os.Stat to check if a dir exists is not really providing
	// any benefits, and can also introduce weird edge cases, so instead just run mkdir everytime
	if err := s.FileSystem.MkdirAll(templatesDir); err != nil {
		return nil, ErrFailedToCreateTemplateDir.WithMsg(err.Error())
	}

	dirs, err := s.FileSystem.ReadDir(templatesDir)
	if err != nil {
		return nil, ErrFailedToReadTemplateDir.WithMsg(err.Error())
	}

	var projects []project.Project

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		dirName := dir.Name()

		templateFile := filepath.Join(templatesDir, dirName, templateFileName)
		bytes, err := s.FileSystem.ReadFile(templateFile)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var p project.Project
		if err = yaml.Unmarshal(bytes, &p); err != nil {
			fmt.Println(err)
			continue
		}

		p.UUID = project.UUID(dirName)
		p = p.WithDefaults()

		if p.IsValid() {
			projects = append(projects, p)
		}
	}

	return projects, nil
}

func (s *YamlStorage) Find(name project.Name) (project.Project, error) {
	projects, err := s.List()
	if err != nil {
		return project.Project{}, err
	}

	for _, p := range projects {
		if p.Name == name {
			p = p.WithDefaults()
			if !p.IsValid() {
				return project.Project{}, ErrProjectNotFound.WithMsg("project", name, "is invalid")
			}
			return p, nil
		}
	}

	return project.Project{}, ErrProjectNotFound.WithMsg("project", name, "not found")
}

func (s *YamlStorage) Save(p *project.Project) error {
	if p.UUID == "" {
		p.UUID = project.UUID(uuid.New().String())
	}

	cfgDir := s.Config.GetConfigDir()
	templateDir := filepath.Join(cfgDir, templatesDirName, string(p.UUID))

	if err := s.FileSystem.MkdirAll(templateDir); err != nil {
		return ErrFailedToCreateTemplateDir.WithMsg(err.Error())
	}

	templateFile := filepath.Join(templateDir, templateFileName)
	bytes, err := yaml.Marshal(p)
	if err != nil {
		return ErrFailedToSerializeProject.WithMsg(err.Error())
	}

	if err := s.FileSystem.WriteFile(templateFile, bytes); err != nil {
		return ErrFailedToSaveProject.WithMsg(err.Error())
	}

	return nil
}

func (s *YamlStorage) Delete(uuid project.UUID) error {
	cfgDir := s.Config.GetConfigDir()
	templateDir := filepath.Join(cfgDir, templatesDirName, string(uuid))
	if err := s.FileSystem.RemoveAll(templateDir); err != nil {
		return ErrFailedToDeleteProject.WithMsg(err.Error())
	}

	return nil
}

func (s *YamlStorage) PrepareTemplateFile(p project.Project) (string, error) {
	cfgDir := s.Config.GetConfigDir()
	return filepath.Join(cfgDir, templatesDirName, string(p.UUID), templateFileName), nil
}
