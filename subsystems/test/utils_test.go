package test

import (
	"GoDocker/subsystems"
	"testing"
)

func TestUtils(t *testing.T) {
	// 使用t测试自带的测试来写
	t.Logf("cpu subsystem mount point %v\n", subsystems.FindCgroupMountpoint("cpu"))
	t.Logf("cpuset subsystem mount point %v\n", subsystems.FindCgroupMountpoint("cpuset"))
	// t.Logf("memory subsystem mount point %v\n", subsystems.FindCgroupMountpoint("memory"))

	path, err := subsystems.GetCgroupPath("cpu", "testCgroup", true)
	if err != nil {
		t.Errorf("GetCgroupPath err:%v\n", err)
	}
	t.Logf("path:%v\n", path)
}
