package mark

import (
	"fmt"
	"io"
	"slices"
	"strings"
	"text/template"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/mark/funcs"
	"github.com/rwxrob/bonzai/to"
)

// Renderer abstracts how a stream of BonzaiMark (zmark) is rendered to
// digital data whether it be text, HTML, PDF, or other binary data.
//
// To maximize compatibility between Renderers, implementations must
// only allow input that complies with the current BonzaiMark
// specification documented in [mark] package. Implementations may
// extend that specification and support more complex markups but
// developers must understand such specialization will be much less
// useful to as many people.
//
// # Templates not included
//
// Note that although BonzaiMark will often be generated from Go
// [pkg/text/template] templates (such as is allowed in Long) that the
// template itself is never a part of the specification, even though
// someone commands (like {{code "go mod init"}}) are designed to help
// with the generation of Markdown from within Go strings.
//
// # Renderers as viewers
//
// Renderers are not intended to fire off a viewer instead leaving that
// work to the caller. Renderers can, however, have very specific ideas
// about how the output will be rendered (ANSI escapes, HTML, etc.).
//
// # Reference implementations and examples
//
//   - [pkg/github.com/rwxrob/bonzai/mark/renderers]
//   - [pkg/github.com/rwxrob/bonzai/cmds/help]
type Renderer interface {
	Render(zmark io.Reader) (io.Reader, error)
}

type HTMLRenderer interface {
	RenderHTML(zmark io.Reader) (io.Reader, error)
}

type ANSIRenderer interface {
	RenderANSI(zmark io.Reader) (io.Reader, error)
}

type TextRenderer interface {
	RenderText(zmark io.Reader) (io.Reader, error)
}

type ManRenderer interface {
	RenderMan(zmark io.Reader) (io.Reader, error)
}

type PDFRenderer interface {
	RenderPDF(zmark io.Reader) (io.Reader, error)
}

// Usage outputs a Markdown view of a Cmd from
// [pkg/github.com/rwxrob/bonzai] package filling the Cmd.Long by
// rendering it as a [pkg/text/template] using itself as the object and
// merging the Cmd.Funcs over [pkg/github.com/rwxrob/mark/funcs].Map to
// provide the [pkg/text/template.Funcs]. This Markdown can be passed to
// any [Renderer] but can also be piped directly to tools that support
// Markdown like [Pandoc].
//
// [Pandoc]: https://pandoc.org/
func Usage(x *bonzai.Cmd) (io.Reader, error) {
	out := new(strings.Builder)
	out.WriteString("# Usage\n\n")
	out.WriteString(CmdTree(x))
	if len(x.Long) > 0 {
		out.WriteString("\n" + to.Dedented(x.Long))
		if x.Long[len(x.Long)-1] != '\n' {
			out.WriteString("\n")
		}
	}
	f := to.MergedMaps(funcs.Map, x.Funcs)
	str, err := Render(x, f, out.String())
	if err != nil {
		return nil, err
	}

	return strings.NewReader(str), nil
}

// UsageString reads input from [Usage] [io.Reader] function associated
// with the command and returns it as a string. It uses
// a [strings.Builder] to efficiently build the output string and
// ignores any errors.
func UsageString(x *bonzai.Cmd) (string, error) {
	var buf strings.Builder
	r, err := Usage(x)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(&buf, r)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func isLastOf(this, caller *bonzai.Cmd) bool {
	// default is always last if not also in command list
	if caller.Def != nil &&
		this == caller.Def &&
		!slices.Contains(caller.Cmds, this) {
		return true
	}
	l := len(caller.Cmds)
	return l > 0 && this == caller.Cmds[l-1]
}

func cmdTree(x *bonzai.Cmd, depth int) string {
	out := new(strings.Builder)
	addbranch := func(c *bonzai.Cmd) error {
		caller := c.Caller()
		clevel := c.Level()
		xlevel := x.Level()
		if clevel-xlevel > depth {
			return nil
		}
		for range clevel - 1 {
			out.WriteString("│ ")
		}
		if clevel > 0 {
			switch {
			case isLastOf(c, caller):
				out.WriteString("└─")
			default:
				out.WriteString("├─")
			}
		}
		if c.IsHidden() {
			out.WriteString("(hidden) ← contains hidden subcommands\n")
			return nil
		}
		name := c.Name
		if len(name) == 0 {
			name = `noname`
		}
		out.WriteString(name)
		if len(c.Short) > 0 {
			out.WriteString(" ← " + c.Short)
		}
		if caller != nil && caller.Def == c {
			out.WriteString(" (default)")
		}
		out.WriteString("\n")
		return nil
	}
	x.WalkDeep(addbranch, nil)
	return out.String()
}

// CmdTreeString generates and returns a formatted string representation
// of the command tree for the [Cmd] instance and all its [Cmd].Cmds
// subcommands. It aligns [Cmd].Short summaries in the output for better
// readability, adjusting spaces based on the position of the dashes.
func CmdTree(x *bonzai.Cmd) string {
	tree := cmdTree(x, 2)
	lines := to.Lines(tree)
	var widest int
	for _, line := range lines {
		if length := strings.IndexRune(line, '←'); length > widest {
			widest = length
		}
	}
	for i, line := range lines {
		parts := strings.Split(line, "←")
		if len(parts) > 1 {
			lines[i] = fmt.Sprintf("    %-*v←%v", widest-6, parts[0], parts[1])
		}
	}
	return strings.Join(lines, "\n") + "\n"
}

// Render processes the input string (in) as a template using the provided
// function map (f) and the data context (it). It returns the rendered
// output as a string or an error if any step fails. No functions beyond
// those passed are merged (unlike [Usage]).
func Render(it any, f template.FuncMap, in string) (string, error) {
	tmpl, err := template.New("t").Funcs(f).Parse(in)
	if err != nil {
		return "", err
	}
	out := new(strings.Builder)
	if err := tmpl.Execute(out, it); err != nil {
		return "", err
	}
	return out.String(), nil
}
