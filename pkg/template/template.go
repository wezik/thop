package template

type Version int

type Template struct {
	FilePath    string    `yaml:"-"`
	Version     Version   `yaml:"version"`
	Name        string    `yaml:"name"`
	SessionName string    `yaml:"session_name,omitempty"`
	Path        string    `yaml:"path"`
	Commands    []string  `yaml:"run"`
	Windows     []*Window `yaml:"windows"`
}

type Window struct {
	Name     string   `yaml:"name"`
	Path     string   `yaml:"path"`
	Layout   string   `yaml:"layout"`
	Commands []string `yaml:"run"`
	Panes    []*Pane  `yaml:"panes"`
}

type Pane struct {
	Active   bool     `yaml:"active,omitempty"`
	Commands []string `yaml:"run"`
	Path     string   `yaml:"path"`
}
