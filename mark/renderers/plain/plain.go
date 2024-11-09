/*
Package plain implements [pkg/github.com/rwxrob/bonzai/mark.Renderer] by just calling [pkt/github.com/rwxrob/bonzai/mark.Render] on the zmark input as a string. Only template replacement is attempted. No additional formatting is applied.
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
	out, err := mark.Render(this, m, string(buf))
	return strings.NewReader(out), err

}
