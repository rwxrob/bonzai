package less

import (
	"io"

	"github.com/rwxrob/bonzai/mark"
	"github.com/rwxrob/bonzai/term"
)

type viewer struct{}

var Viewer mark.Viewer = new(viewer)

func (v *viewer) View(in io.Reader, r mark.Renderer) error {
	// TODO detect less pager and use it or just print it.
	out, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	term.Print(out)
	return nil
}
