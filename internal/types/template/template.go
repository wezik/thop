package template

import (
	"thop/internal/types/command"
	"thop/internal/types/window"
)

type Name string
type Root string
type ActiveWindow string

type Template struct {
	// Template name is used to specify the session name in multiplexer,
	// if not specified, the project name should be used
	Name         Name            `yaml:"name,omitempty"`
	Root         Root            `yaml:"root"`
	Commands     command.Command `yaml:"run,omitempty"`
	Windows      []window.Window `yaml:"windows,omitempty"`
	ActiveWindow ActiveWindow    `yaml:"active_window,omitempty"`
}

// Will set default values for missing fields
func (t *Template) WithDefaults() Template {
	newTemplate := *t

	if len(newTemplate.Windows) == 0 {
		newTemplate.Windows = []window.Window{{}}
	}

	for i := range newTemplate.Windows {
		newTemplate.Windows[i] = newTemplate.Windows[i].WithDefaults()
	}

	return newTemplate
}

func (t *Template) IsValid() bool {
	if len(t.Windows) < 1 { // at least main window is required
		return false
	}

	for _, win := range t.Windows {
		if !win.IsValid() {
			return false
		}
	}

	return true
}
