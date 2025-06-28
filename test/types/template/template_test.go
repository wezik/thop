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
		assert.Equal(t, pane.Pane{}, temp.Windows[0].Panes[0])
	})

	t.Run("fills windows with panes", func(t *testing.T) {
		// given
		temp := template.Template{
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
		temp = temp.WithDefaults()

		// then
		assert.Equal(t, pane.Pane{Commands: []command.Command{"echo foo"}}, temp.Windows[0].Panes[0])
		assert.Equal(t, pane.Pane{}, temp.Windows[1].Panes[0])
	})
}

func Test_IsValid(t *testing.T) {
	t.Run("returns true for valid template", func(t *testing.T) {
		// given
		temp := template.Template{
			Windows: []window.Window{
				{
					Layout: window.LayoutTiled,
					Panes:  []pane.Pane{{}},
				},
			},
		}

		// when
		isValid := temp.IsValid()

		// then
		assert.True(t, isValid)
	})

	t.Run("returns false if no windows", func(t *testing.T) {
		// given
		temp := template.Template{
			Root:    "/foo/bar",
			Windows: []window.Window{},
		}

		// when
		isValid := temp.IsValid()

		// then
		assert.False(t, isValid)
	})

	t.Run("returns false if any window is invalid", func(t *testing.T) {
		// given
		temp := template.Template{
			Root: "/foo/bar",
			Windows: []window.Window{
				{
					Layout: "invalid-layout", // Invalid layout
					Panes:  []pane.Pane{{}},
				},
			},
		}

		// when
		isValid := temp.IsValid()

		// then
		assert.False(t, isValid)
	})
}
