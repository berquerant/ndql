package config

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/berquerant/ndql/pkg/logx"
)

type Config struct {
	Debug       bool `name:"debug" usage:"enable debug logs"`
	Trace       bool `name:"trace" usage:"enable trace logs"`
	Verbose     bool `name:"verbose" short:"v" usage:"enable verbose output"`
	Quiet       bool `name:"quiet" short:"q" usage:"quiet logs except errors"`
	Concurrency uint `name:"concurrency" short:"c" usage:"maximum number of goroutines to process query, 0 means 1"`

	Index     string `name:"index" short:"i" usage:"index source; exclusive with paths"`
	RawOutput bool   `name:"raw" usage:"enable raw output"`

	Mode  Mode     `name:"-"`
	Query string   `name:"-"`
	Path  string   `name:"-"`
	Args  []string `name:"-"`

	Stdout  io.Writer `name:"-"`
	Stderr  io.Writer `name:"-"`
	Sources *Sources  `name:"-"`
}

func (c *Config) Close() error {
	if x := c.Sources; x != nil {
		return x.Close()
	}
	return nil
}

func (c Config) SetupLogger() {
	logx.Setup(c.Stderr, c.Debug, c.Trace, c.Quiet)
}

func (c *Config) SetupQuery() error {
	slog.Debug("Setup query")
	if s := c.Sources; s != nil {
		if q := s.Query; q != nil {
			b, err := q.ReadAll()
			if err != nil {
				return fmt.Errorf("%w: failed to read query", err)
			}
			c.Query = string(b)
		}
	}
	slog.Debug("Query", slog.String("query", c.Query))
	return nil
}
