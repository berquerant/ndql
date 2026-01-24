package run

import (
	"context"
	"log/slog"

	"github.com/berquerant/ndql/pkg/config"
	"github.com/berquerant/ndql/pkg/logx"
)

type runner struct {
	*config.Config
}

func Main(ctx context.Context, c *config.Config) error {
	slog.Debug("RunMain", logx.String("mode", c.Mode), logx.Value("args", c.Args))
	defer func() {
		_ = c.Close()
	}()
	r := &runner{c}
	switch c.Mode {
	case config.ModeVersion:
		return r.version()
	case config.ModeQuery:
		return r.query(ctx)
	case config.ModeDryrun:
		return r.dryrun()
	case config.ModeList:
		return r.list()
	}
	return r.dryrun()
}
