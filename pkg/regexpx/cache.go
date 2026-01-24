package regexpx

import (
	"regexp"

	"github.com/berquerant/cache"
	"github.com/berquerant/ndql/pkg/util"
)

var rCache = util.Must(cache.NewLRU(50, regexp.Compile))

func Compile(expr string) (*regexp.Regexp, error) { return rCache.Get(expr) }
