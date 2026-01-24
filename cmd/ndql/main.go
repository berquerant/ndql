package main

import (
	"context"
	"errors"
	"os"

	_ "github.com/pingcap/tidb/pkg/types/parser_driver"

	"github.com/berquerant/ndql/pkg/errorx"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ndql",
	Short: "Select metadata from files by SQL",
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Help()
	},
}

func init() {
	initFlags(rootCmd)
}

func main() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		var exitErr *errorx.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.Code())
		}
		os.Exit(1)
	}
}
