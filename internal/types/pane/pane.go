package pane

import "thop/internal/types/command"

type Root string

type Pane struct {
	Root     Root              `yaml:"root,omitempty"`
	Commands []command.Command `yaml:"run,omitempty"`
	Active   bool              `yaml:"active,omitempty"`
}
