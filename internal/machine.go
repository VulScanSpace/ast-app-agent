package internal

import (
	"github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/v3/host"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
)

var (
	MachineId string // 如何解决进程重启后，machineId 变化的问题？
)

func init() {
	MachineId, _ = host.HostID()
	// Note: if host.HostID has error, then use initMachineIdWithUUID
	// initMachineIdWithUUID
}

func initMachineIdWithUUID() {
	machinePath := path.Join(os.TempDir(), "ast-app-agent-machine")
	if data, err := ioutil.ReadFile(machinePath); err == nil {
		MachineId = string(data)
	}
	if len(MachineId) == 0 {
		writeMachineId(machinePath)
	}
}

func writeMachineId(machinePath string) {
	MachineId = uuid.NewV4().String()
	err := ioutil.WriteFile(machinePath, []byte(MachineId), fs.ModePerm)
	if err != nil {
		return
	}
}
