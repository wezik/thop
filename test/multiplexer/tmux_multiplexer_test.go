package multiplexer_test

import (
	"errors"
	"testing"
	"thop/internal/logger"
	"thop/internal/multiplexer"
	"thop/internal/types/command"
	"thop/internal/types/pane"
	"thop/internal/types/project"
	"thop/internal/types/template"
	"thop/internal/types/window"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	logger.Init("")
}

type MockTmuxClient struct {
	mock.Mock
}

func (m *MockTmuxClient) AttachSession(sn multiplexer.SessionName) error {
	args := m.Called(sn)
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
	sn multiplexer.SessionName,
	t template.Template,
) (multiplexer.WindowID, multiplexer.PaneID, error) {
	args := m.Called(sn, t)
	return args.Get(0).(multiplexer.WindowID), args.Get(1).(multiplexer.PaneID), args.Error(2)
}

func (m *MockTmuxClient) NewWindow(
	session multiplexer.SessionName,
	root template.Root,
	win window.Window,
) (multiplexer.WindowID, multiplexer.PaneID, error) {
	args := m.Called(session, root, win)
	return args.Get(0).(multiplexer.WindowID), args.Get(1).(multiplexer.PaneID), args.Error(2)
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

func (m *MockTmuxClient) NewPane(
	wID multiplexer.WindowID,
	tr template.Root,
	wr window.Root,
	p pane.Pane,
) (multiplexer.PaneID, error) {
	args := m.Called(wID, tr, wr, p)
	return args.Get(0).(multiplexer.PaneID), args.Error(1)
}

func (m *MockTmuxClient) SetLayout(wID multiplexer.WindowID, wl window.Layout) error {
	args := m.Called(wID, wl)
	return args.Error(0)
}

func (m *MockTmuxClient) SetActivePane(pID multiplexer.PaneID) error {
	args := m.Called(pID)
	return args.Error(0)
}

func (m *MockTmuxClient) SetActiveWindow(wID multiplexer.WindowID) error {
	args := m.Called(wID)
	return args.Error(0)
}

func Test_AttachProject(t *testing.T) {
	t.Run("assembles and attaches to session if it doesn't exist", func(t *testing.T) {
		// given
		sn := multiplexer.SessionName("foo")
		root := template.Root("/home/test")
		w1ID := multiplexer.WindowID("A")
		w2ID := multiplexer.WindowID("B")
		p1ID := multiplexer.PaneID("C")
		p2ID := multiplexer.PaneID("D")
		p3ID := multiplexer.PaneID("E")
		p := project.Project{
			UUID: "uuid",
			Name: "foo",
			Template: template.Template{
				Root:     root,
				Commands: []command.Command{"echo hello"},
				Windows: []window.Window{
					{
						Root:   "/project",
						Panes:  []pane.Pane{{}},
						Layout: window.LayoutTiled,
						Active: true,
					},
					{
						Commands: []command.Command{"ls"},
						Panes: []pane.Pane{
							{
								Active: true,
							},
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
		mockClient.On("HasSession", sn).Return(false, nil).Once()
		mockClient.On("NewSession", sn, p.Template).Return(w1ID, p1ID, nil).Once()
		mockClient.On("NewWindow", sn, root, p.Template.Windows[1]).Return(w2ID, p2ID, nil).Once()
		mockClient.On("NewPane", w2ID, root, p.Template.Windows[1].Root, p.Template.Windows[1].Panes[1]).Return(p3ID, nil).Once()

		mockClient.On("SetLayout", w2ID, p.Template.Windows[1].Layout).Return(nil)

		mockClient.On("SendKeys", p1ID, command.Command("echo hello")).Return(nil).Once()

		mockClient.On("SendKeys", p2ID, command.Command("echo hello")).Return(nil).Once()
		mockClient.On("SendKeys", p2ID, command.Command("ls")).Return(nil).Once()

		mockClient.On("SendKeys", p3ID, command.Command("echo hello")).Return(nil).Once()
		mockClient.On("SendKeys", p3ID, command.Command("ls")).Return(nil).Once()
		mockClient.On("SendKeys", p3ID, command.Command("echo pane")).Return(nil).Once()

		mockClient.On("SetActiveWindow", w1ID).Return(nil).Once()
		mockClient.On("SetActivePane", p2ID).Return(nil).Once()

		mockClient.On("AttachSession", sn).Return(nil).Once()

		multiplexer := multiplexer.TmuxMultiplexer{
			Client: mockClient,
		}

		// when
		err := multiplexer.AttachProject(p)

		// then
		assert.Nil(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("kills the session if assembling fails after session is created", func(t *testing.T) {
		// given
		sessionName := multiplexer.SessionName("foo")
		root := template.Root("/home/test")
		p := project.Project{
			UUID: "uuid",
			Name: "foo",
			Template: template.Template{
				Root:     root,
				Commands: []command.Command{"echo hello"},
				Windows: []window.Window{
					{
						Root:   "/project",
						Panes:  []pane.Pane{{}},
						Layout: window.LayoutTiled,
					},
					{
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
		mockClient.On("NewSession", sessionName, p.Template).Return(multiplexer.WindowID(""), multiplexer.PaneID(""), nil).Once()
		expectedErr := errors.New("some error")
		mockClient.On("NewWindow", sessionName, root, p.Template.Windows[1]).Return(multiplexer.WindowID(""), multiplexer.PaneID(""), expectedErr).Once()
		mockClient.On("KillSession", sessionName).Return(nil).Once()

		multiplexer := multiplexer.TmuxMultiplexer{
			Client: mockClient,
		}

		// when
		err := multiplexer.AttachProject(p)

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
		sn := multiplexer.SessionName("foo")
		root := template.Root("/home/test")
		w1ID := multiplexer.WindowID("A")
		w2ID := multiplexer.WindowID("B")
		p1ID := multiplexer.PaneID("C")
		p2ID := multiplexer.PaneID("D")
		p3ID := multiplexer.PaneID("E")
		p := project.Project{
			UUID: "uuid",
			Name: "foo",
			Template: template.Template{
				Root:     root,
				Commands: []command.Command{"echo hello"},
				Windows: []window.Window{
					{
						Root:   "/project",
						Panes:  []pane.Pane{{}},
						Layout: window.LayoutTiled,
					},
					{
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
		mockClient.On("HasSession", sn).Return(false, nil).Once()
		mockClient.On("NewSession", sn, p.Template).Return(w1ID, p1ID, nil).Once()
		mockClient.On("NewWindow", sn, root, p.Template.Windows[1]).Return(w2ID, p2ID, nil).Once()
		mockClient.On("NewPane", w2ID, root, p.Template.Windows[1].Root, p.Template.Windows[1].Panes[1]).Return(p3ID, nil).Once()

		mockClient.On("SetLayout", w2ID, p.Template.Windows[1].Layout).Return(nil)

		mockClient.On("SendKeys", p1ID, command.Command("echo hello")).Return(nil).Once()

		mockClient.On("SendKeys", p2ID, command.Command("echo hello")).Return(nil).Once()
		mockClient.On("SendKeys", p2ID, command.Command("ls")).Return(nil).Once()

		mockClient.On("SendKeys", p3ID, command.Command("echo hello")).Return(nil).Once()
		mockClient.On("SendKeys", p3ID, command.Command("ls")).Return(nil).Once()
		mockClient.On("SendKeys", p3ID, command.Command("echo pane")).Return(nil).Once()

		mockClient.On("SwitchSession", sn).Return(nil).Once()

		multiplexer := multiplexer.TmuxMultiplexer{
			Client:            mockClient,
			ActiveTmuxSession: "/home/test",
		}

		// when
		err := multiplexer.AttachProject(p)

		// then
		assert.Nil(t, err)
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
