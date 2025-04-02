package subsystems

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

type CpuSetSubSystem struct {
}

// Set 设置cgroupPath中对应Cgroup内存的限制
func (subsystem *CpuSetSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	// 获取当前 subsystem 在虚拟文件系统中的路径
	subsysCgroupPath, err := GetCgroupPath(subsystem.Name(), cgroupPath, true)
	if err != nil {
		return err // 如果获取路径失败，直接返回错误
	}

	// 如果 CpuShare 为空，不需要设置
	if res.CpuSet == "" {
		return nil
	}
	// 查找当前的核心数，如果超过了两个，则报错
	err = subsystem.checkCPUConflict(subsysCgroupPath, res.CpuSet)
	if err != nil {
		return err
	}

	// 设置 cgroup 的内存限制
	if err := os.WriteFile(path.Join(subsysCgroupPath, "cpuset.cpus"), []byte(res.CpuSet), 0644); err != nil {
		return fmt.Errorf("set cgroup cpu.shares fail: %v", err)
	}

	return nil
}

func (subsystem *CpuSetSubSystem) checkCPUConflict(subsysCgroupPath, desiredCPUs string) error {
	// 打开 cpuset.effective_cpus 文件
	f, _ := os.Open(path.Join(subsysCgroupPath, "cpuset.effective_cpus"))

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Printf("Error closing cpuset.effective_cpus: %s\n", err)
		}
	}(f)

	// 读取文件内容
	var effectiveCPUs string
	_, _ = fmt.Fscanf(f, "%s", &effectiveCPUs)
	// 获取有效的 CPU 核心
	effectiveCPUsList := strings.Split(effectiveCPUs, ",")

	// 获取当前设置的 CPU 个数
	desiredCPUsList := strings.Split(desiredCPUs, ",")

	// 比较是否存在冲突
	for _, cpu := range desiredCPUsList {
		// 检查每个想要使用的核心是否在有效的核心列表中
		found := false
		for _, effectiveCPU := range effectiveCPUsList {
			if cpu == effectiveCPU {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("desired CPU %s is not in the effective CPUs list", cpu)
		}
	}

	return nil
}

func (subsystem *CpuSetSubSystem) Remove(cgroupPath string) error {
	if subsysCgroupPath, err := GetCgroupPath(subsystem.Name(), cgroupPath, false); err == nil {
		// 删除对应的目录
		return os.Remove(subsysCgroupPath)
	} else {
		return err
	}
}

func (subsystem *CpuSetSubSystem) Name() string {
	return "cpuset"
}
func (subsystem *CpuSetSubSystem) Apply(cgroupPath string, pid int) error {
	if subsysCgroupPath, err := GetCgroupPath(subsystem.Name(), cgroupPath, false); err == nil {
		if err := os.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("apply cgroup cpuset fail: %v", err)
		}
		return nil
	} else {
		return err
	}
}
