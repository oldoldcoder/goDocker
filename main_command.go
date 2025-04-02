package main

import (
	"GoDocker/container"
	"GoDocker/subsystems"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name:  "run",
	Usage: "create a container with namespace and cgroups limit goDocker run -it",
	Flags: []cli.Flag{
		cli.BoolTFlag{
			Name:  "it",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
	},
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing container command")
		}
		var cmd []string
		// Args是解析后去除run命令以及其他匹配参数之外，其他的参数内容
		for _, arg := range ctx.Args() {
			cmd = append(cmd, arg)
		}
		tty := ctx.BoolT("it")
		res := &subsystems.ResourceConfig{
			MemoryLimit: ctx.String("m"),
			CpuSet:      ctx.String("cpuset"),
			CpuShare:    ctx.String("cpushare"),
		}
		// 调用run function启动容器
		Run(tty, cmd, res)

		return nil
	},
}

var initCommand = cli.Command{
	Name: "init",
	Usage: "Init container process run user’s process in container. " +
		"Do not call it outside",
	Action: func(ctx *cli.Context) error {
		log.Infof("init come on")
		container.RunContainerInitProcess()
		return nil
	},
}
