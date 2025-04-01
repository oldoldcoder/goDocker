package main

import (
	"GoDocker/container"
	"GoDocker/subsystems"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func Run(tty bool, cmd []string, res *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	// 启动一个线程
	if err := parent.Start(); err != nil {
		log.Fatal(err)
	}
	// 创建cgroup manager 并且通过调用set和apply设置资源限制并且使限制在容器上生效
	cgroupManager := subsystems.NewCgroupManager("docker_cgroup")
	defer cgroupManager.Destroy()
	// 设置开启相应的文件夹
	err := cgroupManager.Set(res)
	if err != nil {
		log.Errorf("cgroup manager set error:%v", err)
	}
	// 应用到pid上面
	err = cgroupManager.Apply(parent.Process.Pid)
	if err != nil {
		log.Errorf("cgroup manager apply error:%v", err)
	}
	// 发送信息？
	_ = parent.Wait()
	sendInitCommand(cmd, writePipe)
	os.Exit(-1)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command:%s", command)
	_, err := writePipe.WriteString(command)
	if err != nil {
		log.Errorf("write pipe error:%v", err)
	}
	err = writePipe.Close()
	if err != nil {
		log.Errorf("close pipe error:%v", err)
	}
}
