package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "查看版本",
	Run: func(cmd *cobra.Command, args []string) {
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version: 2.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
