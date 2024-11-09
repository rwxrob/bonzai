/*
Package plain implements [pkg/github.com/rwxrob/bonzai/mark.Renderer]
*/
package plain

import (
	"io"
	"strings"

	"github.com/rwxrob/bonzai/mark"
)

var Renderer mark.Renderer = new(renderer)

type renderer struct{}

func (r *renderer) Render(this any, m *mark.Funcs, zmark io.Reader) (io.Reader, error) {
	buf, err := io.ReadAll(zmark)
	if err != nil {
		return nil, err
	}
	out, err := mark.Render(this, m, buf)
	return strings.NewReader(out), err
}
