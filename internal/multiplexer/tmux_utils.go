package multiplexer

import (
	"os/exec"
	"strings"
	"thop/internal/problem"
	"thop/internal/types/command"
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

func commandToKeys(cmd command.Command) []string {
	var keys []string

	for c := range strings.SplitSeq(string(cmd), "\n") {
		if c != "" {
			keys = append(keys, c)
		}
	}

	return keys
}
