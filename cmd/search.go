package cmd

import (
	"github.com/JustinAzoff/flow-indexer/flowindexer"
	"github.com/spf13/cobra"
)

var cmdSearch = &cobra.Command{
	Use:   "search [args]",
	Short: "Search flows",
	Long:  "Search flows",
	Run: func(cmd *cobra.Command, args []string) {
		flowindexer.RunSearch(dbpath, args)
	},
}

func init() {
	RootCmd.AddCommand(cmdSearch)
}
