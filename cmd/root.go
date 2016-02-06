package cmd

import (
	"github.com/spf13/cobra"
)

var dbpath string

var RootCmd = &cobra.Command{
	Use:   "flow-indexer",
	Short: "flow-indexer indexes flows",
	Long:  "flow-indexer indexes flows",
	Run:   nil,
}

func init() {
	RootCmd.PersistentFlags().StringVar(&dbpath, "dbpath", "flows.db", "Database path")
}
