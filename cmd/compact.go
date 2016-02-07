package cmd

import (
	"github.com/JustinAzoff/flow-indexer/flowindexer"
	"github.com/spf13/cobra"
)

var cmdCompact = &cobra.Command{
	Use:   "compact",
	Short: "Compact the database",
	Long:  "Compact the database",
	Run: func(cmd *cobra.Command, args []string) {
		flowindexer.RunCompact(dbpath)
	},
}

func init() {
	RootCmd.AddCommand(cmdCompact)
}
