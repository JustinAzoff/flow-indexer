package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	VERSION = "0.1.4"
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Output version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(VERSION)
	},
}

func init() {
	RootCmd.AddCommand(cmdVersion)
}
