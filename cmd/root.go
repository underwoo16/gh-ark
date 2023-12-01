package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "gh-diffstack [command]",
	Short:         "Command line tool for managing stacked diffs.",
	Long:          `A command line tool for managing pull requests of stacked diffs through gh cli`,
	Example:       `gh-diffstack pr`,
}

func Execute() error {
	cmd, _, err := rootCmd.Find(os.Args[1:])
	if err != nil || cmd == nil {
		log.Fatal("No command found")
	}

	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(prCmd)
}
