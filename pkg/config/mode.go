package config

import (
	"errors"
	"fmt"
	"log/slog"
)

var ErrUnknownMode = errors.New("UnknownMode")

// Mode specifies the execution mode for ndql.
type Mode string

const (
	StringModeUnknown = "unknown"
	StringModeVersion = "version"
	StringModeQuery   = "query"
	StringModeDryrun  = "dryrun"
	StringModeList    = "list"
)

const (
	ModeUnknown = Mode(StringModeUnknown)
	ModeVersion = Mode(StringModeVersion)
	ModeQuery   = Mode(StringModeQuery)
	ModeDryrun  = Mode(StringModeDryrun)
	ModeList    = Mode(StringModeList)
)

func NewMode(v string) Mode {
	switch v {
	case StringModeVersion:
		return ModeVersion
	case StringModeQuery:
		return ModeQuery
	case StringModeDryrun:
		return ModeDryrun
	case StringModeList:
		return ModeList
	default:
		return ModeUnknown
	}
}

func (m Mode) String() string { return string(m) }

func (c *Config) SetupSources() error {
	slog.Debug("Setup sources")
	s, err := c.newSources(c.Args)
	if errors.Is(err, ErrNoSources) {
		return nil
	}
	if err != nil {
		return err
	}
	if s.HasStdinConflict() {
		return ErrStdinConflict
	}
	c.Sources = s
	return nil
}

func (c Config) newSources(args []string) (*Sources, error) {
	switch c.Mode {
	case ModeVersion:
		return nil, ErrNoSources
	case ModeDryrun:
		return c.newDryrunSources(args)
	case ModeList:
		return c.newListSources(args)
	case ModeQuery:
		return c.newQuerySources(args)
	default:
		return nil, ErrUnknownMode
	}
}

func (c Config) newDryrunSources(args []string) (*Sources, error) {
	switch len(args) {
	case 1:
		querySource := args[0]
		query, err := NewQuerySource(querySource)
		if err != nil {
			return nil, err
		}
		return &Sources{
			Query: query,
		}, nil
	default:
		return nil, fmt.Errorf("%w: no query", ErrInvalidSource)
	}
}

func (c Config) newListSources(args []string) (*Sources, error) {
	switch len(args) {
	case 1:
		pathSource := args[0]
		walker, path, err := NewPathSource(pathSource)
		if err != nil {
			return nil, err
		}
		return &Sources{
			Path: walker,
			path: path,
		}, nil
	default:
		return nil, fmt.Errorf("%w: no path", ErrInvalidSource)
	}
}

func (c Config) newQuerySources(args []string) (*Sources, error) {
	switch len(args) {
	case 1:
		if c.Index == "" {
			return nil, fmt.Errorf("%w: no path no index", ErrInvalidSource)
		}
		querySource := args[0]
		query, err := NewQuerySource(querySource)
		if err != nil {
			return nil, err
		}

		slog.Debug("Use index", slog.String("source", c.Index))
		index, err := NewSource(c.Index)
		if err != nil {
			return nil, err
		}
		return &Sources{
			Query: query,
			Index: index,
		}, nil
	case 2:
		if c.Index != "" {
			return nil, fmt.Errorf("%w: index is exclusive with path", ErrInvalidSource)
		}
		querySource := args[0]
		query, err := NewQuerySource(querySource)
		if err != nil {
			return nil, err
		}

		pathSource := args[1]
		walker, path, err := NewPathSource(pathSource)
		if err != nil {
			return nil, err
		}
		return &Sources{
			Query: query,
			Path:  walker,
			path:  path,
		}, nil
	default:
		return nil, fmt.Errorf("%w: invalid query, path, index", ErrInvalidSource)
	}
}
