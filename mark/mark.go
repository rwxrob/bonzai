package mark

import (
	"io"
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
// Markdown like [Pandoc]. Callers can override the {{cmdtree . }}
// function with their own if more detailed tree representations are
// wanted, but doing so risks breaking renderers that expect the command
// tree to be in the default format.
//
// [Pandoc]: https://pandoc.org/
func Usage(x *bonzai.Cmd) (string, error) {
	out := new(strings.Builder)
	out.WriteString("# Usage\n\n")
	out.WriteString("{{cmdtree .}}\n")
	if len(x.Long) > 0 {
		out.WriteString("\n" + to.Dedented(x.Long))
		if x.Long[len(x.Long)-1] != '\n' {
			out.WriteString("\n")
		}
	}
	f := to.MergedMaps(funcs.Map, x.Funcs)
	return Fill(x, f, out.String())
}

// Fill processes the input string (in) as a [pkg/text/template] using
// the provided function map (f) and the data context (it). It returns
// the rendered output as a string or an error if any step fails. No
// functions beyond those passed are merged (unlike [Usage]).
func Fill(it any, f template.FuncMap, in string) (string, error) {
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
