package schedule

import (
	"time"

	"github.com/mlange-42/track/core"
)

// Renderer is an interface for schedule renderers
type Renderer interface {
	Render(t *core.Track, reporter *core.Reporter, startDate time.Time) (string, error)
}
