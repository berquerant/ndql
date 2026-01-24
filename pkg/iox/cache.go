package iox

import (
	"log/slog"
	"os"

	"github.com/berquerant/cache"
	"github.com/berquerant/ndql/pkg/logx"
	"github.com/berquerant/ndql/pkg/util"
)

const contentCacheSize = 100

// Cache of contents of files.
type ContentCache interface {
	Get(filename string) ([]byte, error)
}

func NewFileContentCache() *FileContentCache {
	return &FileContentCache{
		c: util.Must(cache.NewLRU(contentCacheSize, func(filename string) ([]byte, error) {
			x, err := os.ReadFile(filename)
			if err != nil {
				logx.Trace("FileContentCache", slog.String("filename", filename), logx.Err(err))
			} else {
				logx.Trace("FileContentCache", slog.String("filename", filename), slog.Int("size", len(x)))
			}
			return x, err
		})),
	}
}

type FileContentCache struct {
	c *cache.LRU[string, []byte]
}

var _ ContentCache = &FileContentCache{}

func (c *FileContentCache) Get(filename string) ([]byte, error) {
	return c.c.Get(filename)
}
