package multiplexer_test

import (
	"errors"
	"os/exec"
	"testing"
	"thop/internal/multiplexer"
	"thop/internal/types/pane"
	"thop/internal/types/template"
	"thop/internal/types/window"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCommandExecutor struct {
	mock.Mock
	ExecutedCommands [][]string
}

func (m *MockCommandExecutor) Execute(cmd *exec.Cmd) (string, int, error) {
	args := m.Called(cmd)
	m.ExecutedCommands = append(m.ExecutedCommands, cmd.Args)
	return args.String(0), args.Int(1), args.Error(2)
}

func (m *MockCommandExecutor) ExecuteInteractive(cmd *exec.Cmd) (int, error) {
	args := m.Called(cmd)
	m.ExecutedCommands = append(m.ExecutedCommands, cmd.Args)
	return args.Int(0), args.Error(1)
}

func Test_Client_AttachSession(t *testing.T) {
	t.Run("attaches to session", func(t *testing.T) {
		// given
		mockExecutor := new(MockCommandExecutor)

		mockExecutor.On("Execute", mock.Anything).Return("", 0, nil)

		expectedCmd := [][]string{
			{"tmux", "attach", "-t", "mysession"},
		}

		client := multiplexer.TmuxClientImpl{
			E: mockExecutor,
		}

		// when
		err := client.AttachSession("mysession")

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedCmd, mockExecutor.ExecutedCommands)
	})
}

func Test_Client_SwitchSession(t *testing.T) {
	t.Run("switches session", func(t *testing.T) {
		// given
		mockExecutor := new(MockCommandExecutor)
		mockExecutor.On("Execute", mock.Anything).Return("", 0, nil)

		expectedCmd := [][]string{
			{"tmux", "switch", "-t", "mysession"},
		}

		client := multiplexer.TmuxClientImpl{
			E: mockExecutor,
		}

		// when
		err := client.SwitchSession("mysession")

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedCmd, mockExecutor.ExecutedCommands)
	})

}

func Test_Client_HasSession(t *testing.T) {
	t.Run("returns false on exit code 1 gracefully", func(t *testing.T) {
		// given
		mockExecutor := new(MockCommandExecutor)
		mockExecutor.On("Execute", mock.Anything).Return("", 1, errors.New("exit code 1"))
		expectedCmd := [][]string{
			{"tmux", "has-session", "-t", "mysession"},
		}

		client := multiplexer.TmuxClientImpl{
			E: mockExecutor,
		}

		// when
		exists, err := client.HasSession("mysession")

		// then
		assert.Nil(t, err)
		assert.False(t, exists)
		assert.Equal(t, expectedCmd, mockExecutor.ExecutedCommands)
	})

	t.Run("returns true if session exists", func(t *testing.T) {
		// given
		mockExecutor := new(MockCommandExecutor)
		mockExecutor.On("Execute", mock.Anything).Return("", 0, nil)
		expectedCmd := [][]string{
			{"tmux", "has-session", "-t", "mysession"},
		}

		client := multiplexer.TmuxClientImpl{
			E: mockExecutor,
		}

		// when
		exists, err := client.HasSession("mysession")

		// then
		assert.Nil(t, err)
		assert.True(t, exists)
		assert.Equal(t, expectedCmd, mockExecutor.ExecutedCommands)
	})
}

func Test_Client_NewSession(t *testing.T) {
	t.Run("returns error if session already exists", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 1, errors.New("exit code 1"))
		expectedCmd := [][]string{
			{
				"tmux",
				"new-session",
				"-d",
				"-s",
				"mysession",
				"-c",
				"/home/test",
				"-n",
				"main",
				"-P", "-F", "#{window_id} #{pane_id}",
				"cd /project && exec $SHELL",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		tmpl := template.Template{
			Root: "/home/test",
			Windows: []window.Window{
				{
					Name:  "main",
					Root:  "/project",
					Panes: []pane.Pane{{}},
				},
			},
		}
		_, _, err := client.NewSession("mysession", tmpl)

		// then
		assert.NotNil(t, err)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})

	t.Run("creates new session", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("@1 %2\n", 0, nil)
		expectedCmd := [][]string{
			{
				"tmux",
				"new-session",
				"-d",
				"-s",
				"mysession",
				"-c",
				"/home/test",
				"-n",
				"main",
				"-P", "-F", "#{window_id} #{pane_id}",
				"cd /project && exec $SHELL",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		tmpl := template.Template{
			Root: "/home/test",
			Windows: []window.Window{
				{
					Name:  "main",
					Root:  "/project",
					Panes: []pane.Pane{{}},
				},
			},
		}
		wID, pID, err := client.NewSession("mysession", tmpl)

		// then
		assert.Nil(t, err)
		assert.Equal(t, multiplexer.WindowID("@1"), wID)
		assert.Equal(t, multiplexer.PaneID("%2"), pID)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})

	t.Run("defaults main window root to session root if empty", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("@1 %2\n", 0, nil)
		expectedCmd := [][]string{
			{
				"tmux",
				"new-session",
				"-d",
				"-s",
				"mysession",
				"-c",
				"/home/test",
				"-n",
				"main",
				"-P", "-F", "#{window_id} #{pane_id}",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		tmpl := template.Template{
			Root: "/home/test",
			Windows: []window.Window{
				{
					Name:  "main",
					Root:  "",
					Panes: []pane.Pane{{}},
				},
			},
		}
		_, _, err := client.NewSession("mysession", tmpl)

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})

	t.Run("creates new session with pane root", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("@1 %2\n", 0, nil)
		expectedCmd := [][]string{
			{
				"tmux",
				"new-session",
				"-d",
				"-s",
				"mysession",
				"-c",
				"/home/test",
				"-n",
				"main",
				"-P", "-F", "#{window_id} #{pane_id}",
				"cd /project/pane && exec $SHELL",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		tmpl := template.Template{
			Root: "/home/test",
			Windows: []window.Window{
				{
					Name: "main",
					Root: "/project",
					Panes: []pane.Pane{
						{
							Root: "/project/pane",
						},
					},
				},
			},
		}
		wID, pID, err := client.NewSession("mysession", tmpl)

		// then
		assert.Nil(t, err)
		assert.Equal(t, multiplexer.WindowID("@1"), wID)
		assert.Equal(t, multiplexer.PaneID("%2"), pID)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})

	t.Run("skips -c flag if all roots are empty", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("@1 %2\n", 0, nil)
		expectedCmd := [][]string{
			{
				"tmux",
				"new-session",
				"-d",
				"-s",
				"mysession",
				"-P", "-F", "#{window_id} #{pane_id}",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		tmpl := template.Template{
			Windows: []window.Window{
				{
					Panes: []pane.Pane{{}},
				},
			},
		}
		wID, pID, err := client.NewSession("mysession", tmpl)

		// then
		assert.Nil(t, err)
		assert.Equal(t, multiplexer.WindowID("@1"), wID)
		assert.Equal(t, multiplexer.PaneID("%2"), pID)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})
}

func Test_TmuxClient_SendKeys(t *testing.T) {
	t.Run("sends keys to pane", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 0, nil)
		expectedCmd := [][]string{
			{
				"tmux",
				"send-keys",
				"-t",
				"%id",
				"ls",
				"C-m",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		err := client.SendKeys("%id", "ls")

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})
}

func Test_Client_NewWindow(t *testing.T) {
	t.Run("includes default seession root if window root is empty", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("@1 %2\n", 0, nil)
		expectedCmd := [][]string{
			{
				"tmux",
				"new-window",
				"-d",
				"-t",
				"mysession",
				"-n",
				"main",
				"-c",
				"/home/test",
				"-P", "-F", "#{window_id} #{pane_id}",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		win := window.Window{Name: "main", Root: "", Panes: []pane.Pane{{}}}
		wID, pID, err := client.NewWindow("mysession", "/home/test", win)

		// then
		assert.Nil(t, err)
		assert.Equal(t, multiplexer.WindowID("@1"), wID)
		assert.Equal(t, multiplexer.PaneID("%2"), pID)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})

	t.Run("creates new window", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("@1 %2\n", 0, nil)
		expectedCmd := [][]string{
			{
				"tmux",
				"new-window",
				"-d",
				"-t",
				"mysession",
				"-n",
				"main",
				"-c",
				"/project",
				"-P", "-F", "#{window_id} #{pane_id}",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		win := window.Window{Name: "main", Root: "/project", Panes: []pane.Pane{{}}}
		wID, pID, err := client.NewWindow("mysession", "/home/test", win)

		// then
		assert.Nil(t, err)
		assert.Equal(t, multiplexer.WindowID("@1"), wID)
		assert.Equal(t, multiplexer.PaneID("%2"), pID)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})

	t.Run("creates new window with pane root", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("@1 %2\n", 0, nil)
		expectedCmd := [][]string{
			{
				"tmux",
				"new-window",
				"-d",
				"-t",
				"mysession",
				"-n",
				"main",
				"-c",
				"/project",
				"-P", "-F", "#{window_id} #{pane_id}",
				"cd /project/pane && exec $SHELL",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		win := window.Window{
			Name: "main",
			Root: "/project",
			Panes: []pane.Pane{
				{
					Root: "/project/pane",
				},
			},
		}
		wID, pID, err := client.NewWindow("mysession", "/home/test", win)

		// then
		assert.Nil(t, err)
		assert.Equal(t, multiplexer.WindowID("@1"), wID)
		assert.Equal(t, multiplexer.PaneID("%2"), pID)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})

	t.Run("skips -c flag if all roots are empty", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("@1 %2\n", 0, nil)
		expectedCmd := [][]string{
			{
				"tmux",
				"new-window",
				"-d",
				"-t",
				"mysession",
				"-n",
				"main",
				"-P", "-F", "#{window_id} #{pane_id}",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		win := window.Window{Name: "main", Panes: []pane.Pane{{}}}
		wID, pID, err := client.NewWindow("mysession", "", win)

		// then
		assert.Nil(t, err)
		assert.Equal(t, multiplexer.WindowID("@1"), wID)
		assert.Equal(t, multiplexer.PaneID("%2"), pID)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})
}

func Test_Client_ListSessions(t *testing.T) {
	t.Run("returns list of sessions", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("foo\nbar\nbaz\n", 0, nil).Once()

		expectedCmd := [][]string{
			{"tmux", "list-sessions", "-F", "#S"},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		sessions, err := client.ListSessions()

		// then
		assert.Nil(t, err)
		assert.Equal(t, []multiplexer.SessionName{"foo", "bar", "baz"}, sessions)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})

	t.Run("returns mapped error if command fails", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 1, errors.New("exit code 1")).Once()

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		_, err := client.ListSessions()

		// then
		assert.True(t, multiplexer.ErrFailedToListSessions.Equal(err))
		executor.AssertExpectations(t)
	})
}

func Test_IsTmuxServerRunning(t *testing.T) {
	t.Run("returns true if tmux server is running", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("tmux", 0, nil).Once()

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		running := client.IsTmuxServerRunning()

		// then
		assert.True(t, running)
		assert.Equal(t, executor.ExecutedCommands, [][]string{{"tmux", "run"}})
		executor.AssertExpectations(t)
	})

	t.Run("returns false if tmux server is not running", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("tmux", 1, errors.New("exit code 1")).Once()

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		running := client.IsTmuxServerRunning()

		// then
		assert.False(t, running)
		assert.Equal(t, executor.ExecutedCommands, [][]string{{"tmux", "run"}})
		executor.AssertExpectations(t)
	})
}

func Test_KillSession(t *testing.T) {
	t.Run("kills session", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 0, nil)
		expectedCmd := [][]string{
			{"tmux", "kill-session", "-t", "mysession"},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		err := client.KillSession("mysession")

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})
}

func Test_NewPane(t *testing.T) {
	t.Run("creates new pane with pane root", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("%1\n", 0, nil)
		expectedCmd := [][]string{
			{
				"tmux",
				"split-window",
				"-t",
				"@1",
				"-c",
				"/pane/root",
				"-P", "-F", "#{pane_id}",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		pID, err := client.NewPane("@1", "/template/root", "/window/root", pane.Pane{Root: "/pane/root"})

		// then
		assert.Nil(t, err)
		assert.Equal(t, multiplexer.PaneID("%1"), pID)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})

	t.Run("creates new pane with window root if pane root is empty", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("%1\n", 0, nil)
		expectedCmd := [][]string{
			{
				"tmux",
				"split-window",
				"-t",
				"@1",
				"-c",
				"/window/root",
				"-P", "-F", "#{pane_id}",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		pID, err := client.NewPane("@1", "/template/root", "/window/root", pane.Pane{})

		// then
		assert.Nil(t, err)
		assert.Equal(t, multiplexer.PaneID("%1"), pID)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})

	t.Run("creates new pane with template root if pane and window root are empty", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("%1\n", 0, nil)
		expectedCmd := [][]string{
			{
				"tmux",
				"split-window",
				"-t",
				"@1",
				"-c",
				"/template/root",
				"-P", "-F", "#{pane_id}",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		pID, err := client.NewPane("@1", "/template/root", "", pane.Pane{})

		// then
		assert.Nil(t, err)
		assert.Equal(t, multiplexer.PaneID("%1"), pID)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})

	t.Run("returns mapped error if command fails", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 1, errors.New("exit code 1"))

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		_, err := client.NewPane("@1", "/template/root", "", pane.Pane{})

		// then
		assert.True(t, multiplexer.ErrFailedToCreatePane.Equal(err))
		executor.AssertExpectations(t)
	})

	t.Run("skips -c flag if root is empty", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("%1\n", 0, nil)
		expectedCmd := [][]string{
			{
				"tmux",
				"split-window",
				"-t",
				"@1",
				"-P", "-F", "#{pane_id}",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		pID, err := client.NewPane("@1", "", "", pane.Pane{})

		// then
		assert.Nil(t, err)
		assert.Equal(t, multiplexer.PaneID("%1"), pID)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})
}

func Test_SetLayout(t *testing.T) {
	t.Run("sets layout", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 0, nil)
		expectedCmd := [][]string{
			{"tmux", "select-layout", "-t", "@1", "tiled"},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		err := client.SetLayout("@1", window.LayoutTiled)

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})

	t.Run("returns mapped error if command fails", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 1, errors.New("exit code 1"))

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		err := client.SetLayout("@1", window.LayoutTiled)

		// then
		assert.True(t, multiplexer.ErrFailedToSetLayout.Equal(err))
		executor.AssertExpectations(t)
	})
}

func Test_SetActivePane(t *testing.T) {
	t.Run("sets active pane", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 0, nil).Once()

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		err := client.SetActivePane(multiplexer.PaneID("%1"))

		// then
		assert.Nil(t, err)
		executor.AssertExpectations(t)
		assert.Equal(t, [][]string{{"tmux", "select-pane", "-t", "%1"}}, executor.ExecutedCommands)
	})

	t.Run("returns error if command fails", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 1, errors.New("some error")).Once()

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		err := client.SetActivePane(multiplexer.PaneID("%1"))

		// then
		assert.True(t, multiplexer.ErrFailedToSetActivePane.Equal(err))
		executor.AssertExpectations(t)
		assert.Equal(t, [][]string{{"tmux", "select-pane", "-t", "%1"}}, executor.ExecutedCommands)
	})
}

func Test_SetActiveWindow(t *testing.T) {
	t.Run("sets active window", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 0, nil).Once()
		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		err := client.SetActiveWindow("@1")

		// then
		assert.Nil(t, err)
		executor.AssertExpectations(t)
		assert.Equal(t, [][]string{{"tmux", "select-window", "-t", "@1"}}, executor.ExecutedCommands)
	})

	t.Run("returns error if command fails", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 1, errors.New("some error")).Once()

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		err := client.SetActiveWindow("@1")

		// then
		assert.True(t, multiplexer.ErrFailedToSetActiveWindow.Equal(err))
		executor.AssertExpectations(t)
		assert.Equal(t, [][]string{{"tmux", "select-window", "-t", "@1"}}, executor.ExecutedCommands)
	})
}
