package window

import (
	"thop/internal/types/command"
	"thop/internal/types/pane"
)

type Name string
type Root string

type Window struct {
	Name     Name              `yaml:"name,omitempty"`
	Root     Root              `yaml:"root,omitempty"`
	Commands []command.Command `yaml:"run,omitempty"`
	Panes    []pane.Pane       `yaml:"panes,omitempty"`
}
