package comp

import (
	"github.com/rwxrob/bonzai"
)

type _WithMaxLen struct {
	len   int
	comps Combine
}

// WithMaxLen return a [bonzai.Completer] that limits the length of
// completions to a maximum length.
func WithMaxLen(len int, comps ...bonzai.Completer) _WithMaxLen {
	return _WithMaxLen{len, comps}
}

func (f _WithMaxLen) Complete(x bonzai.Cmd, args ...string) []string {
	var list []string
	for _, completion := range f.comps.Complete(x, args...) {
		if len(completion) <= f.len {
			list = append(list, completion)
		}
	}
	return list
}

type _WithMinLen struct {
	len   int
	comps Combine
}

// WithMinLen return a [bonzai.Completer] that limits the length of
// completions to a minimum length.
func WithMinLen(len int, comps ...bonzai.Completer) _WithMinLen {
	return _WithMinLen{len, comps}
}

func (f _WithMinLen) Complete(x bonzai.Cmd, args ...string) []string {
	var list []string
	for _, completion := range f.comps.Complete(x, args...) {
		if len(completion) >= f.len {
			list = append(list, completion)
		}
	}
	return list
}

type _WithPrefix struct {
	prefix string
	comps  Combine
}

// WithPrefix return a [bonzai.Completer] that filters completions that
// have the given prefix. This removes the requirement for the user to
// type the prefix to get complete completions.
func WithPrefix(prefix string, comps ...bonzai.Completer) _WithPrefix {
	return _WithPrefix{prefix, comps}
}

func (f _WithPrefix) Complete(x bonzai.Cmd, args ...string) []string {
	if len(args) == 0 {
		return f.comps.Complete(x, f.prefix)
	}
	args[0] = f.prefix + args[0]
	return f.comps.Complete(x, args...)
}
