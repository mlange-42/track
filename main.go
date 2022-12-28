package main

import (
	"os"

	"github.com/gookit/color"
	"github.com/mlange-42/track/cli"
	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
)

const version = "0.2.1"

func main() {
	if !color.Support256Color() || !isTerminal() {
		color.Disable()
	}

	track, err := core.NewTrack()
	if err != nil {
		out.Err("%s\n", err.Error())
		os.Exit(1)
	}

	if err := cli.RootCommand(&track, version).Execute(); err != nil {
		out.Err("%s\n", err.Error())
		os.Exit(1)
	}
}

func isTerminal() bool {
	o, _ := os.Stdout.Stat()
	return (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice
}
