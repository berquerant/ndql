package main

import (
	"github.com/berquerant/ndql/pkg/config"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version and exit",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMain(cmd, args, config.ModeVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	initFlags(versionCmd)
}
