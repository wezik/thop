package thop

import (
	"thop/internal/domain/log"
	"thop/internal/domain/multiplexer"
	"thop/internal/domain/selector"
	"thop/internal/domain/session"
	"thop/internal/domain/template"
)

type Search struct {
	Phrase string
}

type Thop struct {
	log             log.Logger
	selector        selector.Selector
	templateService template.TemplateService
	multiplexer     multiplexer.Multiplexer
}

type DefaultDirectory = string

func New(
	log log.Logger,
	selector selector.Selector,
	templateService template.TemplateService,
	multiplexer multiplexer.Multiplexer,
) *Thop {
	return &Thop{
		log:             log,
		selector:        selector,
		templateService: templateService,
		multiplexer:     multiplexer,
	}
}

func (t *Thop) Create(act template.CreateTemplate) (templ *template.Template, err error) {
	t.log.Info("Called Create with name \"" + act.Name + "\" and path \"" + act.Path + "\"")
	return t.templateService.Create(act)
}

func (t *Thop) DeleteSearch(act Search) (err error) {
	t.log.Debug("Called delete with search phrase \"" + act.Phrase + "\"")
	return
}

func (t *Thop) DeleteSelect() (err error) {
	t.log.Debug("Called delete with select")
	return
}

func (t *Thop) EditSearch(act Search) (err error) {
	t.log.Debug("Called edit with search phrase \"" + act.Phrase + "\"")
	return
}

func (t *Thop) EditSelect() (err error) {
	t.log.Debug("Called edit with select")
	return
}

func (t *Thop) KillSearch(act Search) (err error) {
	t.log.Debug("Called kill with search phrase \"" + act.Phrase + "\"")
	return
}

func (t *Thop) KillSelect() (err error) {
	t.log.Debug("Called kill with select")
	return
}

func (t *Thop) OpenSearch(act Search) (err error) {
	t.log.Debug("Called open with search phrase \"" + act.Phrase + "\"")
	return
}

func (t *Thop) OpenSelect() (err error) {
	t.log.Info("Called open with select")

	templates, err := t.templateService.List()
	if err != nil {
		panic(err)
	}

	var entries []*selector.Entry

	templateEntries := make(map[*selector.Entry]*template.Template)
	for _, template := range templates {
		entry := selector.NewEntry(template.Name(), string(template.FilePath()), selector.TagTemplate)
		templateEntries[entry] = template
		entries = append(entries, entry)
	}

	sessions, err := t.multiplexer.ListSessions()
	if err != nil {
		panic(err)
	}

	sessionEntries := make(map[*selector.Entry]*session.Session)
	for _, session := range sessions {
		entry := selector.NewEntry(session.Name(), session.Name(), selector.TagActiveSession)
		// TODO: Figure out a way to tag sessions when created from template to not duplicate entries here
		sessionEntries[entry] = session
		entries = append(entries, entry)
	}

	var result *selector.Entry
	if result, err = t.selector.SelectFrom(entries, selector.OperationOpen); err != nil {
		// TODO: Handling of cancellation
		panic(err)
	}

	if result == nil {
		return
	}

	switch result.Tag() {
	case selector.TagTemplate:
		if templ, ok := templateEntries[result]; ok {
			t.multiplexer.AttachTemplate(templ)
		}
	case selector.TagActiveSession:
		if session, ok := sessionEntries[result]; ok {
			t.multiplexer.AttachSession(session)
		}
	}

	return
}
