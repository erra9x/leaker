# CLAUDE.md — Senior Engineer Rules for leaker

You are a senior Go engineer working on `leaker`, a CLI tool for querying leak databases.

## Non-Negotiable Rules

- All changes must compile: `go build ./...` must pass before you finish
- All tests must pass: `go test ./...` must pass before you finish
- Run `go vet ./...` — fix all warnings
- Run `golangci-lint run ./...` if available — fix any issues in files you touched
- No demo data, placeholder values, or TODO stubs left in committed code
- No `panic()` in library or runner code — always return errors up the call stack

## Code Style

- Explicit over clever — if it needs a comment to understand, write the comment
- Correct over fast — don't optimise unless there's a measured problem
- Use `errors.Is` / `errors.As` for error inspection, not string matching
- Use `fmt.Errorf("context: %w", err)` to wrap errors with context
- Keep changes minimal and scoped — don't refactor things unrelated to the task
- Match the existing style and naming conventions in surrounding code

## Project Layout

- `cmd/` — CLI parsing (kong), banner, version
- `runner/` — business logic: options, config, sources wiring, enumeration
- `runner/sources/` — source implementations (LeakCheck, ProxyNova, etc.)
- `utils/` — file helpers, env, stdin detection, random pick
- `logger/` — leveled logger (global DefaultLogger)

## Testing

- Tests live next to the code they test (`foo_test.go` beside `foo.go`)
- Use `t.TempDir()` for any file system work in tests — never hardcode paths
- Table-driven tests preferred for functions with multiple input cases
- Tests must not make real network calls — use `httptest.NewServer` for HTTP

## Git

- One logical change per commit
- Commit message format: `fix: <what was wrong and what you did>` or `feat: <what you added>`
- Don't commit `go.sum` changes unless `go.mod` also changed
