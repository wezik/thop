package template

type YamlTemplate struct {
	Version     int           `yaml:"version"`
	Name        string        `yaml:"name"`
	SessionName string        `yaml:"session_name,omitempty"`
	Path        string        `yaml:"path"`
	Commands    []string      `yaml:"run,omitempty"`
	Windows     []*YamlWindow `yaml:"windows,omitempty"`
}

type YamlWindow struct {
	Name     string      `yaml:"name,omitempty"`
	Path     string      `yaml:"path,omitempty"`
	Layout   string      `yaml:"layout,omitempty"`
	Commands []string    `yaml:"run,omitempty"`
	Panes    []*YamlPane `yaml:"panes,omitempty"`
}

type YamlPane struct {
	Active   bool     `yaml:"active,omitempty"`
	Commands []string `yaml:"run,omitempty"`
	Path     string   `yaml:"path,omitempty"`
}
