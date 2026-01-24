package cachex

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"text/template"

	"github.com/yuin/gopher-lua/parse"

	"github.com/berquerant/cache"
	"github.com/berquerant/ndql/pkg/util"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	lua "github.com/yuin/gopher-lua"
)

// A LRU in-memory cache.
type Cache[T any] struct {
	c *cache.LRU[string, T]
}

func (c *Cache[T]) Get(key string) (T, error) {
	return c.c.Get(key)
}

func NewCache[T any](f func(string) (T, error)) (*Cache[T], error) {
	x, err := cache.NewLRU(100, f)
	if err != nil {
		return nil, err
	}
	return &Cache[T]{
		c: x,
	}, nil
}

func NewTextTemplateCache(fm template.FuncMap) (*Cache[*template.Template], error) {
	return NewCache(func(text string) (*template.Template, error) {
		return template.New(util.Hash(text)).Funcs(fm).Parse(text)
	})
}

// Return a new cache for temporary files to skip recreating files with the same content.
func NewTmpFileCache(root string) (*Cache[string], error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, err
	}
	return NewCache(func(text string) (string, error) {
		name := util.Hash(text)
		p := filepath.Join(root, name)
		f, err := os.Create(p)
		if err != nil {
			return "", err
		}
		defer func() {
			_ = f.Close()
		}()
		if _, err := fmt.Fprint(f, text); err != nil {
			return "", err
		}
		return p, nil
	})
}

func NewExprCache() (*Cache[*vm.Program], error) {
	return NewCache(func(text string) (*vm.Program, error) {
		return expr.Compile(text, expr.AsKind(reflect.String))
	})
}

func NewLuaCache() (*Cache[*lua.FunctionProto], error) {
	return NewCache(func(text string) (*lua.FunctionProto, error) {
		b := bytes.NewBufferString(text)
		chunk, err := parse.Parse(b, "script.lua")
		if err != nil {
			return nil, err
		}
		return lua.Compile(chunk, "script.lua")
	})
}
