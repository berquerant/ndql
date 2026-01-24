#!/bin/bash

set -e
set -o pipefail

readonly name="$1"
shift

d="$(cd "$(dirname "$0")" || exit 1; pwd)"
readonly bind="${d}/../bin/tools"
mkdir -p "$bind"
readonly bin="${bind}/${name}"

if [[ ! -x "$bin" ]] ; then
    cmd=""
    case "$name" in
        golangci-lint) cmd="github.com/golangci/golangci-lint/v2/cmd/golangci-lint" ;;
        go-licenses) cmd="github.com/google/go-licenses/v2" ;;
        stringer) cmd="golang.org/x/tools/cmd/stringer" ;;
        govulncheck) cmd="golang.org/x/vuln/cmd/govulncheck" ;;
        *)
            echo >&2 "Unknown tool: ${name}"
            exit 1
            ;;
    esac
    go -C "$d" build -o "$bin" "$cmd"
fi

"$bin" "$@"
