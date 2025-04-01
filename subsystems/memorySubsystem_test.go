package subsystems

import (
	"os"
	"path"
	"testing"
)

func TestMemorySubsystem(t *testing.T) {
	memSubSystem := MemorySubsystem{}
	resconfig := ResourceConfig{
		MemoryLimit: "100m",
	}
	testCgroup := "testMemlimit"
	// 在这里创建了
	if err := memSubSystem.Set(testCgroup, &resconfig); err != nil {
		t.Fatalf("cgroup fail %v", err)
	}
	stat, _ := os.Stat(path.Join(FindCgroupMountpoint("memory"), testCgroup))
	t.Logf("cgroup stats: %+v", stat)
	if err := memSubSystem.Apply(testCgroup, os.Getpid()); err != nil {
		t.Fatalf("cgroup Apply %v", err)
	}
	//将进程移回到根Cgroup节点
	if err := memSubSystem.Apply("", os.Getpid()); err != nil {
		t.Fatalf("cgroup Apply %v", err)
	}

	if err := memSubSystem.Remove(testCgroup); err != nil {
		t.Fatalf("cgroup remove %v", err)
	}

}
