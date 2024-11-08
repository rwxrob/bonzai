package mark

import (
	"io"
	"text/template"
)

// Renderer abstracts how a stream of BonzaiMark (zmark) is rendered to
// digital data whether it be text, HTML, PDF, or other binary data.
//
// To maximize compatibility between Renderers, implementations must
// only allow input that complies with the current BonzaiMark specification
// documented in [mark] package.
//
// See the following for examples of implementations:
//
//   - [pkg/github.com/rwxrob/bonzai/mark/renderers]
//   - [pkg/github.com/rwxrob/bonzai/cmds/help]
type Renderer interface {
	Render(this any, zmark io.Reader, tags *template.FuncMap) (io.Reader, error)
}
