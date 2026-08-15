# AGENTS.md — ndql

## Project Overview

**ndql** is a Go CLI tool that lets you query filesystem metadata using SQL syntax. It walks a directory tree, exposes each file/directory as a row (with built-in columns such as `path`, `size`, `is_dir`, `mod_time`, `mode`), and evaluates SQL `SELECT` / `WHERE` expressions against those rows, producing JSON output.

Module: `github.com/berquerant/ndql`  
Go version: see `go.mod`

---

## Repository Layout

```
.
├── cmd/
│   ├── ndql/          # Main CLI entry point
│   │   ├── main.go    # Root cobra command & exit-code handling
│   │   ├── query.go   # `ndql query` subcommand
│   │   ├── dryrun.go  # `ndql dryrun` subcommand (parse + plan only)
│   │   ├── explain.go # `ndql explain` subcommand (embedded docs)
│   │   ├── list.go    # `ndql list` subcommand
│   │   ├── version.go # `ndql version` subcommand
│   │   ├── flag.go    # Shared flag definitions
│   │   └── docs.json  # Embedded documentation (generated)
│   └── gendoc/        # Tool that generates docs.json and docs/
├── pkg/
│   ├── config/        # Config struct, execution modes, source I/O setup
│   ├── node/          # Node (row) type and built-in column definitions
│   ├── tree/          # Core query execution: SQL AST → iterator pipeline
│   │   ├── func.go    # All built-in SQL functions (grep, sh, lua, expr, …)
│   │   └── node.go    # Key/column helpers and node-mapping utilities
│   ├── parse/         # SQL parsing (wraps TiDB parser)
│   ├── run/           # High-level run() helpers wiring config → executor
│   ├── iterx/         # Generic iterator combinators (map, reduce, fanout, …)
│   ├── iox/           # Filesystem walker (WalkerEntry) and I/O utilities
│   ├── cachex/        # Generic cache utilities
│   ├── errorx/        # ExitError and error helpers
│   ├── gopkg/         # DocumentSet used by `explain`
│   ├── jsonmap/       # JSON-map helpers
│   ├── logx/          # slog-based logger setup
│   ├── mapx/          # Ordered map generic utility
│   ├── regexpx/       # Regexp helpers
│   └── util/          # Miscellaneous small utilities
├── docs/              # Generated reference documentation (Markdown)
│   ├── README.md
│   ├── data/          # Built-in column documentation
│   └── syntax/        # Query syntax / function documentation
├── hack/
│   ├── license.sh     # License header checker / reporter
│   └── notice-template.md
├── bin/               # Compiled binaries (git-ignored)
├── version/           # Version string
├── mise.toml          # Tool versions and task runner definitions
├── go.mod
├── go.sum
├── .golangci.yml
└── .goreleaser.yaml
```

---

## Key Concepts

| Concept | Description |
|---|---|
| **Node** | A single row: a typed key→value map. Built-in keys: `path`, `size`, `is_dir`, `mod_time`, `mode`. |
| **Data / Op** | The two faces of a typed value. `Data` is stored in a Node; `Op` exposes arithmetic/comparison operations. |
| **Key** | A `table.column` or bare `column` reference resolved from a Node. |
| **Source** | Input to a query: a filesystem path, stdin (`@stdin` / `@-`), or an index file of pre-serialised Nodes. |
| **Mode** | Execution mode: `query`, `dryrun`, `list`, `version`. |

---

## CLI Subcommands

| Subcommand | Source file | Purpose |
|---|---|---|
| `ndql query QUERY [PATH]` | `cmd/ndql/query.go` | Run a SQL query against a directory tree or index |
| `ndql dryrun QUERY` | `cmd/ndql/dryrun.go` | Parse and plan a query without executing it |
| `ndql list PATH` | `cmd/ndql/list.go` | List all files/dirs under PATH as JSON |
| `ndql explain [KEY]` | `cmd/ndql/explain.go` | Show embedded documentation |
| `ndql version` | `cmd/ndql/version.go` | Print version |

### Common Flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--concurrency` | `-c` | `0` (→1) | Max goroutines for query processing |
| `--index` | `-i` | | Index source (exclusive with PATH argument) |
| `--raw` | | `false` | Raw (non-JSON) output |
| `--debug` / `--trace` | | `false` | Log verbosity |
| `--verbose` | `-v` | `false` | Verbose output |
| `--quiet` | `-q` | `false` | Suppress all logs except errors |

### Source Notation (QUERY and PATH arguments)

| Notation | Meaning |
|---|---|
| `@stdin` or `@-` | Read from stdin |
| `@FILENAME` | Read from the given file |
| anything else | Use the string value as-is |

---

## Build & Development

### Build the binary

```sh
mise run build         # generates docs then builds bin/ndql
bin/ndql help
```

The build script is `bin/build.sh`.

### Run tests

```sh
mise run test          # all packages with race detector and coverage
mise run test-tree     # pkg/tree only
```

Direct Go command:

```sh
go test -race -cover ./...
```

### Lint

```sh
mise run lint          # check-licenses + vet + golangci-lint
mise run vet           # go vet ./...
mise run golangci-lint # golangci-lint config verify + run
```

Lint configuration: `.golangci.yml`

### Code generation

```sh
mise run generate          # go generate ./...
mise run generate-docs     # regenerate cmd/ndql/docs.json and docs/
mise run clean-generated   # delete all *_generated.go files
```

Generated files — **do not edit by hand**:
- `cmd/ndql/docs.json` — regenerated by `gendoc gen json`
- `docs/` — regenerated by `gendoc gen files`
- `*_generated.go` — produced by `go generate`

### Third-party licenses

```sh
mise run NOTICE            # regenerate NOTICE file
mise run check-licenses    # verify headers are present and NOTICE is up-to-date
```

---

## Adding a New Built-in SQL Function

1. Implement the function in `pkg/tree/func.go`, following the existing doc-comment pattern:
   ```go
   // - myFunc(arg: Type) -> ReturnType
   ```
2. Register it in the same file's function dispatch table.
3. Add unit tests in `pkg/tree/` (files named `*_test.go`).
4. Regenerate docs: `mise run generate-docs`.
5. Verify: `mise run test && mise run lint`.

---

## Adding a New CLI Subcommand

1. Create `cmd/ndql/<subcommand>.go` with a `cobra.Command` and an `init()` that registers it on `rootCmd`.
2. Add the corresponding `Mode` constant to `pkg/config/mode.go` and a `newXxxSources` method on `Config`.
3. Implement the execution logic under `pkg/run/`.
4. Run `mise run generate-docs` if the command should appear in embedded docs.

---

## Testing Guidelines

- Tests live alongside their packages (e.g. `pkg/tree/func_test.go`).
- Use `github.com/stretchr/testify` for assertions.
- `pkg/tree/common_test.go` contains shared test helpers for the `tree` package.
- Always run with the race detector: `go test -race ./...`.
- Linting runs without test files (`tests: false` in `.golangci.yml`).

---

## CI

GitHub Actions workflow: `.github/workflows/test.yml`  
Triggers: push / pull request  
Steps: `go test -race -cover ./...`

---

## Release

Managed by GoReleaser: `.goreleaser.yaml`  
Pinned tool versions are managed through `mise.toml`.

