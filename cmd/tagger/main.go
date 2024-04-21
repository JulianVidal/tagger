package tagger

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:   "tagger",
	Short: "tagger - a simple CLI tool to tag anything",
	Long:  `tagger's goal is to be able to add and search files quickly in an intuitive way. Instead of memorising a given structure, use tags to organise files`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	if err := mainCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
