package subsystems

import "testing"

func TestUtils(t *testing.T) {
	// 使用t测试自带的测试来写
	t.Logf("cpu subsystem mount point %v\n", FindCgroupMountpoint("cpu"))
	t.Logf("cpuset subsystem mount point %v\n", FindCgroupMountpoint("cpuset"))
	t.Logf("memory subsystem mount point %v\n", FindCgroupMountpoint("memory"))

	path, err := GetCgroupPath("memory", "testCgroup", true)
	if err != nil {
		t.Errorf("GetCgroupPath err:%v\n", err)
	}
	t.Logf("path:%v\n", path)
}
