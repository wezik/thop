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

var layouts = map[string]Layout{
	LayoutEvenHorizontal: LayoutEvenHorizontal,
	LayoutEvenVertical:   LayoutEvenVertical,
	LayoutMainHorizontal: LayoutMainHorizontal,
	LayoutMainVertical:   LayoutMainVertical,
	LayoutTiled:          LayoutTiled,
}

type Window struct {
	Name     Name            `yaml:"name,omitempty"`
	Root     Root            `yaml:"root,omitempty"`
	Commands command.Command `yaml:"run,omitempty"`
	Layout   Layout          `yaml:"layout,omitempty"`
	Panes    []pane.Pane     `yaml:"panes,omitempty"`
	Active   bool            `yaml:"active,omitempty"`
}

func (w *Window) WithDefaults() Window {
	newWindow := *w

	if len(newWindow.Panes) == 0 {
		newWindow.Panes = []pane.Pane{{}}
	}

	if newWindow.Layout == "" {
		newWindow.Layout = LayoutTiled
	}

	return newWindow
}

func (w *Window) IsValid() bool {
	if _, ok := layouts[string(w.Layout)]; !ok {
		return false
	}

	if len(w.Panes) < 1 { // at least main pane is required
		return false
	}

	for _, p := range w.Panes {
		if !p.IsValid() {
			return false
		}
	}

	return true
}
