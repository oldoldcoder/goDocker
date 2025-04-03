package image

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"syscall"
)

// pivotRoot 修改rootfs系统调用
// @params root 代表新的根文件系统的路径
// TODO 在执行之前务必保证自己使用了MOUNT NAMESPACE 隔离空间否则一定会出错
func pivotRoot(root string) error {
	if err := syscall.Mount(root, root, "bind",
		syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount rootfs to it self error: %v", err)
	}
	// 创建rootfs/.pivot_root 存储old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}
	// pivot_root 到新的rootfs，老的old_root现在挂载在rootfs/.pivot_root上
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root failed: %v", err)
	}
	// 修改当前的目录到根目录
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir /: %v", err)
	}
	pivotDir = filepath.Join("/", ".pivot_root")
	// umount rootfs/.pivot_root
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root failed: %v", err)
	}
	return os.Remove(pivotDir)
}

func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		logrus.Errorf("get pwd: %v", err)
		return
	}
	logrus.Infof("current pwd: %s", pwd)

	err = pivotRoot(pwd)
	if err != nil {
		logrus.Errorf("pivot_root failed: %v", err)
	}
	// mount proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	syscall.Mount("tmpfs", "/dev", "tmpfs", uintptr(defaultMountFlags), "")

}
