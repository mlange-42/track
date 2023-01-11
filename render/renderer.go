package render

import (
	"io"
)

// Renderer is an interface for rendering reports
type Renderer interface {
	Render(io.Writer) error
}
