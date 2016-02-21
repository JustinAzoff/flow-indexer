package cmd

import (
	"github.com/JustinAzoff/flow-indexer/flowindexer"
	"github.com/spf13/cobra"
)

var cmdIndexAll = &cobra.Command{
	Use:   "indexall",
	Short: "Index all flows",
	Long:  "Index all flows once.",
	Run: func(cmd *cobra.Command, args []string) {
		flowindexer.RunIndexAll(config)
	},
}

func init() {
	cmdIndexAll.Flags().StringVarP(&config, "config", "c", "config.json", "configuration filename")
	RootCmd.AddCommand(cmdIndexAll)
}
