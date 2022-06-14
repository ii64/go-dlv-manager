package server

import "fmt"

type Option struct {
	Program     string
	ProgramArgs []string

	BufferSize int
}

func (o *Option) Validate() error {
	if o.Program == "" {
		return fmt.Errorf("program is required")
	}
	if o.BufferSize < 1 {
		o.BufferSize = 4096
	}
	return nil
}
