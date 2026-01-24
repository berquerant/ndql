package mapx

import (
	"fmt"
	"os"

	"mvdan.cc/sh/v3/shell"
)

type EnvMap = Map[string, string]

type Env struct {
	*EnvMap
}

func NewEnv(m *EnvMap) *Env {
	return &Env{m}
}

func (e *Env) GetOrEnv(key string) string {
	if v, ok := e.Get(key); ok {
		return v
	}
	return os.Getenv(key)
}

func (e *Env) Expand(s string) (string, error) {
	return shell.Expand(s, e.GetOrEnv)
}

func (e *Env) AsEnviron() []string {
	e.mux.RLock()
	defer e.mux.RUnlock()
	var (
		i  int
		xs = make([]string, len(e.m))
	)
	for k, v := range e.m {
		xs[i] = fmt.Sprintf("%s=%s", k, v)
		i++
	}
	return append(os.Environ(), xs...)
}
