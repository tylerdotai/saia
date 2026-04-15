# SAiA — Simple AI Agent

A lightweight, self-hosted AI agent runtime built in Go. Think OpenClaw stripped to the bone — always-on, multi-platform messaging, persistent memory, skill invocation, and shell access.

## Features

- **Go daemon** — compiled binary, no runtime dependencies, minimal memory footprint
- **Discord + Telegram** — messaging platforms as the primary interface
- **TUI REPL** — interactive terminal session when you want to go direct
- **SQLite + FTS5** — persistent memory and full-text search across sessions
- **Skills system** — YAML registry + SKILL.md format, agent can invoke and chain skills
- **Shell exec** — run commands on the host, log everything
- **MCP client** — connect to external MCP servers
- **Docker ready** — binary on the host, or containerized with Docker

## Quick Start

```bash
# Clone
git clone https://github.com/tylerdotai/saia.git ~/saia
cd ~/saia

# Install (requires Go 1.22+)
go mod tidy
make build

# Configure
cp config.json.example ~/.saia/config.json
# Edit ~/.saia/config.json with your Discord bot token + model API key

# Run
./build/saiad
```

Or with Docker:

```bash
docker pull saia/saia:latest
docker run -v ~/.saia:/home/saia/.saia saia/saia:latest
```

## Architecture

```
cmd/saiad/     — Daemon: Discord bot + gateway + always-on process
cmd/saia/      — CLI: TUI REPL for direct terminal interaction
internal/
  agent/       — Core agent loop
  config/      — JSON config loading + env var overrides
  db/          — SQLite schema + FTS5 queries
  discord/     — Discord gateway client
  exec/        — Shell command executor
  health/      — HTTP healthcheck server
  logging/     — Structured JSON logging
  memory/      — FTS5-backed memory store
  mcp/         — MCP client connections
  skills/      — Skill registry + invocation engine
  tui/         — Terminal REPL UI
```

## Phase Roadmap

| Phase | Goal |
|---|---|
| **Phase 1** | Foundation: Discord daemon, TUI, SQLite, shell exec, skills |
| **Phase 2** | Multi-platform (Telegram, Signal), MCP client, Docker, memory layer |
| **Phase 3** | Self-improving skills (agent writes its own), SAiA as MCP server |

## Build Rules

Every commit must pass:

```bash
make fmt    # go fmt
make lint   # golangci-lint
make tidy   # go mod tidy
make build  # both binaries
make test   # tests + race detector
```

Minimum 70% code coverage required.

## Docs

- [SPEC.md](SPEC.md) — Full specification
- [CONTRIBUTING.md](CONTRIBUTING.md) — Development guide
- [CONFIG.md](CONFIG.md) — Configuration reference

## License

MIT
