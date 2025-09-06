package selector

import (
	"bytes"
	"os/exec"
	"slices"
	"strings"
	"thop/pkg/log"
	"thop/pkg/platform"
	"thop/pkg/selector"
)

type FzfSelector struct {
	log  log.Logger
	exec platform.ExecFn
}

func NewFzfSelector(
	log log.Logger,
	exec platform.ExecFn,
) selector.Selector {
	return &FzfSelector{
		log:  log,
		exec: exec,
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

	cmd := exec.Command("fzf")
	cmd.Stdin = &input
	cmd.Args = append(cmd.Args, "--prompt", resolveOperation(operation))

	out, exitCode, err := s.exec(cmd)
	if exitCode == 130 {
		s.log.Debug("Selector cancelled")
		return
	} else if err != nil {
		s.log.Debug("Selector failed with error \"" + err.Error() + "\"")
		return
	}

	result, ok := recordMap[out]
	if !ok {
		s.log.Debug("Selector result not found")
		return
	}

	s.log.Debug("Selector selected \"" + result.EntryName() + "\"")

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

func resolveOperation(operation selector.Operation) string {
	switch operation {
	case selector.OperationOpen:
		return "Open > "

	case selector.OperationEdit:
		return "Edit > "

	case selector.OperationDelete:
		return "Delete > "

	case selector.OperationKill:
		return "Kill > "

	default:
		panic("unhandled operation")
	}
}
