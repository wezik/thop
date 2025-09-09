package thop

import (
	"slices"
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

	sessions, err := t.multiplexer.ListSessions()
	if err != nil {
		panic(err)
	}

	// Quick lookup for templates since they are already loaded
	templateEntries := make(map[*selector.Entry]*template.Template)
	var entries []*selector.Entry

	for _, templ := range templates {
		// In case template session is already active, tag it as active template

		// To not iterate twice we check if deletion was successful, since we would need to remove it anyway
		lenBefore := len(sessions)
		sessions = slices.DeleteFunc(sessions, func(session *session.Session) bool {
			return templ.SessionName() == session.Name()
		})
		if lenBefore != len(sessions) {
			entry := selector.NewEntry(templ.Name(), templ.SessionName(), selector.TagActiveTemplate)
			entries = append(entries, entry)
			continue
		}

		// Otherwise create regular template entry
		entry := selector.NewEntry(templ.Name(), string(templ.FilePath()), selector.TagTemplate)
		entries = append(entries, entry)
		templateEntries[entry] = templ
	}

	// Append leftover sessions
	for _, session := range sessions {
		entry := selector.NewEntry(session.Name(), session.Name(), selector.TagActiveSession)
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
		} else {
			panic("Selected template was not added to the lookup map") // Should never happen
		}
	case selector.TagActiveSession, selector.TagActiveTemplate:
		sesh := session.NewSession(result.Name())
		t.multiplexer.AttachSession(sesh)
	}

	return
}
