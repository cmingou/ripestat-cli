# Repository Guidelines

## Project Structure & Module Organization
Core CLI entrypoint lives in `main.go`, with Cobra commands in `cmd/`. API orchestration resides under `internal/ripestat/`, while address classification and formatting helpers sit in `internal/utils/`. Build artifacts land in `build/`, documentation assets in `docs/`, and the smoke-test script is `test_cli.sh`.

## Build, Test, and Development Commands
Use `go build -o ripestat main.go` for a quick local binary; `make install` publishes it to your `$GOBIN`. Cross-platform bundles come from `make all` or `make darwin|linux`, which drop versioned binaries in `build/`. Run `./test_cli.sh` for a fast end-to-end check, `go test ./... -v` for full coverage, `make test` when you prefer the Makefile wrapper, and `make lint` to execute `golangci-lint`.

## Concurrency Controls
RIPEstat allows at most eight concurrent requests per source IP. The CLI defaults to that ceiling and offers both a `--max-concurrency` flag and a `RIPESTAT_MAX_CONCURRENCY` environment variable to tune it (values are clamped between 1 and 8).

## Coding Style & Naming Conventions
Always format with `go fmt ./...`; the project uses default Go tab indentation and imports grouped by `goimports`. Follow Go naming: exported identifiers in CamelCase with doc comments, helpers in lowerCamelCase. Table headers and CLI labels should mirror existing terminology (`COUNTRY`, `IN BGP`, etc.) to keep output consistent.

## Testing Guidelines
Unit coverage lives beside code (`*_test.go` within `cmd/` and `internal/`). End-to-end CLI behavior is asserted in `cmd/integration_test.go` and via `test_cli.sh`. When adding new RIPEstat calls, provide integration tests under `internal/ripestat/` and ensure `go test ./... -cover` clears 80%+ for the touched packages. Name tests `Test<Subject><Behavior>` for clarity.

## Commit & Pull Request Guidelines
Commit messages follow the imperative, sentence-style pattern seen in history (e.g., `Add comprehensive unit tests for input validation`). Squash noisy fixups locally before opening a PR. Each PR should describe intent, list validation steps (`go test ./...`, `make lint`), link related issues, and attach CLI output or screenshots when UI-facing changes alter table rendering. Request review before merging and ensure new binaries are not checked in.

## Security & Configuration Tips
The CLI contacts public RIPEstat endpoints only; no secrets are stored. Guard against accidental rate-limit regressions by running tests with live network access sparingly and mocking responses when feasible.
