package template_test

import (
	"testing"
	"thop/internal/types/command"
	"thop/internal/types/pane"
	"thop/internal/types/template"
	"thop/internal/types/window"

	"github.com/stretchr/testify/assert"
)

func Test_WithDefaults(t *testing.T) {
	t.Run("creates new template with defaults", func(t *testing.T) {
		// given
		temp := template.Template{
			Root: "/foo/bar",
		}

		// when
		temp = temp.WithDefaults()

		// then
		assert.Equal(t, template.Root("/foo/bar"), temp.Root)
		assert.Equal(t, window.Name("window0"), temp.Windows[0].Name)
		assert.Equal(t, pane.Pane{}, temp.Windows[0].Panes[0])
	})

	t.Run("fills windows with panes", func(t *testing.T) {
		// given
		template := template.Template{
			Root: "/foo/bar",
			Windows: []window.Window{
				{
					Name: "foo",
					Panes: []pane.Pane{
						{
							Commands: []command.Command{"echo foo"},
						},
					},
				},
				{
					Name: "foobar",
				},
			},
		}

		// when
		template = template.WithDefaults()

		// then
		assert.Equal(t, pane.Pane{Commands: []command.Command{"echo foo"}}, template.Windows[0].Panes[0])
		assert.Equal(t, pane.Pane{}, template.Windows[1].Panes[0])
	})

	t.Run("fills window names", func(t *testing.T) {
		// given
		template := template.Template{
			Root: "/foo/bar",
			Windows: []window.Window{
				{
					Name: "",
				},
			},
		}

		// when
		newTemplate := template.WithDefaults()

		// then
		assert.Equal(t, window.Name("window0"), newTemplate.Windows[0].Name)
	})
}
