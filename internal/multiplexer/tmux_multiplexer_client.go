package multiplexer

import (
	"fmt"
	"os"
	"os/exec"
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
	HasSession(SessionName) (bool, error)
	IsTmuxServerRunning() bool
	KillSession(SessionName) error
	ListSessions() ([]SessionName, error)
	NewPane(WindowID, template.Root, window.Root, pane.Pane) (PaneID, error)
	NewSession(SessionName, template.Template) (WindowID, PaneID, error)
	NewWindow(SessionName, template.Root, window.Window) (WindowID, PaneID, error)
	SendKeys(PaneID, command.Command) error
	SetLayout(WindowID, window.Layout) error
	SwitchSession(SessionName) error
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
)

func (c *TmuxClientImpl) IsTmuxServerRunning() bool {
	cmd := exec.Command("tmux", "run")
	_, _, err := c.E.Execute(cmd)
	return err == nil
}

func (c *TmuxClientImpl) AttachSession(sn SessionName) error {
	cmd := exec.Command("tmux", "attach", "-t", string(sn))
	cmd.Stdin = os.Stdin // bind tmux session to terminal

	if _, _, err := c.E.Execute(cmd); err != nil {
		return buildExitCodeError(ErrFailedToAttachSession, err)
	}

	return nil
}

func (c *TmuxClientImpl) SwitchSession(sn SessionName) error {
	cmd := exec.Command("tmux", "switch", "-t", string(sn))

	if _, _, err := c.E.Execute(cmd); err != nil {
		return buildExitCodeError(ErrFailedToSwitchSession, err)
	}

	return nil
}

func (c *TmuxClientImpl) HasSession(sn SessionName) (bool, error) {
	cmd := exec.Command("tmux", "has-session", "-t", string(sn))

	_, exitCode, err := c.E.Execute(cmd)
	if err != nil {
		// exit code 1 means session does not exist
		if exitCode == 1 {
			return false, nil
		}
		return false, buildExitCodeError(ErrFailedToCheckSession, err)
	}

	return exitCode == 0, nil
}

func (c *TmuxClientImpl) NewSession(sn SessionName, t template.Template) (WindowID, PaneID, error) {
	mw := t.Windows[0]

	cmd := exec.Command("tmux", "new-session", "-d")
	cmd.Args = append(cmd.Args, "-s", string(sn))
	cmd.Args = append(cmd.Args, "-c", string(t.Root))

	if mw.Name != "" {
		cmd.Args = append(cmd.Args, "-n", string(mw.Name))
	}

	if mw.Root != "" {
		// little hack to start first window at different root than session
		cmd.Args = append(cmd.Args, fmt.Sprintf("cd %s && exec $SHELL", mw.Root))
	}

	cmd.Args = append(cmd.Args, "-P", "-F", "#{window_id} #{pane_id}")

	o, _, err := c.E.Execute(cmd)
	if err != nil {
		return WindowID(o), PaneID(o), buildExitCodeError(ErrFailedToCreateSession, err)
	}

	o = o[:len(o)-1] // trim newline character
	output := strings.Split(o, " ")
	return WindowID(output[0]), PaneID(output[1]), nil
}

func (c *TmuxClientImpl) NewWindow(sn SessionName, tr template.Root, w window.Window) (WindowID, PaneID, error) {
	cmd := exec.Command("tmux", "new-window", "-d")
	cmd.Args = append(cmd.Args, "-t", string(sn))

	if w.Name != "" {
		cmd.Args = append(cmd.Args, "-n", string(w.Name))
	}

	if w.Root != "" {
		cmd.Args = append(cmd.Args, "-c", string(w.Root))
	} else {
		// specify the root explicitly, otherwise it will nest at working directory
		cmd.Args = append(cmd.Args, "-c", string(tr))
	}

	cmd.Args = append(cmd.Args, "-P", "-F", "#{window_id} #{pane_id}")

	o, _, err := c.E.Execute(cmd)
	if err != nil {
		return WindowID(o), PaneID(o), buildExitCodeError(ErrFailedToCreateWindow, err)
	}

	o = o[:len(o)-1] // trim newline character
	output := strings.Split(o, " ")

	return WindowID(output[0]), PaneID(output[1]), nil
}

func (c *TmuxClientImpl) SendKeys(pID PaneID, keys command.Command) error {
	cmd := exec.Command("tmux", "send-keys", "-t", string(pID))
	cmd.Args = append(cmd.Args, string(keys))
	cmd.Args = append(cmd.Args, "C-m")

	if _, _, err := c.E.Execute(cmd); err != nil {
		return buildExitCodeError(ErrFailedToSendKeys, err)
	}

	return nil
}

func (c *TmuxClientImpl) ListSessions() ([]SessionName, error) {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#S")

	output, _, err := c.E.Execute(cmd)
	if err != nil {
		return nil, buildExitCodeError(ErrFailedToListSessions, err)
	}

	var sessionNames []SessionName
	for line := range strings.SplitSeq(output, "\n") {
		sessionNames = append(sessionNames, SessionName(line))
	}

	// drop the last one, it's empty since split is done by newline
	return sessionNames[:len(sessionNames)-1], nil
}

func (c *TmuxClientImpl) KillSession(sn SessionName) error {
	cmd := exec.Command("tmux", "kill-session", "-t", string(sn))

	if _, _, err := c.E.Execute(cmd); err != nil {
		return buildExitCodeError(ErrFailedToKillSession, err)
	}

	return nil
}

func (c *TmuxClientImpl) NewPane(
	wID WindowID,
	tr template.Root,
	wr window.Root,
	p pane.Pane,
) (PaneID, error) {
	cmd := exec.Command("tmux", "split-window", "-t", string(wID))

	var root string
	// pane root > window root > template root
	if p.Root != "" {
		root = string(p.Root)
	} else if wr != "" {
		root = string(wr)
	} else {
		root = string(tr)
	}

	cmd.Args = append(cmd.Args, "-c", root)
	cmd.Args = append(cmd.Args, "-P", "-F", "#{pane_id}")

	o, _, err := c.E.Execute(cmd)
	if err != nil {
		return PaneID(o), buildExitCodeError(ErrFailedToCreatePane, err)
	}

	o = o[:len(o)-1] // trim newline character
	return PaneID(o), nil
}

func (c *TmuxClientImpl) SetLayout(wID WindowID, l window.Layout) error {
	cmd := exec.Command("tmux", "select-layout", "-t", string(wID), string(l))
	if _, _, err := c.E.Execute(cmd); err != nil {
		return buildExitCodeError(ErrFailedToSetLayout, err)
	}

	return nil
}
