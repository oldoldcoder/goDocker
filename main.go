package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

const usage = `goDocker is a simple container runtime implementation of Docker .
The purpose of this project is to learn how docker works and how to write a docker by ourselves
Enjoy it, just for fun.`

func main() {
	app := cli.NewApp()
	app.Name = "goDocker"
	app.Usage = usage
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}

	app.Before = func(ctx *cli.Context) error {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
