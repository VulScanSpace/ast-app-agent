package hook

import (
	"github.com/VulScanSpace/ast-app-agent/pkg/heartbeat"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func WatchSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL) // 其中 SIGKILL = kill -9 <pid> 可能无法截获
	go func() {
		logrus.Info("watching stop signals")
		s := <-c
		switch s {
		case os.Interrupt:
			logrus.Info("got Ctrl-C signal, try stop...")
			heartbeat.Stop("Ctrl-C")
			break
		case syscall.SIGTERM:
			logrus.Info("stop by 'kill <pid>'")
			heartbeat.Stop("stop by 'kill <pid>'")
		case syscall.SIGINT:
			logrus.Info("SIGINT")
			heartbeat.Stop("SIGINT")
		case syscall.SIGKILL:
			logrus.Info("SIGKILL")
			heartbeat.Stop("SIGKILL")
		}

		logrus.Info("===== ready to exit on SIGTERM =====")
		os.Exit(1)
	}()
}
