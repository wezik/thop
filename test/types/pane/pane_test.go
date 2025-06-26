package pane_test

import (
	"testing"
	"thop/internal/types/pane"

	"github.com/stretchr/testify/assert"
)

func TestPane_IsValid(t *testing.T) {
	t.Run("should always return true", func(t *testing.T) {
		// given
		p := pane.Pane{}

		// when
		isValid := p.IsValid()

		// then
		assert.True(t, isValid)
	})
}
