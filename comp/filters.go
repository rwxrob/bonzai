package comp

import (
	"github.com/rwxrob/bonzai"
)

type Filter interface {
	Filter([]string) []string
}

type Pipe []any

func (pipe Pipe) Complete(args ...string) []string {
	list := make([]string, 0)
	for _, elem := range pipe {
		switch typed := elem.(type) {
		case Filter:
			list = typed.Filter(list)
		case bonzai.Completer:
			list = append(list, typed.Complete(args...)...)
		}
	}
	return list
}
