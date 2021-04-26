package main

import (
	"fmt"
	"github.com/jiexun/minijkstool/core"
	v1 "github.com/jiexun/minijkstool/core/v1"
	"gopkg.in/urfave/cli.v2"
	"os"
)

func main() {

	/*pemraw, err := ioutil.ReadFile("./jsk-test/certs/truststore.jks")
	if err != nil {
		fmt.Println(err)
	}
	jiami := base64.StdEncoding.EncodeToString(pemraw)
	fmt.Println(jiami)*/


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
