package __2

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

const cgroupMemoryHierarchyMount = "/sys/fs/cgroup"

func UseCgroupV2Demo() {
	// 创建子进程二次执行
	if os.Args[0] == "/proc/self/exe" {
		// 在单独的命名空间namespace里面，id是1
		fmt.Println("current pid: ", os.Getpid())
		cmd := exec.Command("sh", "-c", "stress --vm-bytes 200m"+" --vm-keep -m 1")
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// 创建第二个子进程
	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// 启动子进程而不阻塞主进程
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		// 在这里执行
		fmt.Printf("子进程id:%v\n", cmd.Process.Pid)
		// 创建，即使存在也不会出错
		err := os.MkdirAll(path.Join(cgroupMemoryHierarchyMount, "memory_limit"), 0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		_ = os.WriteFile(path.Join(cgroupMemoryHierarchyMount, "memory_limit", "cgroup.procs"),
			[]byte(strconv.Itoa(cmd.Process.Pid)), 0644)
		_ = os.WriteFile(path.Join(cgroupMemoryHierarchyMount, "memory_limit",
			"memory.max"), []byte("104857600"), 0644)
	}

	// 等待子进程完成
	_, _ = cmd.Process.Wait()
}
