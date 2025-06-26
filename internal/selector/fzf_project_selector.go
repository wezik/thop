package selector

import (
	"bytes"
	"fmt"
	"os/exec"
	"slices"
	"strings"
	"thop/internal/executor"
	"thop/internal/logger"
	"thop/internal/problem"
	"thop/internal/types/project"
)

type ProjectSelector interface {
	SelectFrom(items []project.Project, prompt string) (*project.Project, error)
}

type FzfProjectSelector struct {
	E executor.CommandExecutor
}

const (
	ErrSelectorCancelled problem.Key = "SELECTOR_CANCELLED"
	ErrSelectorFailed    problem.Key = "SELECTOR_FAILED"
	ErrUnexpectedState   problem.Key = "SELECTOR_UNEXPECTED_STATE"
)

type projectEntry struct {
	Project     *project.Project
	DisplayName string
	Prefix      string
	Order       int
}

func entryFromProject(p *project.Project) (projectEntry, error) {
	if p == nil {
		return projectEntry{}, ErrUnexpectedState.WithMsg("project is nil")
	}

	switch p.Type {
	case project.TypeTmuxSession:
		return projectEntry{
			Project:     p,
			DisplayName: string(p.Name),
			Prefix:      "(Active) ",
			Order:       0, // order active sessions first
		}, nil

	case project.TypeTemplate:
		var displayName string

		if p.Template.Name != "" {
			displayName = string(p.Template.Name)
		} else if p.Name != "" {
			displayName = string(p.Name)
		} else {
			return projectEntry{}, ErrUnexpectedState.WithMsg("project name cannot be empty")
		}

		return projectEntry{
			Project:     p,
			DisplayName: displayName,
			Order:       1,
		}, nil

	default:
		return projectEntry{}, ErrUnexpectedState.WithMsg("unhandled project type")
	}
}

func (s *FzfProjectSelector) SelectFrom(items []project.Project, prompt string) (*project.Project, error) {
	logger.Info(fmt.Sprintf("Selecting from %d projects", len(items)))

	var itemsInternal []projectEntry
	for _, item := range items {
		entry, err := entryFromProject(&item)
		if err != nil {
			return nil, err
		}
		itemsInternal = append(itemsInternal, entry)
	}

	slices.SortFunc(itemsInternal, func(a, b projectEntry) int {
		if a.Order != b.Order {
			// sort ascending by order first
			return b.Order - a.Order
		}

		aName := strings.ToLower(a.DisplayName)
		bName := strings.ToLower(b.DisplayName)

		// sort case-insensitive ascending
		return strings.Compare(bName, aName)
	})

	nameMap := make(map[string]*project.Project)
	var input bytes.Buffer

	for _, item := range itemsInternal {
		fullName := item.Prefix + item.DisplayName
		nameMap[fullName] = item.Project
		input.WriteString(fullName + "\n")
	}

	cmd := exec.Command("fzf")
	cmd.Stdin = &input
	cmd.Args = append(cmd.Args, "--prompt", prompt)

	output, exitCode, err := s.E.Execute(cmd)
	if exitCode == 130 {
		return nil, ErrSelectorCancelled.WithMsg("operation cancelled")
	} else if err != nil {
		return nil, ErrSelectorFailed.WithMsg(err.Error())
	}

	output = output[:len(output)-1] // trim newline character
	selected, ok := nameMap[output]
	if !ok {
		return nil, ErrUnexpectedState.WithMsg("selected project not found")
	}

	logger.Info(fmt.Sprintf("Selected project %s", selected.Name))
	return selected, nil
}
