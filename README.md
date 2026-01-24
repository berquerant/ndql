# ndql

Select metadata from files by SQL.

## Usage

```
❯ ndql help query
Run query.

## QUERY, PATH
- @stdin or @-: from stdin
- @FILENAME: from file
- otherwise: as it is

## Examples

List all files and directories under the dir:

    ndql query 'select *' dir

List all file paths under the dir:

    ndql query 'select path where not is_dir' dir

Add a "count" column for the number of times 'keyword' appears in the file:

    ndql query 'select sh("echo count=$(grep -c keyword $(get path))") where not is_dir'

Roughly list the Go func in the file except tests:

    ndql query 'select path, func from (select grep("func (?P<name>[^(]+)", "func=$name") where not is_dir and path not like "%_test.go")' dir

Add the line number where each func is defined as 'line':

    ndql query 'select path, func from (select grep("func (?P<name>[^(]+)", "func=$name") where not is_dir and path not like "%_test.go")' dir | ndql query 'select sh("echo line=$(grep -n \"func $(get func)\" $(get path) | cut -d: -f1)")' -i@-

Usage:
  ndql query QUERT [PATH] [flags]

Flags:
      --debug          enable debug logs
  -h, --help           help for query
  -i, --index string   index source; exclusive with paths
      --raw            enable raw output
      --trace          enable trace logs
  -v, --verbose        enable verbose output
```

## Documents

See [docs](./docs/README.md)

## Install

``` shell
make
bin/ndql help
```
