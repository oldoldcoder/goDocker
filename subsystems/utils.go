package subsystems

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

// FindCgroupMountpoint <mount ID> <parent ID> <major:minor> <root> <mount point> <mount options> - <filesystem type> <source> <super options>
// FindCgroupMountpoint 格式如上，Fileds[4]是挂载的位置,所以能够返回
func FindCgroupMountpoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("error:", err)
		}
	}(f)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}
	return ""
}

// GetCgroupPath 得到cgroup在文件系统中的绝对位置
func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountpoint(subsystem)
	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path.Join(cgroupRoot, cgroupPath), 0755); err == nil {
			} else {
				return "", fmt.Errorf("error create cgroup %v", err)
			}
		}
		return path.Join(cgroupRoot, cgroupPath), nil
	} else {
		return "", fmt.Errorf("cgroup path error %v", err)
	}
}
