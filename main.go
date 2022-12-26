package main

import (
	"os"

	"github.com/gookit/color"
	"github.com/mlange-42/track/cli"
	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
)

const version = "0.1.0"

func main() {
	if !color.SupportColor() || !isTerminal() {
		color.Disable()
	}

	track, err := core.NewTrack()
	if err != nil {
		out.Err("%s", err.Error())
		os.Exit(1)
	}

	if err := cli.RootCommand(&track, version).Execute(); err != nil {
		out.Err("%s", err.Error())
		os.Exit(1)
	}
}

func isTerminal() bool {
	o, _ := os.Stdout.Stat()
	return (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice
}
