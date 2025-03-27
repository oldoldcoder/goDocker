package main

import (
	"fmt"
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
	var runCommand = cli.Command{
		Name:  "run",
		Usage: "create a container with namespace and cgroups limit goDocker run -it",
		Flags: []cli.Flag{
			cli.BoolTFlag{
				Name:  "it",
				Usage: "enable tty",
			},
		},
		Action: func(ctx *cli.Context) error {
			if len(ctx.Args()) < 1 {
				return fmt.Errorf("missing container command")
			}
			cmd := ctx.Args().Get(0)
			tty := ctx.BoolT("ti")
			// 调用run function启动容器
			Run(tty, cmd)
			return nil
		},
	}
	var initCommand = cli.Command{
		Name: "init",
		Usage: "Init container process run user’s process in container. " +
			"Do not call it outside",
		Action: func(ctx *cli.Context) error {
			log.Infof("init come on")
			cmd := ctx.Args().Get(0)
			log.Infof("cmd: %s", cmd)
			err := container.RuncontainerInitProcess(cmd, nil)
			return err
		},
	}

}
