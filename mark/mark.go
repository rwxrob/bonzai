package mark

import (
	"io"
	"text/template"
)

// Renderer abstracts how a stream of BonzaiMark is rendered to digital
// data whether it be text, HTML, PDF, or other binary data.
// Implementations must only allow compliant BonzaiMark input that
// complies with the Specification (documented in [mark] package
// documentation) and must implement the [FuncMap] template function
// tags. Renderers should not do things better suited for a [Viewer].
// See the collection of [pkg/github.com/rwxrob/bonzai/mark/renderers]
// and the [pkg/github.com/rwxrob/bonzai/cmds/help] for example usage.
type Renderer interface {
	Render(in io.Reader, tag template.FuncMap) (io.Reader, error)
}

// Viewer abstracts how a stream of BonzaiMark is viewed.
// Implementations must only allow compliant BonzaiMark input that
// complies with the [Specification]. Viewers should not do things
// better done in a Renderer. Implementations will vary from those
// that open a local web browser and serve up HTML to those that simply
// print hard-wrapped, plain, white text to the screen (as with Go doc
// output). For example, a Viewer that pages through text on the
// terminal may wish to render that text as [Charmbracelet Glamour] or
// simply as plain text. Do not pass a Renderer that is incompatible
// with a give Viewer, for example, a color ANSI terminal byte stream to
// a web viewer. See the collection of
// [pkg/github.com/rwxrob/bonzai/mark/viewers] and the
// [pkg/github.com/rwxrob/bonzai/cmds/help] for example usage.
//
// [Charmbracelet Glamour]: https://github.com/charmbracelet/glamour
type Viewer interface {
	View(in io.Reader, r Renderer) error
}
