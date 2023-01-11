package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/stretchr/testify/assert"
)

func TestConfirm(t *testing.T) {
	no := "n"
	yes := "y"

	out.StdIn = strings.NewReader(no)
	assert.False(t, confirm("", yes), "Should not be confirmed")

	out.StdIn = strings.NewReader(yes)
	assert.True(t, confirm("", yes), "Should be confirmed")
}

func setupTestCommand() (*core.Track, error) {
	dir, err := os.MkdirTemp("", "track-test")
	if err != nil {
		return nil, err
	}

	track, err := core.NewTrack(&dir)
	if err != nil {
		return nil, err
	}

	return &track, nil
}
