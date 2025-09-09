package template

import (
	"testing"
	"thop/internal/domain/template"

	"github.com/stretchr/testify/assert"
)

func Test_TemplateSessionName(t *testing.T) {
	t.Run("defaults to name if session name is empty", func(t *testing.T) {
		// given
		templ := template.NewTemplate(
			"foo/bar/1",
			template.V1,
			"1",
			"",
			"",
			[]string{},
			[]*template.Window{},
		)

		// when
		sessionName := templ.SessionName()

		// then
		assert.Equal(t, templ.Name(), sessionName)
	})

	t.Run("returns session name if set", func(t *testing.T) {
		// given
		templ := template.NewTemplate(
			"foo/bar/1",
			template.V1,
			"1",
			"session_name",
			"",
			[]string{},
			[]*template.Window{},
		)

		// when
		sessionName := templ.SessionName()

		// then
		assert.Equal(t, "session_name", sessionName)
	})
}
