package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version string = "v0.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
