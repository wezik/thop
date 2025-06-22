package multiplexer

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"thop/internal/executor"
	"thop/internal/problem"
	"thop/internal/types/command"
	"thop/internal/types/pane"
	"thop/internal/types/template"
	"thop/internal/types/window"
)

type TmuxClient interface {
	AttachSession(SessionName) error
	SwitchSession(SessionName) error
	HasSession(SessionName) (bool, error)
	NewSession(SessionName, template.Root, window.Window) error
	NewWindow(SessionName, template.Root, window.Window) error
	SendKeys(SessionName, window.Name, command.Command) error
	ListSessions() ([]SessionName, error)
	IsTmuxServerRunning() bool
	KillSession(SessionName) error
	NewPane(SessionName, window.Name, pane.Pane) error
	SetLayout(SessionName, window.Window) error
}

type TmuxClientImpl struct {
	E executor.CommandExecutor
}

const (
	ErrFailedToAttachSession         problem.Key = "TMUX_FAILED_TO_ATTACH_SESSION"
	ErrFailedToSwitchSession         problem.Key = "TMUX_FAILED_TO_SWITCH_SESSION"
	ErrFailedToCheckSession          problem.Key = "TMUX_FAILED_TO_CHECK_SESSION"
	ErrFailedToCreateSession         problem.Key = "TMUX_FAILED_TO_CREATE_SESSION"
	ErrFailedToCreateWindow          problem.Key = "TMUX_FAILED_TO_CREATE_WINDOW"
	ErrFailedToListSessions          problem.Key = "TMUX_FAILED_TO_LIST_SESSIONS"
	ErrFailedToKillSession           problem.Key = "TMUX_FAILED_TO_KILL_SESSION"
	ErrFailedToSendKeys              problem.Key = "TMUX_FAILED_TO_SEND_KEYS"
	ErrFailedToSetLayout             problem.Key = "TMUX_FAILED_TO_SET_LAYOUT"
	ErrTriedToBuildFromActiveSession problem.Key = "TMUX_TRIED_TO_BUILD_FROM_ACTIVE_SESSION"
	ErrFailedToCreatePane            problem.Key = "TMUX_FAILED_TO_CREATE_PANE"
	ErrInvalidTemplateArgs           problem.Key = "TMUX_INVALID_TEMPLATE_ARGS"
)

func (c *TmuxClientImpl) IsTmuxServerRunning() bool {
	cmd := exec.Command("tmux", "run")
	_, _, err := c.E.Execute(cmd)
	return err == nil
}

func (c *TmuxClientImpl) AttachSession(session SessionName) error {
	if session == "" {
		return ErrInvalidTemplateArgs.WithMsg("session name cannot be empty")
	}

	cmd := exec.Command("tmux", "attach", "-t", string(session))
	cmd.Stdin = os.Stdin // bind tmux session to terminal

	_, _, err := c.E.Execute(cmd)
	if err != nil {
		switch err := err.(type) {
		case *exec.ExitError:
			return ErrFailedToAttachSession.WithMsg(string(err.Stderr))
		default:
			return ErrFailedToAttachSession.WithMsg(err.Error())
		}
	}

	return nil
}

func (c *TmuxClientImpl) SwitchSession(session SessionName) error {
	if session == "" {
		return ErrInvalidTemplateArgs.WithMsg("session name cannot be empty")
	}

	cmd := exec.Command("tmux", "switch", "-t", string(session))

	_, _, err := c.E.Execute(cmd)
	if err != nil {
		switch err := err.(type) {
		case *exec.ExitError:
			return ErrFailedToSwitchSession.WithMsg(string(err.Stderr))
		default:
			return ErrFailedToSwitchSession.WithMsg(err.Error())
		}
	}

	return nil
}

func (c *TmuxClientImpl) HasSession(session SessionName) (bool, error) {
	if session == "" {
		return false, ErrInvalidTemplateArgs.WithMsg("session name cannot be empty")
	}

	cmd := exec.Command("tmux", "has-session", "-t", string(session))

	_, exitCode, err := c.E.Execute(cmd)
	if err != nil {
		// exit code 1 means session does not exist
		if exitCode == 1 {
			return false, nil
		}
		switch err := err.(type) {
		case *exec.ExitError:
			return false, ErrFailedToCheckSession.WithMsg(string(err.Stderr))
		default:
			return false, ErrFailedToCheckSession.WithMsg(err.Error())
		}
	}

	return exitCode == 0, nil
}

func (c *TmuxClientImpl) NewSession(
	session SessionName,
	root template.Root,
	mainWindow window.Window,
) error {
	if anyEmpty(string(session), string(root), string(mainWindow.Name)) {
		return ErrInvalidTemplateArgs.WithMsg("session, root and window name cannot be empty")
	}

	cmd := exec.Command("tmux", "new-session", "-d")
	cmd.Args = append(cmd.Args, "-s", string(session))
	cmd.Args = append(cmd.Args, "-c", string(root))
	cmd.Args = append(cmd.Args, "-n", string(mainWindow.Name))

	if mainWindow.Root != "" {
		// little hack to start first window at different root than session
		cmd.Args = append(cmd.Args, fmt.Sprintf("cd %s && exec $SHELL", mainWindow.Root))
	}

	if _, _, err := c.E.Execute(cmd); err != nil {
		switch err := err.(type) {
		case *exec.ExitError:
			return ErrFailedToCreateSession.WithMsg(string(err.Stderr))
		default:
			return ErrFailedToCreateSession.WithMsg(err.Error())
		}
	}

	return nil
}

func (c *TmuxClientImpl) NewWindow(
	session SessionName,
	root template.Root,
	mainWindow window.Window,
) error {
	if anyEmpty(string(session), string(root), string(mainWindow.Name)) {
		return ErrInvalidTemplateArgs.WithMsg("session, root and window name cannot be empty")
	}

	cmd := exec.Command("tmux", "new-window", "-d")
	cmd.Args = append(cmd.Args, "-t", string(session))
	cmd.Args = append(cmd.Args, "-n", string(mainWindow.Name))

	if mainWindow.Root != "" {
		cmd.Args = append(cmd.Args, "-c", string(mainWindow.Root))
	} else {
		// in certain scenarios tmux will create window in working directory
		// instead of the session root, so specify it explicitly
		cmd.Args = append(cmd.Args, "-c", string(root))
	}

	if _, _, err := c.E.Execute(cmd); err != nil {
		switch err := err.(type) {
		case *exec.ExitError:
			return ErrFailedToCreateWindow.WithMsg(string(err.Stderr))
		default:
			return ErrFailedToCreateWindow.WithMsg(err.Error())
		}
	}

	return nil

}

func (c *TmuxClientImpl) SendKeys(
	session SessionName,
	windowName window.Name,
	keys command.Command,
) error {
	if anyEmpty(string(session), string(windowName), string(keys)) {
		return ErrInvalidTemplateArgs.WithMsg("session, window name and keys cannot be empty")
	}

	cmd := exec.Command("tmux", "send-keys")

	// tmux needs combined name of session:window to send keys to
	cmd.Args = append(cmd.Args, "-t", fmt.Sprintf("%s:%s", session, windowName))
	cmd.Args = append(cmd.Args, string(keys))
	cmd.Args = append(cmd.Args, "C-m")

	if _, _, err := c.E.Execute(cmd); err != nil {
		switch err := err.(type) {
		case *exec.ExitError:
			return ErrFailedToSendKeys.WithMsg(string(err.Stderr))
		default:
			return ErrFailedToSendKeys.WithMsg(err.Error())
		}
	}

	return nil
}

func (c *TmuxClientImpl) ListSessions() ([]SessionName, error) {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#S")

	output, _, err := c.E.Execute(cmd)
	if err != nil {
		switch err := err.(type) {
		case *exec.ExitError:
			return nil, ErrFailedToListSessions.WithMsg(string(err.Stderr))
		default:
			return nil, ErrFailedToListSessions.WithMsg(err.Error())
		}
	}

	var sessionNames []SessionName
	for line := range strings.SplitSeq(output, "\n") {
		sessionNames = append(sessionNames, SessionName(line))
	}

	// drop the last one, it's empty
	return sessionNames[:len(sessionNames)-1], nil
}

func (c *TmuxClientImpl) KillSession(session SessionName) error {
	if session == "" {
		return ErrInvalidTemplateArgs.WithMsg("session name cannot be empty")
	}

	cmd := exec.Command("tmux", "kill-session", "-t", string(session))

	_, _, err := c.E.Execute(cmd)
	if err != nil {
		switch err := err.(type) {
		case *exec.ExitError:
			return ErrFailedToKillSession.WithMsg(string(err.Stderr))
		default:
			return ErrFailedToKillSession.WithMsg(err.Error())
		}
	}

	return nil
}

func (c *TmuxClientImpl) NewPane(session SessionName, windowName window.Name, p pane.Pane) error {
	if anyEmpty(string(session), string(windowName)) {
		return ErrInvalidTemplateArgs.WithMsg("session and window name cannot be empty")
	}

	combinedName := fmt.Sprintf("%s:%s", session, windowName)

	cmd := exec.Command("tmux", "split-window", "-t", combinedName)

	if p.Root != "" {
		cmd.Args = append(cmd.Args, "-c", string(p.Root))
	}

	fmt.Println("running", cmd.Args)

	if _, _, err := c.E.Execute(cmd); err != nil {
		switch err := err.(type) {
		case *exec.ExitError:
			return ErrFailedToCreatePane.WithMsg(string(err.Stderr))
		default:
			return ErrFailedToCreatePane.WithMsg(err.Error())
		}
	}


	return nil
}

func (c *TmuxClientImpl) SetLayout(session SessionName, window window.Window) error {
	if anyEmpty(string(session), string(window.Name), string(window.Layout)) {
		return ErrInvalidTemplateArgs.WithMsg("session, window name and layout cannot be empty")
	}

	combinedName := fmt.Sprintf("%s:%s", session, window.Name)

	cmd := exec.Command("tmux", "select-layout", "-t", combinedName, string(window.Layout))
	_, _, err := c.E.Execute(cmd)
	if err != nil {
		switch err := err.(type) {
		case *exec.ExitError:
			return ErrFailedToSetLayout.WithMsg(string(err.Stderr))
		default:
			return ErrFailedToSetLayout.WithMsg(err.Error())
		}
	}

	return nil
}

func anyEmpty(s ...string) bool {
	return slices.Contains(s, "")
}
