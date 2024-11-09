/*
Package plain implements [pkg/github.com/rwxrob/bonzai/mark.Renderer]
*/
package plain

import (
	"io"
	"text/template"

	"github.com/rwxrob/bonzai/mark"
)

var Renderer mark.Renderer = new(renderer)

type renderer struct{}

func (r *renderer) Render(this any, m *mark.Map, zmark io.Reader) (io.Reader, error) {
	tmpl, err := template.New("greeting").Parse(templateStr)
	if err != nil {
		panic(err)
	}

	return mark.Render(this, m, zmark)
}
