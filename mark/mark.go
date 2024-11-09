package mark

import (
	"io"
	"strings"
	"text/template"
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
// See the following for examples of implementations:
//
//   - [pkg/github.com/rwxrob/bonzai/mark/renderers]
//   - [pkg/github.com/rwxrob/bonzai/cmds/help]
type Renderer interface {
	Render(this any, m *Funcs, zmark io.Reader) (io.Reader, error)
}

type Funcs template.FuncMap

// Render simplifies rendering a [text/template] by processes the input
// (in) and using the provided [Funcs] [template.FuncMap] and the data
// context (this). It returns the rendered output as a string or an
// error if any step fails. The input is not required to be in
// BonzaiMark format (unlike implementations of [Renderer]).
func Render(this any, m *Funcs, in string) (string, error) {
	var fm template.FuncMap
	if m != nil {
		fm = template.FuncMap(*m)
	}
	tmpl, err := template.New("t").Funcs(fm).Parse(in)
	if err != nil {
		return "", err
	}
	out := new(strings.Builder)
	if err := tmpl.Execute(out, this); err != nil {
		return "", err
	}
	return out.String(), nil
}

// MustRender calls [Render] but panics if an error occurs.
func MustRender(this any, m *Funcs, in string) string {
	out, err := Render(this, m, in)
	if err != nil {
		panic(err)
	}
	return out
}
