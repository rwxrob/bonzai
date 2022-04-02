// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package help

import (
	"fmt"
	"unicode"

	"github.com/rwxrob/bonzai/comp"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/fn/filt"
	"github.com/rwxrob/fn/maps"
	"github.com/rwxrob/term"
	"github.com/rwxrob/to"
)

// Cmd provides help documentation for the caller allowing the specific
// section of help wanted to be passed as a tab-completable parameter.
var Cmd = &Z.Cmd{
	Name:    `help`,
	Summary: `display help similar to man page format`,
	Params: []string{
		"name", "title", "summary", "params", "commands", "description",
		"examples", "legal", "copyright", "license", "version",
	},
	Completer: helpCompleter,
	Call: func(x *Z.Cmd, args ...string) error {
		// TODO detect if local web help is preferred over terminal
		if len(args) == 0 {
			args = append(args, "all")
		}
		return nil
	},
}

func helpCompleter(x comp.Command, args ...string) []string {

	// not sure we've completed the command name itself yet
	if len(args) == 0 {
		return []string{x.GetName()}
	}

	// build list of visible commands and params
	list := []string{}
	list = append(list, x.GetParams()...)

	// if the caller has other sections get those
	caller := x.GetCaller()
	if caller != nil {
		other := caller.GetOther()
		if other != nil {
			list = append(list, maps.Keys(other)...)
		}
	}

	if len(args) == 0 {
		return list
	}

	return filt.HasPrefix(list, args[0])
}

// ForTerminal converts the collective help documentation of the given
// command into curses terminal-friendly output which minics UNIX man
// pages as much as possible. Documentation text is expected to be in
// standard BonzaiMark markup (see Format).
//
// If the special "all" section is passed all sections will be
// displayed.
//
// If the "less" pager is detected and the terminal is interactive
// (stdout is to a terminal/tty) will call Z.SysExec to transfer control
// to it and send the output to it making it virtual indistinguishable
// from "man" page output.
//
// If the terminal is non-interactive, simple prints as
// 80-column-wrapped plain text.
func ForTerminal(x *Z.Cmd, section string) {
	switch section {
	case "name":
		fmt.Println(x.Name)
	case "summary":
		fmt.Println(x.Summary)
	case "all":
		fmt.Println(Format("**NAME**"))
		fmt.Println(to.IndentWrapped(x.Title(), 7, 80) + "\n")
		fmt.Println(Format("**SYNOPSIS**"))
		// if commands [COMMAND]
		// if commands and parameters [COMMAND|PARAM]
		// if just parameters [PARAM]
		//fmt.Println(to.IndentWrapped(
		fmt.Println(Format("**PARAMS**"))
		fmt.Println(Format("**COMMANDS**"))
		//Format("**"+x.Name+"** ")+x.(), 7, 80) + "\n")
		fmt.Println(Format("**DESCRIPTION**"))
		//ForTerminal(x, "synopsis")
	default:
		if v, has := x.Other[section]; has {
			fmt.Println(Format(v))
		}
	}
	//"title", "summary", "params", "commands", "description",
	//"examples", "legal", "copyright", "license", "version",
}

// Format is called by Render when producing terminal formatted output
// containing all the BonzaiMark help documentation for a given Bonzai
// command. BonzaiMark is a very minimal subset of CommonMark.
//
// Initial and trailing blank lines and lines with only whitespace will
// be stripped as well as the initial number of spaces for the first
// line of indentation. This allows markup strings to be written in very
// readable ways even when embedded within source code (preferably with
// backtick string literal notation).
//
// Any line beginning with at least four spaces (after trimming
// indentation) will be kept verbatim.
//
// Emphasis will be applied as possible if the following markup is
// detected:
//
//     *italic*
//     **bold**
//     ***bolditalic***
//     <under> (brackets remain)
//
// Note that the format of the emphasis might not always be as
// specifically named. For example, most terminal do not support italic
// fonts and so will instead underline *italic* text, so (as specified
// in HTML5 for <i>, for example) these format names should be taken to
// mean their semantic equivalents.
//
// For historic reasons, the following environment variables will be
// detected and used to replace the specified escapes (see rwxrob/term
// package for details):
//
//     LESS_TERMCAP_mb
//     LESS_TERMCAP_md
//     LESS_TERMCAP_me
//     LESS_TERMCAP_se
//     LESS_TERMCAP_so
//     LESS_TERMCAP_ue
//     LESS_TERMCAP_us
//
func Format(markup string) string {
	out := to.Dedented(markup)
	out = to.Wrapped(out, 80)
	out = Emphasize(out)
	return out
}

// Emphasize replaces minimal Markdown-like syntax with *Italic*,
// **Bold**, ***BoldItalic***, and <bracketed>
func Emphasize(buf string) string {

	// italic = `<italic>`
	// bold = `<bold>`
	// bolditalic = `<bolditalic>`
	// reset = `<reset>`

	term.EmphFromLess()

	nbuf := []rune{}
	prev := ' '
	opentok := false
	otok := ""
	closetok := false
	ctok := ""
	for i := 0; i < len([]rune(buf)); i++ {
		r := []rune(buf)[i]

		if r == '<' {
			nbuf = append(nbuf, '<')
			nbuf = append(nbuf, []rune(term.Under)...)
			for {
				i++
				r = rune(buf[i])
				if r == '>' {
					i++
					break
				}
				nbuf = append(nbuf, r)
			}
			nbuf = append(nbuf, []rune(term.Reset)...)
			nbuf = append(nbuf, '>')
			i--
			continue
		}

		if r != '*' {

			if opentok {
				tokval := " "
				if !unicode.IsSpace(r) {
					switch otok {
					case "*":
						tokval = term.Italic
					case "**":
						tokval = term.Bold
					case "***":
						tokval = term.BoldItalic
					}
				} else {
					tokval = otok
				}
				nbuf = append(nbuf, []rune(tokval)...)
				opentok = false
				otok = ""
			}

			if closetok {
				nbuf = append(nbuf, []rune(term.Reset)...) // practical, not perfect
				ctok = ""
				closetok = false
			}

			prev = r
			nbuf = append(nbuf, r)
			continue
		}

		// everything else for '*'
		if unicode.IsSpace(prev) || opentok {
			opentok = true
			otok += string(r)
			continue
		}

		// only closer conditions remain
		if !unicode.IsSpace(prev) {
			closetok = true
			ctok += string(r)
			continue
		}

		// nothing special
		closetok = false
		nbuf = append(nbuf, r)
	}

	// for tokens at the end of a block
	if closetok {
		nbuf = append(nbuf, []rune(term.Reset)...)
	}

	return string(nbuf)
}
