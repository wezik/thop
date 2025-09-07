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
	*selector.Entry
}

func (s *FzfSelector) SelectFrom(
	entries []*selector.Entry,
	operation selector.Operation,
) (result *selector.Entry, err error) {

	var records []record

	for _, entry := range entries {
		record := toRecord(entry)
		records = append(records, record)
	}

	slices.SortFunc(records, sortByTag)

	recordMap := make(map[string]*selector.Entry)
	var input bytes.Buffer

	for _, record := range records {
		name := record.Name()
		if record.Tag() == selector.TagActiveSession || record.Tag() == selector.TagActiveTemplate {
			name = "(Active) " + name
		}
		name += "\n" // newline character separates entries in fzf

		// avoid duplicates, they should not be allowed into the system, and if they appear ignore them
		if recordMap[name] != nil {
			continue
		}

		recordMap[name] = record.Entry
		input.WriteString(name)
	}

	cmd := exec.Command("fzf")
	cmd.Stdin = &input
	cmd.Args = append(cmd.Args, "--prompt", resolveOperation(operation))

	out, exitCode, err := s.exec(cmd)
	if exitCode == 130 {
		s.log.Info("Selector cancelled")
		err = nil // for now just ignore the error
		return
	} else if err != nil {
		s.log.Warn("Selector failed with error \"" + err.Error() + "\"")
		s.log.Error(err)
		return
	}

	result, ok := recordMap[out]
	if !ok {
		s.log.Debug("Selector result not found")
		return
	}

	s.log.Info("Selector selected \"" + result.Name() + "\"")
	s.log.Debug("Selector selected key \"" + result.Key() + "\"")

	return
}

func toRecord(entry *selector.Entry) record {
	return record{
		Entry: entry,
	}
}

func getTagSortOrder(tag selector.Tag) int {
	switch tag {

	case selector.TagActiveTemplate:
		return 0

	case selector.TagActiveSession:
		return 0

	case selector.TagTemplate:
		return 1

	default:
		return 99
	}
}

func sortByTag(a, b record) int {
	aTag := getTagSortOrder(a.Tag())
	bTag := getTagSortOrder(b.Tag())
	if aTag != bTag {
		// sort ascending by tag first
		return bTag - aTag
	}

	aName := strings.ToLower(a.Name())
	bName := strings.ToLower(b.Name())

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
