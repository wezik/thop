package thop

import (
	"thop/pkg/action"
	"thop/pkg/log"
	"thop/pkg/multiplexer"
	"thop/pkg/platform"
	"thop/pkg/selector"
	"thop/pkg/template"
)

type Thop struct {
	log                 log.Logger
	selector            selector.Selector
	templateFileStorage template.FileStorage
	multiplexer         multiplexer.Multiplexer
	getwd               platform.GetwdFn
}

func New(
	log log.Logger,
	selector selector.Selector,
	templateFileStorage template.FileStorage,
	multiplexer multiplexer.Multiplexer,
	getwd platform.GetwdFn,
) *Thop {
	return &Thop{
		log:                 log,
		selector:            selector,
		templateFileStorage: templateFileStorage,
		multiplexer:         multiplexer,
		getwd:               getwd,
	}
}

func (t *Thop) Create(act action.CreateTemplate) (err error) {
	t.log.Debug("Called Create with name \"" + act.Name + "\" and path \"" + act.Path + "\"")
	path := act.Path
	if path == "" {
		if path, err = t.getwd(); err != nil {
			return
		}
	}

	name := act.Name
	if name == "" {
		name = path
	}

	return
}

func (t *Thop) DeleteSearch(act action.Search) (err error) {
	t.log.Debug("Called DeleteSearch with phrase \"" + act.Phrase + "\"")
	return
}

func (t *Thop) DeleteSelect() (err error) {
	t.log.Debug("Called DeleteSelect")
	return
}

func (t *Thop) EditSearch(act action.Search) (err error) {
	t.log.Debug("Called EditSearch with phrase \"" + act.Phrase + "\"")
	return
}

func (t *Thop) EditSelect() (err error) {
	t.log.Debug("Called EditSelect")
	return
}

func (t *Thop) KillSearch(act action.Search) (err error) {
	t.log.Debug("Called KillSearch with phrase \"" + act.Phrase + "\"")
	return
}

func (t *Thop) KillSelect() (err error) {
	t.log.Debug("Called KillSelect")
	return
}

func (t *Thop) OpenSearch(act action.Search) (err error) {
	t.log.Debug("Called OpenSearch with phrase \"" + act.Phrase + "\"")
	return
}

func (t *Thop) OpenSelect() (err error) {
	t.log.Debug("Called OpenSelect")

	var entries []selector.Entry

	templateFiles, err := t.templateFileStorage.List()
	if err != nil {
		panic(err)
	}

	// TODO: Add pulling of active sessions from multiplexer and in addition to that
	//       match active sessions with template files to not duplicate entries

	for _, templateFile := range templateFiles {
		entries = append(entries, &selector.TemplateEntry{
			File: templateFile,
			IsActive: false,
		})
	}

	var entry selector.Entry
	if entry, err = t.selector.SelectFrom(entries, selector.OperationOpen); err != nil {
		panic(err)
	}

	switch entry := entry.(type) {
	case *selector.TemplateEntry:
		template := t.templateFileStorage.LoadTemplate(entry.File.Path())
		t.multiplexer.AttachTemplate(template)

	case *selector.SessionEntry:
		t.multiplexer.AttachSession(entry.Session)
	}

	return
}
