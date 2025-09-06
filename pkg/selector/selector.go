package selector

import (
	"thop/pkg/session"
	"thop/pkg/template"
)

type Entry interface {
	EntryName() string
}

type TemplateEntry struct {
	*template.File
	IsActive bool
}

func (t *TemplateEntry) EntryName() string {
	return t.Name()
}

type SessionEntry struct {
	*session.Session
}

func (s *SessionEntry) EntryName() string {
	return s.Name()
}

type Operation int

const (
	OperationOpen Operation = iota
	OperationEdit
	OperationDelete
	OperationKill
)

type Selector interface {
	SelectFrom([]Entry, Operation) (Entry, error)
}
