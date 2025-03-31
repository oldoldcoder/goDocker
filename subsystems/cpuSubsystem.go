package subsystems

import (
	"fmt"
	"os"
	"path"
	"strconv"
)

type CpuSubSystem struct {
}

// Set 设置cgroupPath中对应Cgroup内存的限制
func (subsystem *CpuSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	// 获取当前 subsystem 在虚拟文件系统中的路径
	subsysCgroupPath, err := GetCgroupPath(subsystem.Name(), cgroupPath, true)
	if err != nil {
		return err // 如果获取路径失败，直接返回错误
	}

	// 如果 CpuShare 为空，不需要设置
	if res.CpuShare == "" {
		return nil
	}

	// 设置 cgroup 的内存限制
	if err := os.WriteFile(path.Join(subsysCgroupPath, "cpu.shares"), []byte(res.CpuShare), 0644); err != nil {
		return fmt.Errorf("set cgroup cpu.shares fail: %v", err)
	}

	return nil
}

func (subsystem *CpuSubSystem) Remove(cgroupPath string) error {
	if subsysCgroupPath, err := GetCgroupPath(subsystem.Name(), cgroupPath, false); err == nil {
		// 删除对应的目录
		return os.Remove(subsysCgroupPath)
	} else {
		return err
	}
}

func (subsystem *CpuSubSystem) Name() string {
	return "cpu"
}
func (subsystem *CpuSubSystem) Apply(cgroupPath string, pid int) error {
	if subsysCgroupPath, err := GetCgroupPath(subsystem.Name(), cgroupPath, false); err == nil {
		if err := os.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("apply cgroup memory fail: %v", err)
		}
		return nil
	} else {
		return err
	}
}
