package project

import (
	"thop/internal/types"
	"thop/internal/types/template"
)

type UUID string
type Name string

type ProjectType int

const (
	TypeTemplate    ProjectType = iota // will default to TypeTemplate if not set explicitly
	TypeTmuxSession                    // signals tmux session
	TypeInvalid                        // signals invalid project template
)

type Project struct {
	UUID     UUID              `yaml:"-"`
	Name     Name              `yaml:"name"`
	Version  types.Version     `yaml:"version"`
	Template template.Template `yaml:"template"`
	Type     ProjectType       `yaml:"-"`
}

func (p *Project) WithDefaults() Project {
	newProject := *p
	newProject.Template = newProject.Template.WithDefaults()
	return newProject
}

func (p *Project) IsValid() bool {
	// for now other types shouldn't be treated as valid
	if p.Type != TypeTemplate {
		return false
	}

	if p.UUID == "" {
		return false
	}

	if p.Name == "" {
		return false
	}

	return p.Template.IsValid()
}
