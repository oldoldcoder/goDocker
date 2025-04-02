package subsystems

import "github.com/sirupsen/logrus"

type ResourceConfig struct {
	// 内存限制
	MemoryLimit string
	// CPU
	CpuShare string
	// CPU核数
	CpuSet string
}
type Subsystem interface {
	// Name 返回Subsystem的名字
	Name() string
	// Set 设置某个cgroup在这个subsystem中的资源限制
	Set(path string, res *ResourceConfig) error
	// Apply 将某个进程添加到某个cgroup中
	Apply(path string, pid int) error
	// Remove 移除某个cgroup
	Remove(path string) error
}

var (
	Ins = []Subsystem{
		&CpuSetSubSystem{},
		// 内存这个有问题在我们机器上面，先删除掉
		// &MemorySubsystem{},
		&CpuSubSystem{},
	}
)

type CgroupManager struct {
	// 创建的路径
	Path string
	// 具体的资源限制 TODO 这里没有设置貌似
	Resource *ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

// Apply 将pid加入到每个cgroup中
func (c *CgroupManager) Apply(pid int) error {
	for _, subsystem := range Ins {
		err := subsystem.Apply(c.Path, pid)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CgroupManager) Set(res *ResourceConfig) error {
	for _, subsystem := range Ins {
		err := subsystem.Set(c.Path, res)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CgroupManager) Destroy() error {
	for _, subsystem := range Ins {
		if err := subsystem.Remove(c.Path); err != nil {
			logrus.Warnf("remove subsystem %s failed,error is %v", subsystem.Name(), err)
		}
	}
	return nil
}
