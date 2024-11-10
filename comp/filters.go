package comp

import (
	"github.com/rwxrob/bonzai"
)

type withMaxLen struct {
	len       int
	completer bonzai.Completer
}

// WithMaxLen return a [bonzai.Completer] that limits the length of
// completions to a maximum length.
func WithMaxLen(len int, completer bonzai.Completer) withMaxLen {
	return withMaxLen{len, completer}
}

func (f withMaxLen) Complete(x bonzai.Cmd, args ...string) []string {
	var list []string
	for _, completion := range f.completer.Complete(x, args...) {
		if len(completion) <= f.len {
			list = append(list, completion)
		}
	}
	return list
}

type withMinLen struct {
	len       int
	completer bonzai.Completer
}

// WithMinLen return a [bonzai.Completer] that limits the length of
// completions to a minimum length.
func WithMinLen(len int, completer bonzai.Completer) withMinLen {
	return withMinLen{len, completer}
}

func (f withMinLen) Complete(x bonzai.Cmd, args ...string) []string {
	var list []string
	for _, completion := range f.completer.Complete(x, args...) {
		if len(completion) >= f.len {
			list = append(list, completion)
		}
	}
	return list
}

type withPrefix struct {
	prefix    string
	compelter bonzai.Completer
}

// WithPrefix return a [bonzai.Completer] that filters completions that
// have the given prefix. This removes the requirement for the user to
// type the prefix to get complete completions.
func WithPrefix(prefix string, completer bonzai.Completer) withPrefix {
	return withPrefix{prefix, completer}
}

func (f withPrefix) Complete(x bonzai.Cmd, args ...string) []string {
	if len(args) == 0 {
		return f.compelter.Complete(x, f.prefix)
	}
	args[0] = f.prefix + args[0]
	return f.compelter.Complete(x, args...)
}
