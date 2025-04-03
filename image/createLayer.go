package image

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

// NewWorkSpace 创建容器文件系统
func NewWorkSpace(rootURL string, mntURL string) {

}

// CreateReadOnlyLayer 将busybox.tar解压到busybox目录下，作为容器的只读层
func CreateReadOnlyLayer(rootURL string) {
	busyboxURL := rootURL + "busybox/"
	busyboxTarURL := rootURL + "busybox.tar"
	exist, err := PathExists(busyboxTarURL)
	if err != nil {
		log.Infof("Fail to judge whether dir %s exists. %v", busyboxURL, err)
	}
	if exist == false {
		if err := os.Mkdir(busyboxTarURL, 0755); err != nil {
			log.Errorf("fail to create dir %s: %v", busyboxTarURL, err)
		}
		if _, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
			log.Errorf("unTar dir %s error %v", busyboxTarURL, err)
		}
	}
}

// CreateWriteLayer 创建了一个名为writeLayer的文件夹作为容器唯一可写层
func CreateWriteLayer(rootURL string) {
	writeURL := rootURL + "writeLayer/"
	if err := os.MkdirAll(writeURL, os.ModePerm); err != nil {
		log.Errorf("mkdir dir %s error %v", writeURL, err)
	}
}

// CreateMountPoint TODO overlay可能和aufs挂载方式不一样
// CreateMountPoint 创建挂载点，也就是aufs或者overlay两个点
func CreateMountPoint(rootURL string, mntURL string) {
	if err := os.MkdirAll(mntURL, os.ModePerm); err != nil {
		log.Errorf("mkdir dir %s error %v", mntURL, err)
	}
	// 把writeLayer目录和busybox目录mount到mnt目录下
	dirs := "dirs=" + rootURL + "writeLayer:" + rootURL + "busybox"
	cmd := exec.Command("mount", "-t", "overlay", "-o", dirs, "none", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		return true, err
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}
