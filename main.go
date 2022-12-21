package main

import (
	"os"

	"github.com/mlange-42/track/cli"
	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
)

const version = "v0.1.0"

func main() {
	track := core.Track{}
	track.CreateDirs()

	if err := cli.RootCommand(&track, version).Execute(); err != nil {
		out.Err("%s\n", err.Error())
		os.Exit(1)
	}
}
