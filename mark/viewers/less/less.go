package less

import (
	"io"

	"github.com/rwxrob/bonzai/mark"
	"github.com/rwxrob/bonzai/term"
)

type Viewer struct{}

var _ mark.Viewer = new(Viewer)

func (v *Viewer) View(in io.Reader, r mark.Renderer) error {
	// TODO detect less pager and use it or just print it.
	out, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	term.Print(out)
	return nil
}
