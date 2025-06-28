package project_test

import (
	"testing"
	"thop/internal/types/pane"
	"thop/internal/types/project"
	"thop/internal/types/template"
	"thop/internal/types/window"

	"github.com/stretchr/testify/assert"
)

func TestProject_WithDefaults(t *testing.T) {
	t.Run("should set defaults for template", func(t *testing.T) {
		// given
		p := project.Project{
			Template: template.Template{
				Root: "/foo",
			},
		}

		// when
		newP := p.WithDefaults()

		// then
		assert.Len(t, newP.Template.Windows, 1)
		assert.Equal(t, window.Layout(window.LayoutTiled), newP.Template.Windows[0].Layout)
		assert.Len(t, newP.Template.Windows[0].Panes, 1)
	})
}

func TestProject_IsValid(t *testing.T) {
	t.Run("should return true for minimum valid project", func(t *testing.T) {
		// given
		p := project.Project{
			UUID: "some-uuid",
			Name: "my-project",
			Type: project.TypeTemplate,
			Template: template.Template{
				Root: "/foo",
				Windows: []window.Window{
					{
						Layout: window.LayoutTiled,
						Panes:  []pane.Pane{{}},
					},
				},
			},
		}

		// when
		isValid := p.IsValid()

		// then
		assert.True(t, isValid)
	})

	t.Run("should return false for empty UUID", func(t *testing.T) {
		// given
		p := project.Project{
			UUID: "",
			Name: "my-project",
			Type: project.TypeTemplate,
			Template: template.Template{
				Root: "/foo",
				Windows: []window.Window{
					{
						Layout: window.LayoutTiled,
						Panes:  []pane.Pane{{}},
					},
				},
			},
		}

		// when
		isValid := p.IsValid()

		// then
		assert.False(t, isValid)
	})

	t.Run("should return false for empty name", func(t *testing.T) {
		// given
		p := project.Project{
			UUID: "some-uuid",
			Name: "",
			Type: project.TypeTemplate,
			Template: template.Template{
				Root: "/foo",
				Windows: []window.Window{
					{
						Layout: window.LayoutTiled,
						Panes:  []pane.Pane{{}},
					},
				},
			},
		}

		// when
		isValid := p.IsValid()

		// then
		assert.False(t, isValid)
	})

	t.Run("should return false for invalid template", func(t *testing.T) {
		// given
		p := project.Project{
			UUID:     "some-uuid",
			Name:     "my-project",
			Type:     project.TypeTemplate,
			Template: template.Template{
				// Root is missing
			},
		}

		// when
		isValid := p.IsValid()

		// then
		assert.False(t, isValid)
	})

	t.Run("should return false for non-template project type", func(t *testing.T) {
		// given
		p := project.Project{
			Type: project.TypeTmuxSession,
		}

		// when
		isValid := p.IsValid()

		// then
		assert.False(t, isValid)
	})
}
