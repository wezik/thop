package window

import (
	"thop/internal/types/command"
	"thop/internal/types/pane"
)

type Name string
type Root string
type Layout string

const (
	LayoutEvenHorizontal = "even-horizontal"
	LayoutEvenVertical   = "even-vertical"
	LayoutMainHorizontal = "main-horizontal"
	LayoutMainVertical   = "main-vertical"
	LayoutTiled          = "tiled"
)

type Window struct {
	Name     Name              `yaml:"name,omitempty"`
	Root     Root              `yaml:"root,omitempty"`
	Commands []command.Command `yaml:"run,omitempty"`
	Layout   Layout            `yaml:"layout,omitempty"`
	Panes    []pane.Pane       `yaml:"panes,omitempty"`
	Active   bool              `yaml:"active,omitempty"`
}
