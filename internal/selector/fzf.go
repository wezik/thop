package selector

import (
	"bytes"
	"slices"
	"strings"
	"thop/pkg/log"
	"thop/pkg/selector"
)

type FzfSelector struct {
	log log.Logger
}

func NewFzfSelector(log log.Logger) selector.Selector {
	return &FzfSelector{
		log: log,
	}
}

type record struct {
	selector.Entry
	tag         tag
	displayName string
}

type tag int

const (
	tagActive tag = iota
	tagNone
)

func (s *FzfSelector) SelectFrom(
	entries []selector.Entry,
	operation selector.Operation,
) (result selector.Entry, err error) {

	var records []record

	for _, entry := range entries {
		record := toRecord(entry)
		records = append(records, record)
	}

	slices.SortFunc(records, sortByTag)

	recordMap := make(map[string]selector.Entry)
	var input bytes.Buffer

	for _, record := range records {
		name := record.displayName
		if record.tag == tagActive {
			name = name + " (Active)"
		}
		name += "\n"

		recordMap[name] = record.Entry
		input.WriteString(name)
	}

	println(input.String())

	// TODO: Actually call fzf via cmd.Exec here
	return
}

func format(s string) string {
	return strings.ReplaceAll(s, "_", " ")
}

func toRecord(entry selector.Entry) record {
	name := format(entry.EntryName())
	tag := tagNone

	switch entry := entry.(type) {
	case *selector.SessionEntry:
		tag = tagActive

	case *selector.TemplateEntry:
		if entry.IsActive {
			tag = tagActive
		}
	}

	return record{
		Entry:       entry,
		tag:         tag,
		displayName: name,
	}
}

func sortByTag(a, b record) int {
	if a.tag != b.tag {
		// sort ascending by tag first
		return int(b.tag) - int(a.tag)
	}

	aName := strings.ToLower(a.displayName)
	bName := strings.ToLower(b.displayName)

	// sort case-insensitive ascending
	return strings.Compare(bName, aName)
}
