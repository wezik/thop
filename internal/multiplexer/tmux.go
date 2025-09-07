package multiplexer

import (
	"thop/pkg/log"
	"thop/pkg/multiplexer"
	"thop/pkg/session"
	"thop/pkg/template"
)

type TmuxMultiplexer struct {
	log log.Logger
}

func NewTmuxMultiplexer(
	log log.Logger,
) multiplexer.Multiplexer {
	return &TmuxMultiplexer{
		log: log,
	}
}

func (m *TmuxMultiplexer) AttachTemplate(templ *template.Template) (err error) {
	m.log.Debug("Attaching to template \"" + string(templ.FilePath()) + "\"")
	return
}

func (m *TmuxMultiplexer) AttachSession(ses *session.Session) (err error) {
	m.log.Debug("Attaching to session \"" + ses.Name() + "\"")
	return
}
