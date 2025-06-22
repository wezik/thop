package multiplexer

import (
	"fmt"
	"slices"
	"thop/internal/types/project"
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

type SessionName string

func (m *TmuxMultiplexer) AttachProject(p project.Project) error {
	sessionName, err := resolveSessionName(p)
	if err != nil {
		return err
	}

	sessionExists, err := m.Client.HasSession(sessionName)
	if err != nil {
		return err
	}

	if !sessionExists {
		if p.Type == project.TypeTmuxSession {
			return ErrTriedToBuildFromActiveSession.WithMsg("cannot build from active session (it was probably killed while thop was running)")
		}

		// set default values for missing fields if needed
		p.Template = p.Template.WithDefaults()

		if err := m.assembleSession(sessionName, p); err != nil {
			return err
		}
	}

	if m.ActiveTmuxSession != "" {
		fmt.Println("Switching to", sessionName, "session")
		if err := m.Client.SwitchSession(sessionName); err != nil {
			return err
		}
	} else {
		fmt.Println("Attaching to", sessionName, "session")
		if err := m.Client.AttachSession(sessionName); err != nil {
			return err
		}
	}

	return nil
}

func (m *TmuxMultiplexer) ListActiveSessions() ([]project.Project, error) {
	if !m.Client.IsTmuxServerRunning() {
		return []project.Project(nil), nil
	}

	sessionNames, err := m.Client.ListSessions()
	if err != nil {
		return nil, err
	}

	var tmuxProjects []project.Project
	for _, sessionName := range sessionNames {
		tmuxProjects = append(tmuxProjects, project.Project{Name: project.Name(sessionName), Type: project.TypeTmuxSession})
	}

	return tmuxProjects, nil
}

func (m *TmuxMultiplexer) KillSession(p project.Project) error {
	sessionName, err := resolveSessionName(p)
	if err != nil {
		return err
	}

	if err := m.Client.KillSession(sessionName); err != nil {
		return err
	}

	return nil
}

func (m *TmuxMultiplexer) assembleSession(sessionName SessionName, pro project.Project) (err error) {
	sessionRoot := pro.Template.Root
	if sessionRoot == "" {
		return ErrInvalidTemplateArgs.WithMsg("session root cannot be empty")
	}

	if len(pro.Template.Windows) == 0 {
		return ErrInvalidTemplateArgs.WithMsg("project template needs at least one window to be created")
	}

	mainWindow := pro.Template.Windows[0]

	// first window gets created together with the session
	if err = m.Client.NewSession(sessionName, sessionRoot, mainWindow); err != nil {
		return err
	}

	// from this point on, if any error occurs, kill the session
	defer func() {
		if err != nil {
			m.Client.KillSession(sessionName)
		}
	}()

	for i, window := range pro.Template.Windows {
		// first window is created together with the session, so skip it
		if i != 0 {
			if err = m.Client.NewWindow(sessionName, sessionRoot, window); err != nil {
				return err
			}
		}

		// first pane is created together with the window, so skip it
		for _, p := range window.Panes[1:] {
			if err = m.Client.NewPane(sessionName, window.Name, p); err != nil {
				return err
			}

			// set layout after each pane to ensure, layout is up to its limits
			if err = m.Client.SetLayout(sessionName, window); err != nil {
				return err
			}
		}

		paneIDs, err := m.Client.ListPanes(sessionName, window.Name)
		if err != nil {
			return err
		}

		for paneIndex, paneID := range paneIDs {
			if len(paneIDs) != len(window.Panes) {
				panic(fmt.Errorf("invalid state: panes count does not match window panes count"))
			}

			commands := slices.Concat(pro.Template.Commands, window.Commands, window.Panes[paneIndex].Commands)

			for _, keys := range commands {
				if err = m.Client.SendKeys(paneID, keys); err != nil {
					return err
				}
			}
		}
	}

	fmt.Println("Session", sessionName, "created")
	return nil
}

func resolveSessionName(p project.Project) (SessionName, error) {
	if p.Template.Name != "" {
		return SessionName(p.Template.Name), nil
	}

	if p.Name == "" {
		return "", ErrInvalidTemplateArgs.WithMsg("project name cannot be empty")
	}

	return SessionName(p.Name), nil
}
