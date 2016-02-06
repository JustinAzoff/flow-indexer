package cmd

import (
	"github.com/JustinAzoff/flow-indexer/flowindexer"
	"github.com/spf13/cobra"
)

var cmdExpandCIDR = &cobra.Command{
	Use:   "expandcidr [args]",
	Short: "Expand a CIDR range from those seen in the database",
	Long:  "Expand a CIDR range from those seen in the database",
	Run: func(cmd *cobra.Command, args []string) {
		flowindexer.RunExpandCIDR(dbpath, args)
	},
}

func init() {
	RootCmd.AddCommand(cmdExpandCIDR)
}
