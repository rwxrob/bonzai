package futil

import "fmt"

type ErrorExists struct {
	P string
}

func (e ErrorExists) Error() string {
	return fmt.Sprintf("error: %v exists", e.P)
}
