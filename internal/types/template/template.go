package template

import (
	"strconv"
	"thop/internal/types/command"
	"thop/internal/types/pane"
	"thop/internal/types/window"
)

type Name string
type Root string
type ActiveWindow string

type Template struct {
	// Template name is used to specify the session name in multiplexer,
	// if not specified, the project name should be used
	Name         Name              `yaml:"name,omitempty"`
	Root         Root              `yaml:"root"`
	Commands     []command.Command `yaml:"run,omitempty"`
	Windows      []window.Window   `yaml:"windows,omitempty"`
	ActiveWindow ActiveWindow      `yaml:"active_window,omitempty"`
}

// Will set default values for missing fields
func (t *Template) WithDefaults() Template {
	newTemplate := *t

	if len(newTemplate.Windows) == 0 {
		newTemplate.Windows = []window.Window{{}}
	}

	for i := range newTemplate.Windows {
		win := &newTemplate.Windows[i]

		if win.Name == "" {
			win.Name = window.Name("window" + strconv.Itoa(i))
		}

		if len(win.Panes) == 0 {
			win.Panes = []pane.Pane{
				{
					Name: pane.Name(string(win.Name) + "-pane" + strconv.Itoa(0)),
				},
			}
		}

		for j := range win.Panes {
			pan := &win.Panes[j]

			if pan.Name == "" {
				pan.Name = pane.Name(string(win.Name) + "-pane" + strconv.Itoa(j))
			}
		}
	}

	return newTemplate
}
