package main

import (
	"fmt"

	"github.com/berquerant/ndql/pkg/config"
	"github.com/spf13/cobra"
)

var dryrunCmd = &cobra.Command{
	Use:   "dry QUERY",
	Short: "Parse query and exit",
	Long: fmt.Sprintf(`Parse query and exit.

## QUERY
%s

## Examples

    ndql dry 'select path, mod_time where not is_dir'

as json:

    ndql dry 'select path, mod_time where not is_dir' -v
`, config.DescribeSourceUsage()),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMain(cmd, args, config.ModeDryrun)
	},
}

func init() {
	rootCmd.AddCommand(dryrunCmd)
	initFlags(dryrunCmd)
}
