package main

import (
	"fmt"
	"github.com/jiexun/minijkstool/core"
	v1 "github.com/jiexun/minijkstool/core/v1"
	"os"

	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{
		Name:    "minijks",
		Version: "0.5.0",
		Usage:   "inspect, unpack and pack and packfile Java keystore files",
		Commands: []*cli.Command{
			core.InspectCommand,
			core.UnpackCommand,
			core.PackCommand,
			v1.CertFilePackCommand,
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
