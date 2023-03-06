package inject

import (
	"fmt"
	"github.com/VulScanSpace/ast-app-agent/internal"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

var (
	waitingInjectedPid map[string]int64
	injectedPid        map[string]bool
	Timeout            int64
	JAttach            string
	IastCmd            string
)

type Inject struct{}

func init() {
	currentDir, _ := os.Getwd()
	IastCmd = fmt.Sprintf("-javaagent:%s=%s&%s", path.Join(currentDir, "libs", "ast-agent.jar"), path.Join(currentDir, "libs", "IAST-Engine.jar"), "INSTALL")
	waitingInjectedPid = map[string]int64{}
	injectedPid = map[string]bool{}
	// Timeout TODO read timeout from config or remote api
	Timeout = 60 * 10
	switch CurrentOs {
	case Windows:
		// TODO add arch check
		JAttach = "/opt/ast-app/bin/jattach-windows.exe"
		break
	case Linux:
		JAttach = "/opt/ast-app/bin/jattach-linux"
		break
	case Darwin:
		JAttach = "/opt/ast-app/bin/jattach-darwin"
	}
}

func (Inject) Run() {
	supportInject, javaProcessList := CheckDeploymentRequirements()

	if len(javaProcessList) == 0 {
		logrus.Info("no java process.")
	}

	if !supportInject {
		logrus.Info("mem limit")
	}

	removeTimeoutPid()
	for _, javaPid := range javaProcessList {
		if _, hasKey := waitingInjectedPid[javaPid]; hasKey {
			logrus.Infof("injecting iast to %s", javaPid)
			continue
		}
		if _, hasKey := injectedPid[javaPid]; hasKey {
			logrus.Debugf("already injected iast to %s", javaPid)
			continue
		}

		go inject(javaPid)
	}
}

func removeTimeoutPid() {
	now := time.Now().Unix()
	var timeoutPid []string
	for pid, startTime := range waitingInjectedPid {
		if now-startTime > Timeout {
			timeoutPid = append(timeoutPid, pid)
		}
	}

	for _, pid := range timeoutPid {
		stopInject(pid)
	}
}

func inject(pid string) {
	logrus.Infof("start inject iast to %s", pid)
	waitingInjectedPid[pid] = time.Now().Unix()

	attachArgs := []string{pid, "load", "instrument", "false", IastCmd}
	logrus.Infof("run inject iast to %s. cmd: %s %v", pid, JAttach, attachArgs)
	err, output, errMsg := internal.RunShellCmd(JAttach, attachArgs...)

	delete(waitingInjectedPid, pid)

	if err == nil {
		injectedPid[pid] = true
		logrus.Infof("finished inject iast to %s", pid)
	} else {
		logrus.Errorf("failure inject iast to %s, reason: %v\noutput:\n %s\nerror:\n %s", pid, err, output, errMsg)
	}
}

// stopInject ast-app agent use jAttach, needn't kill inject process
func stopInject(pid string) {
	logrus.Infof("remove timeout inject process, pid: %s", pid)
	delete(waitingInjectedPid, pid)
}
