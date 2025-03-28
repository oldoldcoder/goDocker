package subsystems

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type ResourceConfig struct {
	// 内存限制
	MemoryLimit string
	// CPU
	CpuShare string
	// CPU核数
	CpuSet string
}
type Cgroupv2 interface {
	// Create 创建一个cgroupv2
	Create() error
	// Name 返回Cgroupv2的名字
	Name() string
	// Set 设置某个cgroup在这个subsystem中的资源限制
	Set(res *ResourceConfig) error
	// Apply 将某个进程添加到某个cgroup中
	Apply(pid int) error
	// Remove 移除某个cgroup
	Remove() error
}

type Cgroup struct {
	path            string
	config          *ResourceConfig
	resourceFileMap map[string]string
}

func NewCgroup() *Cgroup {
	cgroup := new(Cgroup)
	cgroup.config = nil
	cgroup.path = ""
	cgroup.resourceFileMap = make(map[string]string)
	cgroup.resourceFileMap["MemoryLimit"] = "memory.max"
	cgroup.resourceFileMap["CpuShare"] = "cpu.shares"
	cgroup.resourceFileMap["CpuSet"] = "cpu.cfs_quota_us"
	return cgroup
}
func (cgroup *Cgroup) Create() error {
	src := rand.NewSource(time.Now().UnixNano()) // 创建新的随机数源
	rng := rand.New(src)                         // 创建随机数生成器
	// 在相应目录下创建一个cgroup
	cgroup.path = "/sys/fs/cgroup/limit" + strconv.Itoa(rng.Int())
	if err := os.MkdirAll(cgroup.path, 0755); err != nil {
		return fmt.Errorf("failed to create cgroup: %v", err)
	}
	return nil
}
func (cgroup *Cgroup) Name() string {
	return cgroup.path
}

// Set 按照config来配置具体的
func (cgroup *Cgroup) Set(res *ResourceConfig) error {
	cgroup.config = res

}
