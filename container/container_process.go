package container

import (
	"github.com/sirupsen/logrus"
	"os"
	"syscall"
)

func RunContainerInitProcess(command string, args []string) {

	logrus.Infof("command %s", command)
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NODEV | syscall.MS_NOSUID
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		logrus.Errorf("exec command %s err %v", command, err)
	}

	return
}
