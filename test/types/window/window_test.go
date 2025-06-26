package window_test

import (
	"testing"
	"thop/internal/types/pane"
	"thop/internal/types/window"

	"github.com/stretchr/testify/assert"
)

func TestWindow_WithDefaults(t *testing.T) {
	t.Run("should set default values for empty window", func(t *testing.T) {
		// given
		w := window.Window{}

		// when
		newW := w.WithDefaults()

		// then
		assert.Equal(t, window.Layout(window.LayoutTiled), newW.Layout)
		assert.Len(t, newW.Panes, 1)
		assert.Equal(t, pane.Pane{}, newW.Panes[0])
	})

	t.Run("should not override existing values", func(t *testing.T) {
		// given
		w := window.Window{
			Name:   "my-window",
			Layout: window.LayoutEvenHorizontal,
			Panes: []pane.Pane{
				{Root: "/tmp"},
			},
		}

		// when
		newW := w.WithDefaults()

		// then
		assert.Equal(t, window.Name("my-window"), newW.Name)
		assert.Equal(t, window.Layout(window.LayoutEvenHorizontal), newW.Layout)
		assert.Len(t, newW.Panes, 1)
		assert.Equal(t, pane.Pane{Root: "/tmp"}, newW.Panes[0])
	})
}

func TestWindow_IsValid(t *testing.T) {
	t.Run("should return true for a valid window", func(t *testing.T) {
		// given
		w := window.Window{
			Layout: window.LayoutTiled,
			Panes:  []pane.Pane{{}},
		}

		// when
		isValid := w.IsValid()

		// then
		assert.True(t, isValid)
	})

	t.Run("should return false for an invalid layout", func(t *testing.T) {
		// given
		w := window.Window{
			Layout: "invalid-layout",
			Panes:  []pane.Pane{{}},
		}

		// when
		isValid := w.IsValid()

		// then
		assert.False(t, isValid)
	})

	t.Run("should return false if there are no panes", func(t *testing.T) {
		// given
		w := window.Window{
			Layout: window.LayoutTiled,
			Panes:  []pane.Pane{},
		}

		// when
		isValid := w.IsValid()

		// then
		assert.False(t, isValid)
	})

	t.Run("name is not required", func(t *testing.T) {
		// given
		w := window.Window{
			Layout: window.LayoutTiled,
			Panes:  []pane.Pane{{}},
		}

		// when
		isValid := w.IsValid()

		// then
		assert.True(t, isValid)
	})
}
