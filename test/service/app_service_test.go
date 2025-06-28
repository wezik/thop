package service_test

import (
	"errors"
	"fmt"
	"testing"
	"thop/internal/config"
	"thop/internal/problem"
	"thop/internal/service"
	"thop/internal/storage"
	"thop/internal/types/project"
	"thop/internal/types/template"
	"thop/internal/types/window"
	"thop/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CreateProject(t *testing.T) {
	t.Run("creates project", func(t *testing.T) {
		// given
		stMock := new(test.MockStorage)
		svc := &service.AppService{
			Storage: stMock,
		}

		cwd := template.Root("/home/test")
		name := project.Name("foobar")

		stMock.On("Save", mock.Anything).Return(nil)

		// when
		err := svc.CreateProject(cwd, name)

		// then
		assert.Nil(t, err)
	})

	t.Run("errors with invalid data", func(t *testing.T) {
		type TestCase struct {
			cwd  template.Root
			name project.Name
			err  problem.Key
		}
		for _, args := range []TestCase{
			{"", "", service.ErrEmptyProjectName},
			{"", "foo", service.ErrEmptyRootPath},
			{"/foo/bar", "", service.ErrEmptyProjectName},
		} {
			t.Run(fmt.Sprintf("for %s and %s", args.cwd, args.name), func(t *testing.T) {
				// given
				svc := new(service.AppService)
				cwd := args.cwd
				name := args.name

				// when
				err := svc.CreateProject(cwd, name)

				// then
				assert.True(t, args.err.Equal(err))
			})
		}
	})
}

func Test_OpenProject(t *testing.T) {
	t.Run("runs selector and attaches to project when name is empty", func(t *testing.T) {
		// given
		projects := []project.Project{{UUID: "1234", Name: "foobar"}}

		slMock := new(test.MockProjectSelector)
		slMock.On("SelectFrom", projects, mock.Anything).Return(&projects[0], nil).Once()

		stMock := new(test.MockStorage)
		stMock.On("List").Return(projects, nil).Once()

		muMock := new(test.MockMultiplexer)
		muMock.On("AttachProject", projects[0]).Return(nil).Once()

		muMock.On("ListActiveSessions").Return([]project.Project(nil), nil).Once()

		svc := &service.AppService{
			Selector:    slMock,
			Multiplexer: muMock,
			Storage:     stMock,
			E:           nil,
		}

		// when
		err := svc.OpenProject("")

		// then
		assert.Nil(t, err)
		slMock.AssertExpectations(t)
		stMock.AssertExpectations(t)
		muMock.AssertExpectations(t)
	})

	t.Run("lists active sessions and allows to attach to them", func(t *testing.T) {
		// given
		projects := []project.Project{
			{UUID: "1234", Name: "foobar"},
		}
		sessions := []project.Project{
			{Name: "foobar", Type: project.TypeTmuxSession},
			{Name: "barfoo", Type: project.TypeTmuxSession},
		}
		// First "session" should be filtered out
		combined := append(projects, sessions...)

		slMock := new(test.MockProjectSelector)
		slMock.On("SelectFrom", combined, mock.Anything).Return(&combined[1], nil).Once()

		stMock := new(test.MockStorage)
		stMock.On("List").Return(projects, nil).Once()

		muMock := new(test.MockMultiplexer)
		muMock.On("AttachProject", combined[1]).Return(nil).Once()
		muMock.On("ListActiveSessions").Return(sessions, nil).Once()

		svc := &service.AppService{
			Selector:    slMock,
			Multiplexer: muMock,
			Storage:     stMock,
			E:           nil,
		}

		// when
		err := svc.OpenProject("")

		// then
		assert.Nil(t, err)
		slMock.AssertExpectations(t)
		stMock.AssertExpectations(t)
		muMock.AssertExpectations(t)
	})

	t.Run("tries to find project if name is provided", func(t *testing.T) {
		// given
		p := project.Project{UUID: "1234", Name: "foobar"}

		stMock := new(test.MockStorage)
		stMock.On("Find", p.Name).Return(p, nil).Once()

		muMock := new(test.MockMultiplexer)
		muMock.On("AttachProject", p).Return(nil).Once()

		svc := &service.AppService{
			Selector:    nil,
			Multiplexer: muMock,
			Storage:     stMock,
			E:           nil,
		}

		// when
		err := svc.OpenProject(p.Name)

		// then
		assert.Nil(t, err)
		stMock.AssertExpectations(t)
		muMock.AssertExpectations(t)
	})

	t.Run("propagates find errors", func(t *testing.T) {
		// given
		expected := errors.New("expected error")

		stMock := new(test.MockStorage)
		stMock.On("Find", project.Name("foobar")).Return(project.Project{}, expected).Once()

		svc := &service.AppService{
			Selector:    nil,
			Multiplexer: nil,
			Storage:     stMock,
			E:           nil,
		}

		// when
		err := svc.OpenProject("foobar")

		// then
		assert.Equal(t, expected, err)
		stMock.AssertExpectations(t)
	})

	t.Run("propagates list errors", func(t *testing.T) {
		// given
		expected := errors.New("expected error")

		slMock := new(test.MockProjectSelector)
		slMock.On("SelectFrom", mock.Anything, mock.Anything).Return("", nil).Once()

		stMock := new(test.MockStorage)
		stMock.On("List").Return([]project.Project{}, expected).Once()

		svc := &service.AppService{
			Selector:    slMock,
			Multiplexer: nil,
			Storage:     stMock,
			E:           nil,
		}

		// when
		err := svc.OpenProject("")

		// then
		assert.Equal(t, expected, err)
		stMock.AssertExpectations(t)
	})

	t.Run("propagates selector errors", func(t *testing.T) {
		// given
		expected := errors.New("expected error")
		listReturn := []project.Project{{UUID: "1234", Name: "foobar"}}

		slMock := new(test.MockProjectSelector)
		slMock.On("SelectFrom", mock.Anything, mock.Anything).Return(&project.Project{}, expected).Once()

		stMock := new(test.MockStorage)
		stMock.On("List").Return(listReturn, nil).Once()

		muMock := new(test.MockMultiplexer)
		muMock.On("ListActiveSessions").Return([]project.Project(nil), nil).Once()

		svc := &service.AppService{
			Selector:    slMock,
			Multiplexer: muMock,
			Storage:     stMock,
			E:           nil,
		}

		// when
		err := svc.OpenProject("")

		// then
		assert.Equal(t, expected, err)
		slMock.AssertExpectations(t)
		stMock.AssertExpectations(t)
	})

	t.Run("finds active session and opens it", func(t *testing.T) {
		// given
		projects := []project.Project{
			{UUID: "1234", Name: "foobar"},
		}
		sessions := []project.Project{
			{Name: "foobar", Type: project.TypeTmuxSession},
			{Name: "barfoo", Type: project.TypeTmuxSession},
		}
		combined := append(projects, sessions...)

		stMock := new(test.MockStorage)
		errNotFound := storage.ErrProjectNotFound.WithMsg("project", "foobar", "not found")
		stMock.On("Find", project.Name("foobar")).Return(project.Project{}, errNotFound).Once()

		muMock := new(test.MockMultiplexer)
		muMock.On("AttachProject", combined[1]).Return(nil).Once()
		muMock.On("ListActiveSessions").Return(sessions, nil).Once()

		svc := &service.AppService{
			Selector:    nil,
			Multiplexer: muMock,
			Storage:     stMock,
			E:           nil,
		}

		// when
		err := svc.OpenProject("foobar")

		// then
		assert.Nil(t, err)
		stMock.AssertExpectations(t)
		muMock.AssertExpectations(t)
	})
}

func Test_DeleteProject(t *testing.T) {
	t.Run("runs selector and deletes project when name is empty", func(t *testing.T) {
		// given
		projects := []project.Project{{UUID: "1234", Name: "foobar"}}

		slMock := new(test.MockProjectSelector)
		slMock.On("SelectFrom", projects, mock.Anything).Return(&projects[0], nil).Once()

		stMock := new(test.MockStorage)
		stMock.On("List").Return(projects, nil).Once()
		stMock.On("Delete", projects[0].UUID).Return(nil).Once()

		svc := &service.AppService{
			Selector:    slMock,
			Multiplexer: nil,
			Storage:     stMock,
			E:           nil,
		}

		// when
		err := svc.DeleteProject("")

		// then
		assert.Nil(t, err)
		slMock.AssertExpectations(t)
		stMock.AssertExpectations(t)
	})

	t.Run("tries to find project if name is provided", func(t *testing.T) {
		// given
		project := project.Project{UUID: "1234", Name: "foobar"}

		stMock := new(test.MockStorage)
		stMock.On("Find", project.Name).Return(project, nil).Once()
		stMock.On("Delete", project.UUID).Return(nil).Once()

		svc := &service.AppService{
			Selector:    nil,
			Multiplexer: nil,
			Storage:     stMock,
			E:           nil,
		}

		// when
		err := svc.DeleteProject(project.Name)

		// then
		assert.Nil(t, err)
		stMock.AssertExpectations(t)
	})

	t.Run("propagates find errors", func(t *testing.T) {
		// given
		expected := errors.New("expected error")

		stMock := new(test.MockStorage)
		stMock.On("Find", project.Name("foobar")).Return(project.Project{}, expected).Once()

		svc := &service.AppService{
			Selector:    nil,
			Multiplexer: nil,
			Storage:     stMock,
			E:           nil,
		}

		// when
		err := svc.DeleteProject("foobar")

		// then
		assert.Equal(t, expected, err)
		stMock.AssertExpectations(t)
	})

	t.Run("propagates list errors", func(t *testing.T) {
		// given
		expected := errors.New("expected error")

		slMock := new(test.MockProjectSelector)
		slMock.On("SelectFrom", mock.Anything, mock.Anything).Return("", nil).Once()

		stMock := new(test.MockStorage)
		stMock.On("List").Return([]project.Project{}, expected).Once()

		svc := &service.AppService{
			Selector:    slMock,
			Multiplexer: nil,
			Storage:     stMock,
			E:           nil,
		}

		// when
		err := svc.DeleteProject("")

		// then
		assert.Equal(t, expected, err)
		stMock.AssertExpectations(t)
	})

	t.Run("propagates selector errors", func(t *testing.T) {
		// given
		expected := errors.New("expected error")
		listReturn := []project.Project{{UUID: "1234", Name: "foobar"}}

		slMock := new(test.MockProjectSelector)
		slMock.On("SelectFrom", mock.Anything, mock.Anything).Return(&project.Project{}, expected).Once()

		stMock := new(test.MockStorage)
		stMock.On("List").Return(listReturn, nil).Once()

		svc := &service.AppService{
			Selector:    slMock,
			Multiplexer: nil,
			Storage:     stMock,
			E:           nil,
		}

		// when
		err := svc.DeleteProject("")

		// then
		assert.Equal(t, expected, err)
		slMock.AssertExpectations(t)
		stMock.AssertExpectations(t)
	})
}

func Test_EditProject(t *testing.T) {
	t.Run("runs selector and launches editor when name is empty", func(t *testing.T) {
		// given
		projects := []project.Project{{UUID: "1234", Name: "foobar"}}

		template := template.Template{
			Root:    "/home/test",
			Windows: []window.Window{{Name: "main", Root: "/project"}},
		}

		projects[0].Template = template

		templateFile := "/home/test/template.yaml"

		slMock := new(test.MockProjectSelector)
		slMock.On("SelectFrom", projects, mock.Anything).Return(&projects[0], nil).Once()

		stMock := new(test.MockStorage)
		stMock.On("ListWithInvalid").Return(projects, nil).Once()
		stMock.On("PrepareTemplateFile", projects[0].UUID).Return(templateFile, nil).Once()

		executorMock := new(test.MockExecutor)
		executorMock.On("ExecuteInteractive", mock.Anything).Return(0, nil).Once()

		svc := &service.AppService{
			Selector:    slMock,
			Multiplexer: nil,
			Storage:     stMock,
			Config:      &config.Config{Editor: "vim"},
			E:           executorMock,
		}

		// when
		err := svc.EditProject("")

		// then
		assert.Nil(t, err)
		slMock.AssertExpectations(t)
		stMock.AssertExpectations(t)
		executorMock.AssertExpectations(t)
	})

	t.Run("tries to find project if name is provided", func(t *testing.T) {
		// given
		project := project.Project{UUID: "1234", Name: "foobar"}

		template := template.Template{
			Root:    "/home/test",
			Windows: []window.Window{{Name: "main", Root: "/project"}},
		}

		project.Template = template

		templateFile := "/home/test/template.yaml"

		stMock := new(test.MockStorage)
		stMock.On("Find", project.Name).Return(project, nil).Once()
		stMock.On("PrepareTemplateFile", project.UUID).Return(templateFile, nil).Once()

		editorMock := new(test.MockExecutor)
		editorMock.On("ExecuteInteractive", mock.Anything).Return(0, nil).Once()

		svc := &service.AppService{
			Selector:    nil,
			Multiplexer: nil,
			Storage:     stMock,
			Config:      &config.Config{Editor: "vim"},
			E:           editorMock,
		}

		// when
		err := svc.EditProject(project.Name)

		// then
		assert.Nil(t, err)
		stMock.AssertExpectations(t)
		editorMock.AssertExpectations(t)
	})

	t.Run("propagates find errors", func(t *testing.T) {
		// given
		expected := errors.New("expected error")

		stMock := new(test.MockStorage)
		stMock.On("Find", project.Name("foobar")).Return(project.Project{}, expected).Once()

		svc := &service.AppService{
			Selector:    nil,
			Multiplexer: nil,
			Storage:     stMock,
			E:           nil,
		}

		// when
		err := svc.EditProject("foobar")

		// then
		assert.Equal(t, expected, err)
		stMock.AssertExpectations(t)
	})

	t.Run("propagates list errors", func(t *testing.T) {
		// given
		expected := errors.New("expected error")

		stMock := new(test.MockStorage)
		stMock.On("ListWithInvalid").Return([]project.Project{}, expected).Once()

		svc := &service.AppService{
			Selector:    nil,
			Multiplexer: nil,
			Storage:     stMock,
			E:           nil,
		}

		// when
		err := svc.EditProject("")

		// then
		assert.Equal(t, expected, err)
		stMock.AssertExpectations(t)
	})

	t.Run("propagates selector errors", func(t *testing.T) {
		// given
		expected := errors.New("expected error")
		listReturn := []project.Project{{UUID: "1234", Name: "foobar"}}

		slMock := new(test.MockProjectSelector)
		slMock.On("SelectFrom", mock.Anything, mock.Anything).Return(&project.Project{}, expected).Once()

		stMock := new(test.MockStorage)
		stMock.On("ListWithInvalid").Return(listReturn, nil).Once()

		svc := &service.AppService{
			Selector:    slMock,
			Multiplexer: nil,
			Storage:     stMock,
			E:           nil,
		}

		// when
		err := svc.EditProject("")

		// then
		assert.Equal(t, expected, err)
		slMock.AssertExpectations(t)
		stMock.AssertExpectations(t)
	})
}

func Test_KillSession(t *testing.T) {
	t.Run("runs selector and kills session when name is empty", func(t *testing.T) {
		// given
		projects := []project.Project{{UUID: "1234", Name: "foobar", Type: project.TypeTmuxSession}}

		slMock := new(test.MockProjectSelector)
		slMock.On("SelectFrom", projects, mock.Anything).Return(&projects[0], nil).Once()

		muMock := new(test.MockMultiplexer)
		muMock.On("ListActiveSessions").Return(projects, nil).Once()
		muMock.On("KillSession", projects[0]).Return(nil).Once()

		svc := &service.AppService{
			Selector:    slMock,
			Multiplexer: muMock,
			Storage:     nil,
			E:           nil,
		}

		// when
		err := svc.KillSession("")

		// then
		assert.Nil(t, err)
		slMock.AssertExpectations(t)
		muMock.AssertExpectations(t)
	})

	t.Run("tries to find active session if name is provided", func(t *testing.T) {
		// given
		projects := []project.Project{
			{UUID: "1234", Name: "foobar", Type: project.TypeTmuxSession},
			{UUID: "5678", Name: "barfoo", Type: project.TypeTmuxSession},
		}

		stMock := new(test.MockStorage)

		muMock := new(test.MockMultiplexer)
		muMock.On("ListActiveSessions").Return(projects, nil).Once()
		muMock.On("KillSession", projects[1]).Return(nil).Once()

		svc := &service.AppService{
			Selector:    nil,
			Multiplexer: muMock,
			Storage:     nil,
			E:           nil,
		}

		// when
		err := svc.KillSession("barfoo")

		// then
		assert.Nil(t, err)
		stMock.AssertExpectations(t)
		muMock.AssertExpectations(t)
	})
}
