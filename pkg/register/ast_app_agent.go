package register

import "github.com/sirupsen/logrus"
import "github.com/shirou/gopsutil/v3/host"

// Register ast-app-agent to ast-app-server
func Register() {
	info, _ := host.Info()
	logrus.Infof("register agent for machine: %s", info)
}
