package service

import (
	"os/exec"
	"thop/internal/config"
	"thop/internal/executor"
	"thop/internal/multiplexer"
	"thop/internal/problem"
	"thop/internal/selector"
	"thop/internal/storage"
	"thop/internal/types"
	"thop/internal/types/project"
	"thop/internal/types/template"
)

type Service interface {
	CreateProject(template.Root, project.Name) error
	OpenProject(project.Name) error
	DeleteProject(project.Name) error
	EditProject(project.Name) error
	KillSession(project.Name) error
}

type AppService struct {
	Selector    selector.ProjectSelector
	Multiplexer multiplexer.Multiplexer
	Storage     storage.Storage
	Config      *config.Config
	E           executor.CommandExecutor
}

const (
	ErrEditorNotSet             problem.Key = "THOP_EDITOR_NOT_SET"
	ErrEmptyProjectName         problem.Key = "THOP_EMPTY_PROJECT_NAME"
	ErrEmptyRootPath            problem.Key = "THOP_EMPTY_ROOT_PATH"
	ErrSelectedNonExisting      problem.Key = "THOP_SELECTED_NON_EXISTING"
	ErrSessionNotFound          problem.Key = "THOP_SESSION_NOT_FOUND"
	ErrProjectOrSessionNotFound problem.Key = "THOP_PROJECT_OR_SESSION_NOT_FOUND"
)

const (
	TemplateVersion = types.V1
)

func (s *AppService) CreateProject(root template.Root, name project.Name) error {
	if name == "" {
		return ErrEmptyProjectName.WithMsg("project name cannot be empty")
	}

	if root == "" {
		return ErrEmptyRootPath.WithMsg("root path cannot be empty")
	}

	template := template.Template{
		Root: root,
	}

	p := project.Project{
		Name:     name,
		Version:  TemplateVersion,
		Template: template.WithDefaults(),
	}

	return s.Storage.Save(&p)
}

func (s *AppService) OpenProject(name project.Name) error {
	if name != "" {
		p, err := s.Storage.Find(name)

		if err == nil {
			return s.Multiplexer.AttachProject(p)
		}

		if !storage.ErrProjectNotFound.Equal(err) {
			return err
		}

		// try to find active session if no template is found
		active, err := s.Multiplexer.ListActiveSessions()
		if err != nil {
			return err
		}

		for _, session := range active {
			if session.Name == name {
				return s.Multiplexer.AttachProject(session)
			}
		}

		return ErrProjectOrSessionNotFound.WithMsg(name)
	}

	projects, err := s.Storage.List()
	if err != nil {
		return err
	}

	sessions, err := s.Multiplexer.ListActiveSessions()
	if err != nil {
		return err
	}

	projects = append(projects, sessions...)

	found, err := s.Selector.SelectFrom(projects, "Select project to open > ")
	if err != nil {
		return err
	}

	return s.Multiplexer.AttachProject(*found)
}

func (s *AppService) DeleteProject(name project.Name) error {
	p, err := s.findOrSelect(name, "Select project to delete > ")
	if err != nil {
		return err
	}

	return s.Storage.Delete(p.UUID)
}

func (s *AppService) EditProject(name project.Name) error {
	p, err := s.findOrSelect(name, "Select project to edit > ")
	if err != nil {
		return err
	}

	templatePath, err := s.Storage.PrepareTemplateFile(p)
	if err != nil {
		return err
	}

	editor := s.Config.GetEditor()
	if editor == "" {
		return ErrEditorNotSet.WithMsg("$EDITOR environment variable is not set")
	}

	cmd := exec.Command(editor, templatePath)
	_, err = s.E.ExecuteInteractive(cmd)
	return err
}

func (s *AppService) KillSession(name project.Name) error {
	sessions, err := s.Multiplexer.ListActiveSessions()
	if err != nil {
		return err
	}

	if name != "" {
		for _, session := range sessions {
			if session.Name == name {
				return s.Multiplexer.KillSession(session)
			}
		}
		return ErrSessionNotFound.WithMsg(name)
	}

	selected, err := s.Selector.SelectFrom(sessions, "Select session to kill > ")
	if err != nil {
		return err
	}

	return s.Multiplexer.KillSession(*selected)
}

// common logic used by most commands
func (s *AppService) findOrSelect(name project.Name, prompt string) (project.Project, error) {
	if name != "" {
		return s.Storage.Find(name)
	}

	projects, err := s.Storage.List()
	if err != nil {
		return project.Project{}, err
	}

	selected, err := s.Selector.SelectFrom(projects, prompt)

	if err != nil {
		return project.Project{}, err
	}

	return *selected, nil
}
