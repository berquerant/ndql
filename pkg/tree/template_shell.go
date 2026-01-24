package tree

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/berquerant/ndql/pkg/cachex"
	"github.com/berquerant/ndql/pkg/util"
)

//
// shell gen template
//
// sh(script)
//
// ## Environment variables
// Available, referenced by $NAME
//
// ## Functions
// ### get
// Get value from a node like 'get key'.
// If key is not found, returns an empty string.
//
// ### get_or
// Get value from a node like 'get key default_value'.
// If key is not found, returns default_value.

type ShellGenTemplate struct {
	text  string
	shell string
}

func NewShellGenTemplate(text string) *ShellGenTemplate {
	return &ShellGenTemplate{
		text:  text,
		shell: "bash",
	}
}

var _ GenTemplate = &ShellGenTemplate{}

var shellGenTemplateCache = util.Must(cachex.NewTmpFileCache(util.TempDir("shell_template")))

func (g ShellGenTemplate) Generate(ctx context.Context, n *N) ([]byte, error) {
	t, err := shellGenTemplateCache.Get(g.generateScript())
	if err != nil {
		return nil, fmt.Errorf("%w: cannot get shell template", errors.Join(ErrGenTemplate, err))
	}

	slog.Debug("ShellGenTemplate", slog.String("file", t))
	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, g.shell, t)
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	cmd.Env = NodeAsEnviron(n)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%w: failed to run shell template %s", errors.Join(ErrGenTemplate, err), t)
	}
	return bytes.TrimSpace(out.Bytes()), nil
}

func (g ShellGenTemplate) generateScript() string {
	return strings.Join([]string{
		fmt.Sprintf("#!/bin/%s", g.shell),
		fmt.Sprintf(shellGenTemplateFunc, TableKeySeparator),
		g.text,
	}, "\n")
}

const shellGenTemplateFunc = `key_has_table() {
  echo "$1" | grep -q '%[1]s'
}
key_from_name() {
  echo "$1" | sed 's|\.|%[1]s|g'
}
name_from_key() {
  echo "$1" | sed 's|%[1]s|\.|g'
}
key_suffix() {
  name_from_key "$1" | cut -d "." -f 2-
}
get() {
  local -r name="$1"
  local -r key="$(key_from_name "$name")"
  if key_has_table "$key" ; then
    echo "${!key}"
    return
  fi
  local -r suffix="$(key_suffix "$key")"
  for varname in $(compgen -v | grep -E "%[1]s${suffix}$") ; do
    if ! key_has_table "$varname" ; then
      continue
    fi
    echo "${!varname}"
    return
  done
  echo "${!suffix}"
}
get_or() {
  local -r name="$1"
  local -r default_value="$2"
  local r="$(get "$name")"
  if [[ "$r" == "" ]] ; then
    echo "$default_value"
  else
    echo "$r"
  fi
}`
