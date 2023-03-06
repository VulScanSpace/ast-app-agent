package heartbeat

import (
	"github.com/VulScanSpace/ast-app-agent/internal"
	"github.com/VulScanSpace/ast-app-agent/pkg/inject"
	"github.com/sirupsen/logrus"
)

type HeartBeat struct{}

func (HeartBeat) Run() {
	javaProcessList := inject.ReadJavaProcess()
	aliveAgentSize := inject.CheckAgentStatus()

	// TODO send heartbeat msg to remote server
	logrus.Infof("machine: %s, javaProcess: %d, agentAlive: %d, status: up", internal.MachineId, len(javaProcessList), aliveAgentSize)

	// TODO read heartbeat api and to do action.
}

func Stop(reason string) {
	javaProcessList := inject.ReadJavaProcess()
	aliveAgentSize := inject.CheckAgentStatus()
	logrus.Infof("machine: %s, javaProcess: %d, agentAlive: %d, status: down, reason: %s", internal.MachineId, len(javaProcessList), aliveAgentSize, reason)
}
