package mark

import "io"

type Renderer interface {
	Render(in io.Reader) error
}

type Viewer interface {
	View(in io.Reader, r Renderer) error
}
