package mark

import (
	"io"
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
// # First fill template input
//
// All implementations must first fill the incoming data using the
// functions passed in f template.FuncMap (if any). The string returned
// from any function in template.FuncMap must therefore only return
// Markdown (preferably BonzaiMark for best compatibility). The filled
// template data is then rendered using whatever method implemented by
// the Renderer. For example, {{code "foo"}} returns `foo` instead of
// something else like <code>foo</code> or an equivalent with ANSI
// escapes.
//
// # Reference implementations and examples
//
//   - [pkg/github.com/rwxrob/bonzai/mark/renderers]
//   - [pkg/github.com/rwxrob/bonzai/cmds/help]
type Renderer interface {
	Render(this any, f template.FuncMap, zmark io.Reader) (io.Reader, error)
}
