package cmd

import (
	"github.com/JustinAzoff/flow-indexer/flowindexer"
	"github.com/spf13/cobra"
)

var config string

var cmdDaemon = &cobra.Command{
	Use:   "daemon [args]",
	Short: "Start daemon",
	Long:  "Start daemon that indexes flows and exposes HTTP api",
	Run: func(cmd *cobra.Command, args []string) {
		flowindexer.RunDaemon(config)
	},
}

func init() {
	cmdDaemon.Flags().StringVarP(&config, "config", "c", "config.json", "configuration filename")
	RootCmd.AddCommand(cmdDaemon)
}
