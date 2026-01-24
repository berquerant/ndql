package main

import (
	"fmt"

	"github.com/berquerant/ndql/pkg/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "ls PATH",
	Short: "List paths and exit",
	Long: fmt.Sprintf(`List paths and exit.

This command is equivalent to:

    ndql query 'select *' PATH

## PATH
%s`, config.DescribeSourceUsage()),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMain(cmd, args, config.ModeList)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	initFlags(listCmd)
}
