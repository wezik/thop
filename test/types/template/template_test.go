package template_test

import (
	"testing"
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
		assert.Equal(t, pane.Name("window0-pane0"), temp.Windows[0].Panes[0].Name)
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
							Name: "bar",
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
		assert.Equal(t, pane.Name("bar"), template.Windows[0].Panes[0].Name)
		assert.Equal(t, pane.Name("foobar-pane0"), template.Windows[1].Panes[0].Name)
	})

	t.Run("fills window names", func(t *testing.T) {
		// given
		template := template.Template{
			Root: "/foo/bar",
			Windows: []window.Window{
				{
					Name: "",
					Panes: []pane.Pane{
						{
							Name: "bar",
						},
					},
				},
			},
		}

		// when
		newTemplate := template.WithDefaults()

		// then
		assert.Equal(t, window.Name("window0"), newTemplate.Windows[0].Name)
	})

	t.Run("fills pane names", func(t *testing.T) {
		// given
		template := template.Template{
			Root: "/foo/bar",
			Windows: []window.Window{
				{
					Name: "foo",
					Panes: []pane.Pane{
						{
							Name: "",
						},
					},
				},
			},
		}

		// when
		newTemplate := template.WithDefaults()

		// then
		assert.Equal(t, pane.Name("foo-pane0"), newTemplate.Windows[0].Panes[0].Name)
	})
}
