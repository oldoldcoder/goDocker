package main

import (
	"GoDocker/container"
	log "github.com/sirupsen/logrus"
	"os"
)

func Run(tty bool, cmd string) {
	parent := container.NewParentProcess(tty, cmd)
	// 启动一个线程
	if err := parent.Start(); err != nil {
		log.Fatal(err)
	}
	_ = parent.Wait()
	os.Exit(-1)
}
