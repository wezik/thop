package multiplexer

import (
	"thop/pkg/session"
	"thop/pkg/template"
)

type Multiplexer interface {
	AttachTemplate(*template.Template) error
	AttachSession(*session.Session) error
}
