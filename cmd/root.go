/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/VulScanSpace/ast-app-agent/internal"
	"github.com/VulScanSpace/ast-app-agent/pkg/heartbeat"
	"github.com/VulScanSpace/ast-app-agent/pkg/hook"
	"github.com/VulScanSpace/ast-app-agent/pkg/inject"
	"github.com/VulScanSpace/ast-app-agent/pkg/register"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ast-app-agent",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if daemon {
			if internal.ProcessIsAlive() {
				logrus.Info("ast-app-agent is already running...")
				return
			}
			go internal.StartSock()

			hook.WatchSignal()
			register.Register()

			c := cron.New()
			schedule, _ := cron.ParseStandard("@every 1s")
			c.Schedule(schedule, heartbeat.HeartBeat{})
			schedule, _ = cron.ParseStandard("@every 1m")
			c.Schedule(schedule, inject.Inject{})
			c.Start()
			select {}
		} else {
			inject.Inject{}.Run()
		}
	},
}

var (
	daemon bool
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ast-app-agent.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolVarP(&daemon, "daemon", "s", false, "run with daemon mode")
}
