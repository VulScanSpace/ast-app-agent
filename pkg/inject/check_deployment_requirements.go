package inject

import (
	"fmt"
	"github.com/VulScanSpace/ast-app-agent/internal"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
	"strings"
)

type OsType string

const (
	Unknown OsType = "unknown"
	Windows OsType = "windows"
	Linux   OsType = "linux"
	Darwin  OsType = "darwin"
)

var (
	CurrentOs    OsType
	IsContainer  bool
	CGroupName   string
	CGroupPath   string
	IsPrivileged bool
	JavaHome     string
	JpsPath      string
)

func init() {
	CurrentOs = readMachineInfo()
	if CurrentOs != Windows {
		IsContainer, CGroupName, CGroupPath = readCGroup()
	}
	JavaHome = checkJdk()
	JpsPath = checkJps(JavaHome)
}

// CheckDeploymentRequirements memLimit、javaProcess, etc.
func CheckDeploymentRequirements() (bool, []string) {
	if supportInject := checkMemLimit(); !supportInject {
		return false, nil
	}
	return true, ReadJavaProcess()
}

func readMachineInfo() OsType {
	logrus.Debug("start read MachineInfo")
	switch runtime.GOOS {
	case "windows":
		return Windows
	case "linux":
		return Linux
	case "darwin":
		return Darwin
	default:
		logrus.Errorf("machine info read failure. os: %s, arch: %s, go: %s", runtime.GOOS, runtime.GOARCH, runtime.Version())
		return Unknown
	}
}

func checkJdk() string {
	err, output, errMsg := internal.RunShellCmd("bash", "-c", "which java")
	if err != nil {
		logrus.Errorf("check JDK environment failure, reason: %v, msg: %s", err, errMsg)
		return ""
	}
	output = strings.TrimSpace(output)
	if len(output) == 0 {
		logrus.Info("JDK environment does not exist.")
		return ""
	}
	logrus.Debugf("JDK environment exist. path: %s", output)
	return strings.TrimSuffix(output, "/java")
}

/**
要在容器中区分容器进程和其他容器进程，可以使用不同的 cgroup 名称和路径。例如，可以为每个容器创建一个名为 mycontainer 的 cgroup，但是将它们放在不同的路径中，例如 /sys/fs/cgroup/pids/container1/mycontainer 和 /sys/fs/cgroup/pids/container2/mycontainer。这样，可以使用路径来区分容器进程和其他容器进程。
*/
func readCGroup() (hasCGroup bool, name, path string) {
	err, output, errMsg := internal.RunShellCmd("grep", "-E", "'^1:name=systemd:/.+'", "/proc/self/cgroup")

	if err != nil {
		logrus.Debugf("check CGroup failure, reason: %v, msg: %s", err, errMsg)
		return
	}

	output = strings.TrimSpace(output)
	if output == "" {
		return
	}

	items := strings.Split(output, ":")
	_, name, path = items[0], items[1], items[2]
	name = strings.TrimPrefix(name, "name=")
	hasCGroup = true
	return
}

func checkJps(javaHome string) (jpsPath string) {
	err, output, errMsg := internal.RunShellCmd("bash", "-c", "which jps")
	output = strings.TrimSpace(output)
	if err != nil || len(output) == 0 {
		logrus.Errorf("check jps process failure, reason: %v, msg: %s", err, errMsg)
		if len(javaHome) > 0 {
			os.Stat(path.Join(javaHome, "jps"))
		}
		return ""
	}

	logrus.Debugf("jps process exist. path: %s", output)
	return strings.TrimSpace(output)
}

func ReadJavaProcess() []string {
	// 检查内存信息
	if len(JpsPath) > 0 {
		if processList := readJavaProcessWithJps(); len(processList) == 0 {
			return nil
		} else {
			return processList
		}
	} else {
		return readProcessWithOutJps("java")
	}
}

func readJavaProcessWithJps() (processIds []string) {
	err, output, errMsg := internal.RunShellCmd(JpsPath, "-l")
	if err != nil {
		logrus.Errorf("ReadJavaProcess failure, reason: %v, msg: %s", err, errMsg)
		return
	}
	output = strings.TrimSpace(output)
	if len(output) == 0 {
		logrus.Info("Java process does not exist.")
		return
	}

	var processId, processName string
	processes := strings.Split(output, "\n")
	for _, processInfo := range processes {
		processItems := strings.Split(processInfo, " ")
		processId = processItems[0]

		if len(processItems) >= 2 {
			processName = strings.Join(processItems[1:], " ")
			// TODO exclude: inject self、jps、jShell、jstat, etc; it should be remote config.
			if len(processName) > 0 && strings.Contains(processName, "jdk.jcmd/sun.tools.jps.Jps") {
				continue
			}
			if strings.Contains(processName, "jdk.jshell.execution.RemoteExecutionControl") {
				continue
			}
			if strings.Contains(processName, "jdk.jshell/jdk.internal.jshell.tool.JShellToolProvider") {
				continue
			}
		}

		processIds = append(processIds, processId)
	}
	logrus.Debugf("finished read alive java process, count: %d", len(processIds))
	return
}

func readProcessWithOutJps(name string) []string {
	var (
		err    error
		output string
		errMsg string
	)
	if CurrentOs == Windows {
		err, output, errMsg = internal.RunShellCmd("cmd", "/c", fmt.Sprintf("wmic process where \"caption like '%%%s%%'\" get processid, commandline", name))
	} else {
		if IsContainer {
			// TODO add process get logic
		} else {
			err, output, _ = internal.RunShellCmd("sh", "-c", fmt.Sprintf("ps aux | grep -v grep | grep -v defunct | grep %s", name))
		}
	}

	if err != nil {
		logrus.Errorf("read %s process failure, reason: %v, msg: %s", name, err, errMsg)
		return nil
	}

	processes := strings.Split(output, "\n")
	for _, process := range processes {
		if strings.Contains(process, "java") {
			fmt.Println(process)
		}
	}
	return nil
}

/**
可以通过检查 `/proc/1/cgroup` 文件来判断当前 shell 是否在容器中运行。如果文件包含 `docker` 或 `lxc` 字符串，则表示当前 shell 在容器中运行。如果文件包含 `machine` 或 `systemd` 字符串，则表示当前 shell 在虚拟机中运行。如果文件包含 `/` 字符，则表示当前 shell 在 chroot 环境中运行。如果文件为空，则表示当前 shell 在宿主机中运行。
*/
//func checkIsDocker() bool {
//	err, output, errMsg := utils.RunShellCmd("cat", "/proc/1/cgroup")
//	output = strings.TrimSpace(output)
//	if err != nil || output == "" {
//		logrus.Errorf("check privileged failure, reason: %s", errMsg)
//		return false
//	}
//
//	if strings.Contains("")
//		cgroupItems := strings.Split(output, ":")
//	_, cgroupName, cgroupPath := cgroupItems[0], cgroupItems[1], cgroupItems[2]
//	cgroupName = strings.TrimPrefix(cgroupName, "name=")
//}

func checkMemLimit() bool {
	v, _ := mem.VirtualMemory()
	logrus.Debugf("memory check: total[%v], free[%v], used-percent[%f%%]\n", v.Total, v.Free, v.UsedPercent)
	return v.UsedPercent > 0.8
}

func CheckAgentStatus() (aliveAgentSize int) {
	return 0
}
