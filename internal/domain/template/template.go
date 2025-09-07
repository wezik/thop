package template

type Version int

const (
	V1 Version = iota + 1
)

type Template struct {
	filePath    FilePath
	version     Version
	name        string
	sessionName string
	path        string
	commands    []string
	windows     []*Window
}

func DefaultTemplate(name, path string) *Template {
	return &Template{
		version: V1,
		name:    name,
		path:    path,
		windows: []*Window{
			{
				panes: []*Pane{{}},
			},
		},
	}
}

func NewTemplate(
	filePath FilePath,
	version Version,
	name string,
	sessionName string,
	path string,
	commands []string,
	windows []*Window,
) *Template {
	return &Template{
		filePath:    filePath,
		version:     version,
		name:        name,
		sessionName: sessionName,
		path:        path,
		commands:    commands,
		windows:     windows,
	}
}

func (t *Template) FilePath() FilePath {
	return t.filePath
}

func (t *Template) Version() Version {
	return t.version
}

func (t *Template) Name() string {
	return t.name
}

func (t *Template) SessionName() string {
	return t.sessionName
}

func (t *Template) Path() string {
	return t.path
}

func (t *Template) Commands() []string {
	return t.commands
}

func (t *Template) Windows() []*Window {
	return t.windows
}

type Window struct {
	name     string
	path     string
	layout   string
	commands []string
	panes    []*Pane
}

func NewWindow(
	name string,
	path string,
	layout string,
	commands []string,
	panes []*Pane,
) *Window {
	return &Window{
		name:     name,
		path:     path,
		layout:   layout,
		commands: commands,
		panes:    panes,
	}
}

func (w *Window) Name() string {
	return w.name
}

func (w *Window) Path() string {
	return w.path
}

func (w *Window) Layout() string {
	return w.layout
}

func (w *Window) Commands() []string {
	return w.commands
}

func (w *Window) Panes() []*Pane {
	return w.panes
}

type Pane struct {
	active   bool
	commands []string
	path     string
}

func NewPane(
	active bool,
	commands []string,
	path string,
) *Pane {
	return &Pane{
		active:   active,
		commands: commands,
		path:     path,
	}
}

func (p *Pane) Active() bool {
	return p.active
}

func (p *Pane) Commands() []string {
	return p.commands
}

func (p *Pane) Path() string {
	return p.path
}
