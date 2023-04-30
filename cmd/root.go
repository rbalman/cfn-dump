package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var pattern string

var rootCmd = &cobra.Command{
	Use:     "cfnd",
	Version: "0.0.1-beta",
	Short:   "Take snapshot of current cfn dependency and analyze it as json",
	Long:    `Dump cloudformation stacks export/import dependencies as a json file and analyze it for further processing`,
}

func init() {
	dumpCmd.PersistentFlags().StringVarP(&pattern, "pattern", "p", "", "specify the pattern of the cfn stacks that you want dump dependency of")
	rootCmd.AddCommand(dumpCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
