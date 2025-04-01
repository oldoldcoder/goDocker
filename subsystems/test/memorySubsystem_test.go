package test

import (
	"GoDocker/subsystems"
	"os"
	"path"
	"testing"
)

func TestMemorySubsystem(t *testing.T) {
	CpuSubSystem := subsystems.CpuSubSystem{}
	resconfig := subsystems.ResourceConfig{
		CpuShare: "512",
	}
	testCgroup := "testCpuShares"
	// 在这里创建了
	if err := CpuSubSystem.Set(testCgroup, &resconfig); err != nil {
		t.Fatalf("cgroup fail %v", err)
	}
	t.Logf("path is %s", path.Join(subsystems.FindCgroupMountpoint("cpu"), testCgroup))
	stat, _ := os.Stat(path.Join(subsystems.FindCgroupMountpoint("cpu"), testCgroup))
	t.Logf("cgroup stats: %+v", stat)
	if err := CpuSubSystem.Apply(testCgroup, os.Getpid()); err != nil {
		t.Fatalf("cgroup Apply %v", err)
	}
	//将进程移回到根Cgroup节点
	if err := CpuSubSystem.Apply("", os.Getpid()); err != nil {
		t.Fatalf("cgroup Apply %v", err)
	}
	// 这里会remove掉
	if err := CpuSubSystem.Remove(testCgroup); err != nil {
		t.Fatalf("cgroup remove %v", err)
	}

}
