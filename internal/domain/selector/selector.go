package selector

type Tag string

const (
	TagTemplate       Tag = "template"
	TagActiveTemplate Tag = "active-template"
	TagActiveSession  Tag = "active-session"
)

type Entry struct {
	name string
	key  string
	tag  Tag
}

func NewEntry(name string, key string, tag Tag) *Entry {
	return &Entry{
		name: name,
		key:  key,
		tag:  tag,
	}
}

func (e *Entry) Name() string {
	return e.name
}

func (e *Entry) Key() string {
	return e.key
}

func (e *Entry) Tag() Tag {
	return e.tag
}

type Operation int

const (
	OperationOpen Operation = iota
	OperationEdit
	OperationDelete
	OperationKill
)

type Selector interface {
	SelectFrom([]*Entry, Operation) (*Entry, error)
}
