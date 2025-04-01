package container

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func RunContainerInitProcess() error {
	// 从管道读取用户的命令，然后执行用户的命令
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("run container get user command error, cmdArray is nil")
	}
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		return fmt.Errorf("run container get path error, err:%v", err)
	}
	logrus.Infof("find path:%v", path)
	// 替换当前进程，占用pid 1的进程
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}

	return nil
}
func readUserCommand() []string {
	pipe := os.NewFile(3, "pipe")
	msg, err := io.ReadAll(pipe)
	if err != nil {
		logrus.Errorf("read pipe err %v", err)
		return nil
	}
	return strings.Split(string(msg), " ")
}

/*
*
设置挂载的点
*/
func setUpMount() {

}

func pivotRoot() error {
	return nil
}
