package thop

import (
	"strconv"
	"thop/internal/domain/log"
	"thop/internal/domain/multiplexer"
	"thop/internal/domain/platform"
	"thop/internal/domain/selector"
	"thop/internal/domain/template"
	"time"
)

type CreateTemplate struct {
	Name string
	Path string
}

type Search struct {
	Phrase string
}

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

func (t *Thop) Create(act CreateTemplate) (err error) {
	t.log.Info("Called Create with name \"" + act.Name + "\" and path \"" + act.Path + "\"")
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

	templ := template.DefaultTemplate(name, path)
	err = t.templateFileStorage.Save(templ)

	return
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

	start := time.Now()

	templateFiles, err := t.templateFileStorage.List()
	if err != nil {
		panic(err)
	}

	t.log.Debug("Discovered " + strconv.Itoa(len(templateFiles)) + " template files in the tree in " + time.Since(start).String())

	start = time.Now()

	// load all templates (necessary to get active sessions and names)
	templates := make(map[string]*template.Template)
	for _, templateFile := range templateFiles {
		templ, err := t.templateFileStorage.LoadTemplate(templateFile.Path())
		if err != nil {
			// TODO: Collect failed templates and notify? Otherwise just ignore
			continue
		}

		templates[string(templateFile.Path())] = templ
	}

	t.log.Debug("Loaded " + strconv.Itoa(len(templates)) + " templates in " + time.Since(start).String())

	// TODO: Add pulling of active sessions from multiplexer and in addition to that
	//       match active sessions with template files to not duplicate entries

	var entries []*selector.Entry

	for key, templ := range templates {
		if templ != nil {
			entry := selector.NewEntry(templ.Name(), key, selector.TagTemplate)
			entries = append(entries, entry)
		}
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
		if templ, ok := templates[result.Key()]; ok {
			t.multiplexer.AttachTemplate(templ)
		}
	}

	return
}
