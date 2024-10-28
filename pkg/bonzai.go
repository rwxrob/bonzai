package bonzai

import (
	"log"

	"github.com/rwxrob/bonzai/pkg/core/vars"
)

var Vars vars.Driver

func init() {
	m := vars.NewMap()
	if err := m.Init(); err != nil {
		log.Print(err)
		return
	}
	Vars = m
}

const (
	FAILED   = -1
	NOTFOUND = 0
	SUCCESS  = 1
)

// Completer specifies a struct with a [Completer.Complete] function
// that will complete the first argument (usually a command of some kind)
// based on the remaining arguments. The [Complete] method must never
// panic and always return at least an empty slice of strings. By
// convention Completers that do not make use of or other arguments
// should use an underscore identifier since they are ignored.
type Completer interface {
	Complete(x any, args ...string) []string
}
