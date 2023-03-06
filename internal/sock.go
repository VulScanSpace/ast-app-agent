package internal

import (
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"path"
)

var (
	AgentSock string
)

func init() {
	AgentSock = path.Join(os.TempDir(), "ast-app-agent.sock")
}

func StartSock() {
	_ = os.Remove(AgentSock)
	unixSock, err := net.ResolveUnixAddr("unix", AgentSock)
	if err != nil {
		logrus.Errorf("StartSock ResolveUnixAddr error, reason: %v", err)
		return
	}

	listen, err := net.ListenUnix("unix", unixSock)
	for {
		_, _ = listen.Accept()
	}
}

func ProcessIsAlive() bool {
	unixSock, err := net.ResolveUnixAddr("unix", AgentSock)
	if err != nil {
		return false
	}
	_, err = net.DialUnix("unix", nil, unixSock)
	return err == nil
}
