# Contributing to SAiA

Thank you for your interest in contributing to SAiA.

## Development Setup

```bash
# Clone the repository
git clone https://github.com/tylerdotai/saia.git
cd saia

# Install Go 1.22+
# Verify: go version

# Install golangci-lint
go install github.com/golangci-lint/cmd/golangci-lint@latest

# Install dependencies
go mod tidy

# Run all pre-commit checks
make verify
```

## Pre-Commit Checklist

Before opening a PR, ensure all of these pass:

```bash
make fmt    # Format code
make lint   # Lint
make tidy   # Tidy go.mod/go.sum
make build  # Build both binaries
make test   # Tests + race detector
```

## Project Structure

- `cmd/saiad/` — Daemon entry point (Discord bot + gateway)
- `cmd/saia/` — CLI/TUI entry point
- `internal/` — Core packages
- `skills/` — Bundled skill definitions

## Style Guide

- Run `go fmt` before committing
- Document all exported functions
- Return errors; don't log and discard
- Log at the appropriate level (debug/info/warn/error)
- Never commit secrets or real credentials

## Testing

- Unit tests for all public package functions
- Integration tests for Discord, database, and skill loading
- Minimum 70% code coverage required
- Run with `make test`

## Commit Messages

Format: `type: short description`

Types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `perf`

Examples:
- `feat: add Telegram adapter`
- `fix: handle Discord rate limit gracefully`
- `docs: update SPEC.md with Phase 2 features`

## Reporting Issues

- Search existing issues before creating a new one
- Include Go version, OS, and error output
- Include minimal reproduction steps
