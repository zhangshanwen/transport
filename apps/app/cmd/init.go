package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{Use: "app"}

func Execute() {
	rootCmd.AddCommand(runCmd, stopCmd)
	_ = rootCmd.Execute()
}
