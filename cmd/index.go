package cmd

import (
	"github.com/JustinAzoff/flow-indexer/flowindexer"
	"github.com/spf13/cobra"
)

var backend_type string
var cmdIndex = &cobra.Command{
	Use:   "index [args]",
	Short: "Index flows",
	Long:  "Index flows",
	Run: func(cmd *cobra.Command, args []string) {
		flowindexer.RunIndex(dbpath, backend_type, args)
	},
}

func init() {
	cmdIndex.Flags().StringVarP(&backend_type, "backend", "b", "bro", "Log Backend")
	RootCmd.AddCommand(cmdIndex)
}
