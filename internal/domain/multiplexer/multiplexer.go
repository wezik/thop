package multiplexer

import (
	"thop/internal/domain/session"
	"thop/internal/domain/template"
)

type Multiplexer interface {
	AttachTemplate(*template.Template) error
	AttachSession(*session.Session) error
}
