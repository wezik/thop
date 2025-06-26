package multiplexer

import (
	"os/exec"
	"thop/internal/problem"
	"thop/internal/types/project"
)

type SessionName string
type PaneID string
type WindowID string

func buildExitCodeError(key problem.Key, err error) error {
	switch err := err.(type) {
	case *exec.ExitError:
		return key.WithMsg(string(err.Stderr))
	default:
		return key.WithMsg(err.Error())
	}
}

func resolveSessionName(p project.Project) SessionName {
	if p.Template.Name != "" {
		return SessionName(p.Template.Name)
	}

	return SessionName(p.Name)
}
