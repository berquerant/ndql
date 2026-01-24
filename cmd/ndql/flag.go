package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/berquerant/ndql/pkg/config"
	"github.com/berquerant/ndql/pkg/run"
	"github.com/berquerant/ndql/pkg/util"
	"github.com/berquerant/structconfig"
	"github.com/spf13/cobra"
)

func initFlags(cmd *cobra.Command) {
	util.FailOnError(structconfig.New[config.Config]().SetFlags(cmd.Flags()))
}

func newConfig(cmd *cobra.Command, args []string) (*config.Config, error) {
	var c config.Config
	if err := structconfig.New[config.Config]().FromFlags(&c, cmd.Flags()); err != nil {
		return nil, err
	}

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.SetupLogger()
	return &c, nil
}

func runMain(cmd *cobra.Command, args []string, mode config.Mode) error {
	c, err := newConfig(cmd, args)
	if err != nil {
		_ = c.Close()
		return err
	}
	c.Args = args
	c.Mode = mode
	ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGPIPE)
	defer stop()
	return run.Main(ctx, c)
}
