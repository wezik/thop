package thop

import "thop/pkg/command"

type Thop struct {
}

func New() *Thop {
	return &Thop{}
}

func (t *Thop) Create(createTemplate command.CreateTemplate) error {
	return nil
}
