package main

import (
	"context"
	"errors"
	"os"

	"github.com/berquerant/ndql/pkg/errorx"
	"github.com/berquerant/ndql/pkg/logx"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug logs")
}

var rootCmd = &cobra.Command{
	Use:   "gendoc",
	Short: "Generate documents from comments",
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		logx.Setup(os.Stderr, debug, false, false)
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Help()
	},
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
