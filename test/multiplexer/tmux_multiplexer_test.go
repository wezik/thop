package multiplexer_test

import (
	"errors"
	"testing"
	"thop/internal/multiplexer"
	"thop/internal/types/command"
	"thop/internal/types/pane"
	"thop/internal/types/project"
	"thop/internal/types/template"
	"thop/internal/types/window"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTmuxClient struct {
	mock.Mock
}

func (m *MockTmuxClient) AttachSession(session multiplexer.SessionName) error {
	args := m.Called(session)
	return args.Error(0)
}

func (m *MockTmuxClient) SwitchSession(session multiplexer.SessionName) error {
	args := m.Called(session)
	return args.Error(0)
}

func (m *MockTmuxClient) HasSession(session multiplexer.SessionName) (bool, error) {
	args := m.Called(session)
	return args.Bool(0), args.Error(1)
}

func (m *MockTmuxClient) NewSession(
	session multiplexer.SessionName,
	root template.Root,
	win window.Window,
) error {
	args := m.Called(session, root, win)
	return args.Error(0)
}

func (m *MockTmuxClient) NewWindow(
	session multiplexer.SessionName,
	root template.Root,
	win window.Window,
) error {
	args := m.Called(session, root, win)
	return args.Error(0)
}

func (m *MockTmuxClient) SendKeys(
	paneID multiplexer.PaneID,
	keys command.Command,
) error {
	args := m.Called(paneID, keys)
	return args.Error(0)
}

func (m *MockTmuxClient) ListSessions() ([]multiplexer.SessionName, error) {
	args := m.Called()
	return args.Get(0).([]multiplexer.SessionName), args.Error(1)
}

func (m *MockTmuxClient) IsTmuxServerRunning() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockTmuxClient) KillSession(session multiplexer.SessionName) error {
	args := m.Called(session)
	return args.Error(0)
}

func (m *MockTmuxClient) NewPane(session multiplexer.SessionName, win window.Name, p pane.Pane, fallbackRoot string) error {
	args := m.Called(session, win, p, fallbackRoot)
	return args.Error(0)
}

func (m *MockTmuxClient) ListPanes(session multiplexer.SessionName, windowName window.Name) ([]multiplexer.PaneID, error) {
	args := m.Called(session, windowName)
	return args.Get(0).([]multiplexer.PaneID), args.Error(1)
}

func (m *MockTmuxClient) SetLayout(session multiplexer.SessionName, window window.Window) error {
	args := m.Called(session, window)
	return args.Error(0)
}

func Test_AttachProject(t *testing.T) {
	t.Run("assembles and attaches to session if it doesn't exist", func(t *testing.T) {
		// given
		sessionName := multiplexer.SessionName("foo")
		root := template.Root("/home/test")
		window1Name := window.Name("main")
		window2Name := window.Name("baz")
		project := project.Project{
			UUID: "uuid",
			Name: "foo",
			Template: template.Template{
				Root:     root,
				Commands: []command.Command{"echo hello"},
				Windows: []window.Window{
					{
						Name:   window1Name,
						Root:   "/project",
						Panes:  []pane.Pane{{}},
						Layout: window.LayoutTiled,
					},
					{
						Name:     window2Name,
						Commands: []command.Command{"ls"},
						Panes: []pane.Pane{
							{},
							{
								Commands: []command.Command{"echo pane"},
							},
						},
						Layout: window.LayoutMainVertical,
					},
				},
			},
		}

		mockClient := new(MockTmuxClient)
		mockClient.On("HasSession", sessionName).Return(false, nil).Once()
		mockClient.On("NewSession", sessionName, root, project.Template.Windows[0]).Return(nil).Once()
		mockClient.On("NewWindow", sessionName, root, project.Template.Windows[1]).Return(nil).Once()
		mockClient.On("NewPane", sessionName, window2Name, project.Template.Windows[1].Panes[1], string(root)).Return(nil).Once()

		mockClient.On("SetLayout", sessionName, project.Template.Windows[1]).Return(nil)

		mockClient.On("ListPanes", sessionName, window1Name).Return([]multiplexer.PaneID{"A"}, nil).Once()
		mockClient.On("ListPanes", sessionName, window2Name).Return([]multiplexer.PaneID{"B", "C"}, nil).Once()

		mockClient.On("SendKeys", multiplexer.PaneID("A"), command.Command("echo hello")).Return(nil).Once()

		mockClient.On("SendKeys", multiplexer.PaneID("B"), command.Command("echo hello")).Return(nil).Once()
		mockClient.On("SendKeys", multiplexer.PaneID("B"), command.Command("ls")).Return(nil).Once()

		mockClient.On("SendKeys", multiplexer.PaneID("C"), command.Command("echo hello")).Return(nil).Once()
		mockClient.On("SendKeys", multiplexer.PaneID("C"), command.Command("ls")).Return(nil).Once()
		mockClient.On("SendKeys", multiplexer.PaneID("C"), command.Command("echo pane")).Return(nil).Once()

		mockClient.On("AttachSession", sessionName).Return(nil).Once()

		multiplexer := multiplexer.TmuxMultiplexer{
			Client: mockClient,
		}

		// when
		err := multiplexer.AttachProject(project)

		// then
		assert.Nil(t, err, "Expected no error")
		mockClient.AssertExpectations(t)
	})

	t.Run("kills the session if assembling fails after session is created", func(t *testing.T) {
		// given
		sessionName := multiplexer.SessionName("foo")
		root := template.Root("/home/test")
		window1Name := window.Name("main")
		project := project.Project{
			UUID: "uuid",
			Name: "foo",
			Template: template.Template{
				Root:     root,
				Commands: []command.Command{"echo hello"},
				Windows: []window.Window{
					{
						Name:   window1Name,
						Root:   "/project",
						Panes:  []pane.Pane{{}},
						Layout: window.LayoutTiled,
					},
				},
			},
		}

		mockClient := new(MockTmuxClient)
		mockClient.On("HasSession", sessionName).Return(false, nil).Once()
		mockClient.On("NewSession", sessionName, root, project.Template.Windows[0]).Return(nil).Once()
		expectedErr := errors.New("some error")
		mockClient.On("ListPanes", sessionName, window1Name).Return([]multiplexer.PaneID(nil), expectedErr).Once()
		mockClient.On("KillSession", sessionName).Return(nil).Once()

		multiplexer := multiplexer.TmuxMultiplexer{
			Client: mockClient,
		}

		// when
		err := multiplexer.AttachProject(project)

		// then
		assert.Equal(t, expectedErr, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("attaches if session exists", func(t *testing.T) {
		// given
		project := project.Project{
			UUID: "foo",
			Name: "foo",
			Template: template.Template{
				Root: "/home/test",
			},
		}

		mockClient := new(MockTmuxClient)
		mockClient.On("HasSession", multiplexer.SessionName("foo")).Return(true, nil).Once()
		mockClient.On("AttachSession", multiplexer.SessionName("foo")).Return(nil).Once()

		multiplexer := multiplexer.TmuxMultiplexer{
			Client: mockClient,
		}

		// when
		err := multiplexer.AttachProject(project)

		// then
		assert.Nil(t, err, "Expected no error")
		mockClient.AssertExpectations(t)
	})

	t.Run("assembles and switches to session if it doesn't exist and shell is in active session", func(t *testing.T) {
		// given
		sessionName := multiplexer.SessionName("foo")
		root := template.Root("/home/test")
		window1Name := window.Name("main")
		window2Name := window.Name("baz")
		project := project.Project{
			UUID: "uuid",
			Name: "foo",
			Template: template.Template{
				Root:     root,
				Commands: []command.Command{"echo hello"},
				Windows: []window.Window{
					{
						Name:   window1Name,
						Root:   "/project",
						Panes:  []pane.Pane{{}},
						Layout: window.LayoutTiled,
					},
					{
						Name:     window2Name,
						Commands: []command.Command{"ls"},
						Panes: []pane.Pane{
							{},
							{
								Commands: []command.Command{"echo pane"},
							},
						},
						Layout: window.LayoutTiled,
					},
				},
			},
		}

		mockClient := new(MockTmuxClient)
		mockClient.On("HasSession", sessionName).Return(false, nil).Once()
		mockClient.On("NewSession", sessionName, root, project.Template.Windows[0]).Return(nil).Once()
		mockClient.On("NewWindow", sessionName, root, project.Template.Windows[1]).Return(nil).Once()
		mockClient.On("NewPane", sessionName, window2Name, project.Template.Windows[1].Panes[1], string(root)).Return(nil).Once()

		mockClient.On("SetLayout", sessionName, project.Template.Windows[1]).Return(nil)

		mockClient.On("ListPanes", sessionName, window1Name).Return([]multiplexer.PaneID{"A"}, nil).Once()
		mockClient.On("ListPanes", sessionName, window2Name).Return([]multiplexer.PaneID{"B", "C"}, nil).Once()

		mockClient.On("SendKeys", multiplexer.PaneID("A"), command.Command("echo hello")).Return(nil).Once()

		mockClient.On("SendKeys", multiplexer.PaneID("B"), command.Command("echo hello")).Return(nil).Once()
		mockClient.On("SendKeys", multiplexer.PaneID("B"), command.Command("ls")).Return(nil).Once()

		mockClient.On("SendKeys", multiplexer.PaneID("C"), command.Command("echo hello")).Return(nil).Once()
		mockClient.On("SendKeys", multiplexer.PaneID("C"), command.Command("ls")).Return(nil).Once()
		mockClient.On("SendKeys", multiplexer.PaneID("C"), command.Command("echo pane")).Return(nil).Once()

		mockClient.On("SwitchSession", sessionName).Return(nil).Once()

		multiplexer := multiplexer.TmuxMultiplexer{
			Client:            mockClient,
			ActiveTmuxSession: "/home/test",
		}

		// when
		err := multiplexer.AttachProject(project)

		// then
		assert.Nil(t, err, "Expected no error")
		mockClient.AssertExpectations(t)
	})

	t.Run("switches to session if it exist and shell is in active session", func(t *testing.T) {
		// given
		project := project.Project{
			UUID: "foo",
			Name: "foo",
			Template: template.Template{
				Root: "/home/test",
			},
		}

		mockClient := new(MockTmuxClient)
		mockClient.On("HasSession", multiplexer.SessionName("foo")).Return(true, nil).Once()
		mockClient.On("SwitchSession", multiplexer.SessionName("foo")).Return(nil).Once()

		multiplexer := multiplexer.TmuxMultiplexer{
			Client:            mockClient,
			ActiveTmuxSession: "/home/test",
		}

		// when
		err := multiplexer.AttachProject(project)

		// then
		assert.Nil(t, err, "Expected no error")
		mockClient.AssertExpectations(t)
	})

	t.Run("uses template name for session if provided", func(t *testing.T) {
		// given
		project := project.Project{
			UUID: "foo",
			Name: "foo",
			Template: template.Template{
				Name: "bar",
				Root: "/home/test",
			},
		}

		mockClient := new(MockTmuxClient)
		mockClient.On("HasSession", multiplexer.SessionName("bar")).Return(true, nil).Once()
		mockClient.On("AttachSession", multiplexer.SessionName("bar")).Return(nil).Once()

		multiplexer := multiplexer.TmuxMultiplexer{
			Client: mockClient,
		}

		// when
		err := multiplexer.AttachProject(project)

		// then
		assert.Nil(t, err, "Expected no error")
		mockClient.AssertExpectations(t)
	})

	t.Run("returns error if project has no name", func(t *testing.T) {
		multiplexer := multiplexer.TmuxMultiplexer{
			Client: nil,
		}

		err := multiplexer.AttachProject(project.Project{Name: "", Template: template.Template{Name: ""}})
		assert.NotNil(t, err, "Expected error when project has no name")
	})
}

func Test_ListActiveSessions(t *testing.T) {
	t.Run("returns empty list if tmux server is not running", func(t *testing.T) {
		// given
		mockClient := new(MockTmuxClient)
		mockClient.On("IsTmuxServerRunning").Return(false).Once()

		m := multiplexer.TmuxMultiplexer{
			Client: mockClient,
		}

		// when
		activeSessions, err := m.ListActiveSessions()

		// then
		assert.Nil(t, err)
		assert.Equal(t, []project.Project(nil), activeSessions)
		mockClient.AssertExpectations(t)
	})

	t.Run("returns empty list if client returns empty list", func(t *testing.T) {
		// given
		mockClient := new(MockTmuxClient)
		mockClient.On("ListSessions").Return([]multiplexer.SessionName{}, nil).Once()
		mockClient.On("IsTmuxServerRunning").Return(true).Once()

		multiplexer := multiplexer.TmuxMultiplexer{
			Client: mockClient,
		}

		// when
		sessions, err := multiplexer.ListActiveSessions()

		// then
		assert.Nil(t, err)
		assert.Equal(t, []project.Project(nil), sessions)
		mockClient.AssertExpectations(t)
	})

	t.Run("returns list of session type projects", func(t *testing.T) {
		// given
		mockClient := new(MockTmuxClient)
		mockClient.On("ListSessions").Return([]multiplexer.SessionName{"foo", "bar"}, nil).Once()
		mockClient.On("IsTmuxServerRunning").Return(true).Once()

		multiplexer := multiplexer.TmuxMultiplexer{
			Client: mockClient,
		}

		// when
		sessions, err := multiplexer.ListActiveSessions()

		// then
		assert.Nil(t, err)
		for i, session := range []project.Project{{Name: "foo", Type: project.TypeTmuxSession}, {Name: "bar", Type: project.TypeTmuxSession}} {
			assert.Equal(t, session, sessions[i])
		}
		mockClient.AssertExpectations(t)
	})
}
