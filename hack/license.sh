#!/bin/bash

set -e
set -o pipefail

readonly d="$(cd "$(dirname "$0")" || exit 1 ; pwd)"

go_licenses() {
    "${d}/../tools/run.sh" go-licenses "$@"
}

readonly ignore='--ignore "github.com/berquerant/ndql/"'
readonly target='./cmd/ndql'

report() {
    go_licenses report "$target" $ignore --template="${d}/notice-template.md"
}

check() {
    go_licenses check "$target" $ignore
}

readonly cmd="$1"
case "$cmd" in
    report) report ;;
    check) check ;;
    *)
        echo 'Available command: report,check'
        exit 1
        ;;
esac
