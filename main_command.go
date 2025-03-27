package main

import (
	"GoDocker/container"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"syscall"
)

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
		tty := ctx.BoolT("it")
		// 调用run function启动容器
		Run(tty, cmd)

		// 确保程序退出时挂载 /proc
		defer func() {
			// 在这里挂载 /proc
			if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
				log.Errorf("Failed to mount /proc: %v", err)
			} else {
				log.Infof("/proc mounted successfully")
			}
		}()

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
		container.RunContainerInitProcess(cmd, nil)
		return nil
	},
}
