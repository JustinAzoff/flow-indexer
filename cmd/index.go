package cmd

import (
	"github.com/JustinAzoff/flow-indexer/flowindexer"
	"github.com/spf13/cobra"
)

var cmdIndex = &cobra.Command{
	Use:   "index [args]",
	Short: "Index flows",
	Long:  "Index flows",
	Run: func(cmd *cobra.Command, args []string) {
		flowindexer.RunIndex(dbpath, args)
	},
}

func init() {
	RootCmd.AddCommand(cmdIndex)
}
