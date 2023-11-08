package service

import (
	"wsfuzz/cmd"

	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "services",
	Short: "特定服务爆破",
}

func init() {
	cmd.RootCmd.AddCommand(serviceCmd)

	serviceCmd.AddCommand(serviceUptimeKumaCmd)
}
