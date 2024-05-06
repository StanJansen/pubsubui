package main

import (
	"flag"

	"github.com/stanjansen/pubsubui/internal/app"
)

func main() {
	c := app.Config{}
	flag.StringVar(&c.PubsubEmulatorHost, "emulator", "", "The host and port of the Pub/Sub emulator")
	flag.StringVar(&c.Project, "project", "", "The GCP project name")
	flag.Parse()

	app := app.New(c)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
