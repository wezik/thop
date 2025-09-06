package multiplexer

import (
	"thop/pkg/template"
	"thop/pkg/session"
)

type Multiplexer interface {
	AttachTemplate(*template.Template) (error)
	AttachSession(*session.Session) (error)
}
