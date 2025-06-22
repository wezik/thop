package multiplexer_test

import (
	"errors"
	"os/exec"
	"testing"
	"thop/internal/multiplexer"
	"thop/internal/types/pane"
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
	t.Run("returns error if session name is empty", func(t *testing.T) {
		// given
		client := &multiplexer.TmuxClientImpl{}

		// when
		err := client.AttachSession("")

		// then
		assert.True(t, multiplexer.ErrInvalidTemplateArgs.Equal(err))
	})

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
	t.Run("returns error if session name is empty", func(t *testing.T) {
		// given
		client := multiplexer.TmuxClientImpl{
			E: nil,
		}

		// when
		err := client.SwitchSession("")

		// then
		assert.True(t, multiplexer.ErrInvalidTemplateArgs.Equal(err))
	})

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
	t.Run("returns error if session name is empty", func(t *testing.T) {
		// given
		client := multiplexer.TmuxClientImpl{
			E: nil,
		}

		// when
		_, err := client.HasSession("")

		// then
		assert.True(t, multiplexer.ErrInvalidTemplateArgs.Equal(err))
	})

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
	t.Run("returns error when missing required fields", func(t *testing.T) {
		// given
		client := multiplexer.TmuxClientImpl{
			E: nil,
		}

		win := window.Window{Name: "win"}

		// expect
		err := client.NewSession("", "root", win)
		assert.NotNil(t, err, "expected error for empty session name")

		// and
		err = client.NewSession("sess", "", win)
		assert.NotNil(t, err, "expected error for empty session root")

		// and
		err = client.NewSession("sess", "root", window.Window{})
		assert.NotNil(t, err, "expected error for empty window name")
	})

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
				"cd /project && exec $SHELL",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		win := window.Window{Name: "main", Root: "/project"}
		err := client.NewSession("mysession", "/home/test", win)

		// then
		assert.NotNil(t, err)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})

	t.Run("creates new session", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 0, nil)
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
				"cd /project && exec $SHELL",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		win := window.Window{Name: "main", Root: "/project"}
		err := client.NewSession("mysession", "/home/test", win)

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})

	t.Run("defaults main window root to session root if empty", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 0, nil)
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
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		win := window.Window{Name: "main", Root: ""}
		err := client.NewSession("mysession", "/home/test", win)

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})
}

func Test_TmuxClient_SendKeys(t *testing.T) {
	t.Run("returns error when missing required fields", func(t *testing.T) {
		// given
		client := multiplexer.TmuxClientImpl{
			E: nil,
		}

		// expect
		err := client.SendKeys("", "ls")
		assert.NotNil(t, err, "expected error for empty pane id")

		// and
		err = client.SendKeys("%id", "")
		assert.NotNil(t, err, "expected error for empty keys")
	})

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
	t.Run("returns error when missing required fields", func(t *testing.T) {
		// given
		client := multiplexer.TmuxClientImpl{
			E: nil,
		}

		// expect
		err := client.NewWindow("", "/project", window.Window{Name: "window", Root: "/root"})
		assert.NotNil(t, err, "expected error for empty session name")

		// and
		err = client.NewWindow("sess", "/project", window.Window{Name: "", Root: "/root"})
		assert.NotNil(t, err, "expected error for empty window name")

		// and
		err = client.NewWindow("sess", "", window.Window{Name: "window", Root: "/root"})
		assert.NotNil(t, err, "expected error for empty session root")
	})

	t.Run("includes default seession root if window root is empty", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 0, nil)
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
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		win := window.Window{Name: "main", Root: ""}
		err := client.NewWindow("mysession", "/home/test", win)

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})

	t.Run("creates new window", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 0, nil)
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
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		win := window.Window{Name: "main", Root: "/project"}
		err := client.NewWindow("mysession", "/home/test", win)

		// then
		assert.Nil(t, err)
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
	t.Run("returns error if session name is empty", func(t *testing.T) {
		// given
		client := multiplexer.TmuxClientImpl{
			E: nil,
		}

		// when
		err := client.KillSession("")

		// then
		assert.True(t, multiplexer.ErrInvalidTemplateArgs.Equal(err))
	})

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
	t.Run("returns error when missing required fields", func(t *testing.T) {
		// given
		client := multiplexer.TmuxClientImpl{
			E: nil,
		}

		// expect
		err := client.NewPane("", "win", pane.Pane{})
		assert.NotNil(t, err, "expected error for empty session name")

		// and
		err = client.NewPane("sess", "", pane.Pane{})
		assert.NotNil(t, err, "expected error for empty window name")
	})

	t.Run("creates new pane", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 0, nil)
		expectedCmd := [][]string{
			{
				"tmux",
				"split-window",
				"-d",
				"-t",
				"mysession:win",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		err := client.NewPane("mysession", "win", pane.Pane{})

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedCmd, executor.ExecutedCommands)
	})

	t.Run("creates new pane in directory", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 0, nil)
		expectedCmd := [][]string{
			{
				"tmux",
				"split-window",
				"-d",
				"-t",
				"mysession:win",
				"-c",
				"/project",
			},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		err := client.NewPane("mysession", "win", pane.Pane{Root: "/project"})

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
		err := client.NewPane("mysession", "win", pane.Pane{})

		// then
		assert.True(t, multiplexer.ErrFailedToCreatePane.Equal(err))
		executor.AssertExpectations(t)
	})
}

func Test_ListPanes(t *testing.T) {
	t.Run("returns error when missing required fields", func(t *testing.T) {
		// given
		client := multiplexer.TmuxClientImpl{
			E: nil,
		}

		// expect
		_, err := client.ListPanes("", "win")
		assert.NotNil(t, err, "expected error for empty session name")

		// and
		_, err = client.ListPanes("mysession", "")
		assert.NotNil(t, err, "expected error for empty window name")
	})

	t.Run("returns list of panes", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("foo\nbar\nbaz\n", 0, nil).Once()

		expectedCmd := [][]string{
			{"tmux", "list-panes", "-t", "mysession:win", "-F", "#{pane_id}"},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		panes, err := client.ListPanes("mysession", "win")

		// then
		assert.Nil(t, err)
		assert.Equal(t, []multiplexer.PaneID{"foo", "bar", "baz"}, panes)
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
		_, err := client.ListPanes("mysession", "win")

		// then
		assert.True(t, multiplexer.ErrFailedToListPanes.Equal(err))
		executor.AssertExpectations(t)
	})
}

func Test_SetLayout(t *testing.T) {
	t.Run("returns error if missing required fields", func(t *testing.T) {
		// given
		client := multiplexer.TmuxClientImpl{
			E: nil,
		}

		// expect
		err := client.SetLayout("", window.Window{})
		assert.NotNil(t, err, "expected error for empty session name")

		// and
		err = client.SetLayout("mysession", window.Window{})
		assert.NotNil(t, err, "expected error for empty window")
	})

	t.Run("sets layout", func(t *testing.T) {
		// given
		executor := new(MockCommandExecutor)
		executor.On("Execute", mock.Anything).Return("", 0, nil)
		expectedCmd := [][]string{
			{"tmux", "select-layout", "-t", "mysession:main", "tiled"},
		}

		client := multiplexer.TmuxClientImpl{
			E: executor,
		}

		// when
		err := client.SetLayout("mysession", window.Window{Name: "main", Layout: window.LayoutTiled})

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
		err := client.SetLayout("mysession", window.Window{Name: "main", Layout: window.LayoutTiled})

		// then
		assert.True(t, multiplexer.ErrFailedToSetLayout.Equal(err))
		executor.AssertExpectations(t)
	})
}
