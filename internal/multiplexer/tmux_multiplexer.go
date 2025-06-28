package multiplexer

import (
	"fmt"
	"os"
	"slices"
	"thop/internal/logger"
	"thop/internal/types/pane"
	"thop/internal/types/project"
	"thop/internal/types/window"
)

type Multiplexer interface {
	AttachProject(project.Project) error
	ListActiveSessions() ([]project.Project, error)
	KillSession(project.Project) error
}

type TmuxMultiplexer struct {
	ActiveTmuxSession string
	Client            TmuxClient
}

func NewTmuxMultiplexer(c TmuxClient) Multiplexer {
	return &TmuxMultiplexer{Client: c, ActiveTmuxSession: os.Getenv("TMUX")}
}

func (m *TmuxMultiplexer) AttachProject(p project.Project) error {
	sn := resolveSessionName(p)
	exists, err := m.Client.HasSession(sn)
	if err != nil {
		return err
	}

	if !exists {
		if p.Type != project.TypeTemplate {
			return ErrTriedToBuildFromActiveSession.WithMsg("tried to build from non-template type project")
		}

		if err := m.assembleSession(sn, p); err != nil {
			return err
		}
	}

	if m.ActiveTmuxSession != "" {
		logger.Message(fmt.Sprintf("Switching to %s session", sn))
		if err := m.Client.SwitchSession(sn); err != nil {
			return err
		}
	} else {
		logger.Message(fmt.Sprintf("Attaching to %s session", sn))
		if err := m.Client.AttachSession(sn); err != nil {
			return err
		}
	}

	return nil
}

func (m *TmuxMultiplexer) ListActiveSessions() ([]project.Project, error) {
	if !m.Client.IsTmuxServerRunning() {
		return nil, nil
	}

	sessionNames, err := m.Client.ListSessions()
	if err != nil {
		return nil, err
	}

	var tmuxProjects []project.Project
	for _, sessionName := range sessionNames {
		// simple project with only name and type signaling its not regular template
		p := project.Project{
			Name: project.Name(sessionName),
			Type: project.TypeTmuxSession,
		}

		tmuxProjects = append(tmuxProjects, p)
	}

	logger.Info(fmt.Sprintf("Loaded %d active tmux sessions", len(tmuxProjects)))
	return tmuxProjects, nil
}

func (m *TmuxMultiplexer) KillSession(p project.Project) error {
	sn := resolveSessionName(p)
	err := m.Client.KillSession(sn)
	logger.Message(fmt.Sprintf("Killed session %s", sn))
	return err
}

func (m *TmuxMultiplexer) assembleSession(sn SessionName, p project.Project) (err error) {
	windows := map[WindowID]window.Window{}
	windowIdsToPaneIds := map[WindowID][]PaneID{}
	panes := map[PaneID]pane.Pane{}

	// first window gets created together with the session
	wID, pID, err := m.Client.NewSession(sn, p.Template)
	if err != nil {
		return err
	}

	// save main window and main pane
	windows[wID] = p.Template.Windows[0]
	panes[pID] = p.Template.Windows[0].Panes[0]
	windowIdsToPaneIds[wID] = []PaneID{pID}

	// from this point on, if any error occurs, attempt to kill the session
	defer func() {
		if err != nil {
			m.Client.KillSession(sn)
		}
	}()

	t := p.Template

	// initialize additional windows
	// first window is created together with the session, so skip it
	for _, w := range t.Windows[1:] {
		wID, pID, err := m.Client.NewWindow(sn, t.Root, w)
		if err != nil {
			return err
		}

		windows[wID] = w
		// save main pane of the window
		panes[pID] = w.Panes[0]
		windowIdsToPaneIds[wID] = []PaneID{pID}
	}

	// initialize additional panes
	for wID, w := range windows {
		// first pane is created together with the window, so skip it
		for _, p := range w.Panes[1:] {
			pID, err := m.Client.NewPane(wID, t.Root, w.Root, p)
			if err != nil {
				return err
			}

			panes[pID] = p
			windowIdsToPaneIds[wID] = append(windowIdsToPaneIds[wID], pID)

			// set layout after each pane to ensure, layout is up to its limits
			if err = m.Client.SetLayout(wID, w.Layout); err != nil {
				return err
			}
		}
	}

	// execute shell commands
	for wID, pIDs := range windowIdsToPaneIds {
		for _, pID := range pIDs {
			w := windows[wID]
			p := panes[pID]

			tKeys := commandToKeys(t.Commands)
			wKeys := commandToKeys(w.Commands)
			pKeys := commandToKeys(p.Commands)

			cmds := slices.Concat(tKeys, wKeys, pKeys)
			for _, keys := range cmds {
				if err = m.Client.SendKeys(pID, keys); err != nil {
					return err
				}
			}
		}
	}

	// set active windows and panes
	for wID, pIDs := range windowIdsToPaneIds {

		for _, pID := range pIDs {
			p := panes[pID]
			if p.Active {
				if err = m.Client.SetActivePane(pID); err != nil {
					return err
				}
			}
		}

		w := windows[wID]
		if w.Active {
			if err = m.Client.SetActiveWindow(wID); err != nil {
				return err
			}
		}
	}

	logger.Message(fmt.Sprintf("Session %s assembled", sn))
	return nil
}
