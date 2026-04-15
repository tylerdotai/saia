# SAiA — Simple AI Agent

## What It Is

SAiA is a lightweight, self-hosted AI agent runtime built in Go. Think OpenClaw stripped to the bone — always-on, multi-platform messaging, persistent memory, skill invocation, and shell access. Designed for developers who want full control without the overhead.

**Principles:**
- One binary. One config. No runtime dependencies.
- Go for speed and minimal footprint.
- SQLite for persistence. FTS5 for memory search.
- Human-readable config. Developer-friendly onboarding.
- Ship phases. Don't boil the ocean.

---

## Background / Why

Tyler's dissatisfaction with OpenClaw's complexity + a desire to build from scratch = SAiA.

Key problems SAiA solves:
1. **Agents are passive** — they wait to be asked. SAiA's skills system lets the agent take initiative.
2. **Too much bloat** — Electron/Node.js overhead for what should be a lean daemon.
3. **Skills are static** — SKILL.md files are human-written. SAiA keeps the format but builds a proper invocation engine.
4. **No self-improvement** — Planned for Phase 3: the agent writes its own skills from successful runs.

---

## Stack

| Layer | Choice | Rationale |
|---|---|---|
| Language | Go | Compiled, single binary, minimal memory, fast startup |
| Runtime config | JSON | Bot tokens, API keys, database path |
| Declarative config | YAML | Skills registry, agent definitions, routing rules |
| Database | SQLite + FTS5 | Based on Hermes schema; persistent memory, session history |
| Install | bun / npm | Developer-friendly, already on clawbox |
| Package manager | npx / go mod | npm for JS tooling, Go for the core |
| Process manager | systemd | Always-on daemon on clawbox |

---

## Project Structure

```
saia/                          # Project code (~/saia/)
├── cmd/
│   ├── saiad/                 # Daemon entry point (discord bot + gateway)
│   └── saia/                  # CLI entry point (TUI REPL)
├── internal/
│   ├── agent/                 # Core agent loop
│   ├── config/                # Config loading (JSON + YAML + env overrides)
│   ├── db/                    # SQLite + FTS5 schema + queries
│   ├── discord/               # Discord bot client + rate limiter
│   ├── exec/                  # Shell command executor
│   ├── health/                # Healthcheck HTTP server
│   ├── logging/               # Structured JSON logging + rotation
│   ├── memory/                # Memory management
│   ├── mcp/                   # MCP client connections
│   ├── skills/                # Skill registry + invocation
│   └── tui/                   # Terminal REPL UI
├── .github/
│   └── workflows/
│       ├── ci.yml             # Lint + test + coverage
│       └── release.yml        # Cross-compile + Docker build + push on tag
├── Dockerfile                 # Multi-stage build, distroless base (~10MB)
├── docker-compose.yml         # Local dev: SAiA + SearXNG + Redis sidecars
├── skills/                    # Bundled skills (SKILL.md format)
├── .saia/                     # Runtime data (gitignored)
│   ├── saia.db               # SQLite database
│   ├── skills/               # User-added skills
│   └── logs/                 # Daemon logs
├── Makefile                   # make build|test|lint|install|uninstall
├── CONFIG.md                  # Configuration reference
├── README.md
├── CONTRIBUTING.md
└── SPEC.md
```

**Runtime data at:** `~/.saia/` (not `~/.saia/` — dot prefix keeps it hidden)
**Project code at:** `~/saia/`

---

## Phase Roadmap

### Phase 1 — Foundation (THIS SPEC)
**Goal:** Minimal viable daemon. Discord connected. TUI works. Shell exec works. SQLite initialized.

- [ ] Go project scaffold (`~/saia/`)
- [ ] `saiad` daemon — starts on boot via systemd
- [ ] `saia` CLI — opens TUI REPL
- [ ] Discord bot client (new bot token, guild commands)
- [ ] SQLite schema (sessions, messages, memory, skills registry)
- [ ] Shell exec (local commands, full access on host)
- [ ] Config loading (JSON for tokens/keys, YAML for skills/routing)
- [ ] Basic skill invocation (YAML registry → SKILL.md → execute)

### Phase 2 — Expansion
**Goal:** Multi-platform, better memory, MCP connections.

- [ ] Telegram adapter (in addition to Discord)
- [ ] Signal, WhatsApp (roadmap, not committed)
- [ ] Memory layer — FTS5 search, LLM summarization
- [ ] MCP client — connect to external MCP servers
- [ ] Skill self-registration — agent can add new skills to registry
- [ ] Persistent shell sessions (tmux backend)
- [ ] Docker support — `Dockerfile` + `docker-compose.yml`
  - **Host mode (default):** binary runs on host → full shell exec on host
  - **Container mode:** toggled via `config.json` → `exec.backend: "host" | "docker"`
  - Docker users get `docker pull saia/saia:latest` + `docker-compose.yml` for sidecars

### Phase 3 — Intelligence
**Goal:** Agent writes its own skills. Self-improving loop.

- [ ] Self-improvement loop — successful runs → SKILL.md generation
- [ ] Skill evaluation — agent assesses skill quality, refines or deletes
- [ ] Agent-as-MCP-server — other tools can call SAiA
- [ ] Broader model support (OpenRouter, Ollama, etc.)

---

## Database Schema (SQLite + FTS5)

Based on Hermes Agent's schema, simplified for Phase 1.

```sql
-- Sessions: one per conversation thread
CREATE TABLE sessions (
  id TEXT PRIMARY KEY,
  platform TEXT NOT NULL,          -- discord, telegram, cli
  platform_id TEXT NOT NULL,       -- channel/chat ID
  model TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  metadata TEXT                    -- JSON: thread info, user info
);

-- Messages: full conversation history
CREATE TABLE messages (
  id TEXT PRIMARY KEY,
  session_id TEXT NOT NULL REFERENCES sessions(id),
  role TEXT NOT NULL,              -- user, assistant, system, tool
  content TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  metadata TEXT                    -- JSON: tool calls, errors, timing
);

-- Memory: FTS5-indexed persistent memory
CREATE VIRTUAL TABLE memory USING fts5(
  id UNINDEXED,
  content,
  context_id,
  created_at
);

CREATE TABLE memory_meta (
  id TEXT PRIMARY KEY,
  context_id TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  last_accessed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  access_count INTEGER DEFAULT 0
);

-- Skills registry
CREATE TABLE skills (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  description TEXT,
  trigger TEXT,                    -- keywords or intent pattern
  file_path TEXT NOT NULL,        -- path to SKILL.md
  enabled INTEGER DEFAULT 1,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  last_used_at DATETIME
);

-- Skill invocations (for self-improvement tracking in Phase 3)
CREATE TABLE skill_invocations (
  id TEXT PRIMARY KEY,
  skill_id TEXT NOT NULL REFERENCES skills(id),
  session_id TEXT REFERENCES sessions(id),
  success INTEGER,
  output TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Config metadata
CREATE TABLE meta (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL
);
```

---

## Config Format

### `~/.saia/config.json` — Runtime Config
```json
{
  "discord": {
    "token": "YOUR_DISCORD_BOT_TOKEN",
    "guild_id": "YOUR_GUILD_ID",
    "allowed_channels": ["channel-id-1", "channel-id-2"]
  },
  "model": {
    "provider": "minimax",
    "api_key": "YOUR_MINIMAX_API_KEY",
    "endpoint": "https://api.minimax.chat/v1",
    "model": "MiniMax-M2.7"
  },
  "database": {
    "path": "~/.saia/saia.db"
  },
  "exec": {
    "shell": "/bin/bash",
    "cwd": "/home/dexter",
    "timeout_seconds": 30
  },
  "mcp": {
    "servers": []
  },
  "logging": {
    "level": "info",
    "path": "~/.saia/logs/saiad.log"
  }
}
```

### `~/.saia/skills.yaml` — Skills Registry
```yaml
skills:
  - name: web_search
    description: Search the web using SearXNG
    trigger: "search|look up|find on web"
    file_path: "~/.saia/skills/web_search.md"
    enabled: true

  - name: shell_exec
    description: Execute shell commands on the host
    trigger: "run|execute|shell|command"
    file_path: "~/.saia/skills/shell_exec.md"
    enabled: true

  - name: memory_recall
    description: Search long-term memory
    trigger: "remember|what do you know|recall"
    file_path: "~/.saia/skills/memory_recall.md"
    enabled: true

  - name: skill_list
    description: List all available skills
    trigger: "what can you do|list skills|help"
    file_path: "~/.saia/skills/skill_list.md"
    enabled: true
```

### Skill File Format (SKILL.md)
```
# Skill: web_search

## Description
Searches the web using the local SearXNG instance.

## Triggers
"search for X", "look up X", "find X on web"

## Actions
1. Extract search query from user message
2. Call SearXNG API (localhost:8888)
3. Return top 5 results as formatted text

## Output Format
Formatted markdown with titles, URLs, and snippets.

## Notes
Uses local SearXNG instance. Falls back to web_fetch if SearXNG is down.
```

---

## Phase 1 Build Order

### Step 1: Project Scaffold
```
~/saia/
├── cmd/saiad/main.go
├── cmd/saia/main.go
├── internal/
│   ├── config/
│   ├── db/
│   ├── agent/
│   ├── discord/
│   ├── tui/
│   ├── exec/
│   ├── skills/
│   └── memory/
├── go.mod
├── go.sum
└── README.md
```

### Step 2: Config + Database
- Load `~/.saia/config.json` on startup
- Initialize SQLite at `~/.saia/saia.db`
- Create tables if they don't exist

### Step 3: Discord Bot Client
- Connect to Discord gateway with bot token
- Handle message events (create, update, delete)
- Respond to mentions or DMs
- Register slash commands: `/saia prompt <message>`, `/saia status`

### Step 4: TUI REPL
- `saia` command opens interactive terminal UI
- Persistent input loop
- Colored output, session context shown
- Escape / Ctrl+C to exit

### Step 5: Agent Loop
- Load skills from `~/.saia/skills/`
- Match incoming message against skill triggers
- If skill matched → execute skill action
- If no skill matched → call LLM with system prompt + memory context
- Stream response back to platform (Discord or TUI)

### Step 6: Shell Exec
- Execute shell commands on local host
- Timeout enforcement
- Return stdout/stderr as tool result
- Security: log all executions

### Step 7: Skills Registry
- Load `~/.saia/skills.yaml`
- Watch for changes (inotify or poll)
- SKILL.md parser (frontmatter + action steps)
- Skill invocation engine

### Step 8: Systemd Unit
```ini
[Unit]
Description=SAiA Daemon
After=network.target

[Service]
Type=simple
User=dexter
ExecStart=/home/dexter/saia/build/saiad
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

---

## API Design (Internal)

### Core Agent Interface
```go
type Agent interface {
  Start(ctx context.Context) error
  Stop() error
  Prompt(ctx context.Context, p Prompt) (*Response, error)
  GetSkills() []Skill
  InvokeSkill(ctx context.Context, skillID string, input string) (*SkillResult, error)
}
```

### Prompt Structure
```go
type Prompt struct {
  SessionID  string
  Platform   string    // discord, telegram, cli
  ChannelID  string
  UserID     string
  Content    string
  Role       string    // user, system
  Metadata   map[string]any
}
```

### Response Structure
```go
type Response struct {
  Content    string
  SkillUsed  string
  ToolCalls  []ToolCall
  SessionID  string
  Metadata   map[string]any
}
```

---

## Quality & Operations

### Build Rules (Pre-Commit Gate)

Every commit must pass before merging to main:

```bash
# 1. Format
go fmt ./...

# 2. Lint
golangci-lint run ./...

# 3. Dependencies clean
go mod tidy

# 4. Compile
go build ./...

# 5. Test + race detector
go test ./... -race -count=1

# 6. Coverage gate (minimum 70%)
go test ./... -coverprofile=coverage.out -covermode=atomic
go tool cover -func=coverage.out | grep total
# Must be >= 70%
```

### Makefile Targets

```makefile
make build        # Build both saiad and saia binaries
make test         # Run all tests with race detector
make lint          # Run golangci-lint
make coverage      # Generate coverage report
make install       # Install binaries to /usr/local/bin
make uninstall     # Remove binaries
make clean         # Remove build artifacts
make run           # Run saiad from source (dev mode)
make docs          # Generate code docs
```

### GitHub Actions CI

- **On every PR:** lint → build → test → coverage report → comment on PR
- **On every tag:** cross-compile Linux amd64 + arm64 → create GitHub Release → upload binaries
- **Branch protection:** require CI green before merge to main

### Error Response Format

Discord errors must never expose raw panic traces. All errors are wrapped:

```go
type ErrorResponse struct {
  Code    string `json:"code"`     // e.g., "SHELL_TIMEOUT", "MODEL_UNAVAILABLE"
  Message string `json:"message"`   // Human-safe description
  Retry   bool   `json:"retry"`    // Whether the user can retry
}
```

- Message send failures → log + skip (Discord handles retry)
- Shell timeouts → reply with "Command timed out after 30s"
- Model errors → reply with "I'm having trouble thinking right now, try again"
- Panic → log full stack, reply with "Something broke on my end — logged for review"

### Healthcheck Endpoint

Required for systemd watchdog + monitoring:

```
GET /health → 200 OK { "status": "ok", "version": "0.1.0", "uptime": "2h34m" }
GET /health/ready → 200 OK { "discord": "connected", "db": "ok", "model": "ok" }
```

Port: `8080` by default (configurable via `config.json`)

### Graceful Shutdown

`saiad` must handle SIGTERM and SIGINT cleanly:

1. Stop accepting new messages
2. Finish any in-flight model calls (with 10s timeout)
3. Flush SQLite writes
4. Close Discord gateway connection
5. Exit with code 0

```go
// Signal handling
signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
<-sigChan
// graceful shutdown sequence
```

### Logging

- Format: JSON to file + plain text to stdout (journalctl readable)
- Levels: `debug`, `info`, `warn`, `error` (configurable)
- Log files rotated daily, max 7 days (delete oldest)
- All shell execs logged: timestamp + user + command + duration + exit code
- Location: `~/.saia/logs/saiad.log`

### Rate Limiting

- Discord: respect `X-RateLimit-*` headers, never burst past limits
- Model API: configurable requests-per-minute cap
- Shell exec: max 10 concurrent processes (configurable)

### Environment Variable Overrides

Config values can be overridden by env vars (useful for CI, containers):

```bash
SAIA_DISCORD_TOKEN=...       # Overrides discord.token
SAIA_MODEL_API_KEY=...       # Overrides model.api_key
SAIA_DB_PATH=...             # Overrides database.path
SAIA_LOG_LEVEL=debug         # Overrides logging.level
SAIA_PORT=8081               # Overrides healthcheck port
```

### Version Flag

Both `saiad` and `saia` support `--version`:

```
$ saia --version
saia version 0.1.0 (go1.22.2)
$ saiad --version
saiad version 0.1.0 (go1.22.2)
```

Git tag `v0.1.0` sets the version at build time via `-ldflags`.

### Message Retry

Discord message sends that fail (rate limit, network) are retried up to 3 times with exponential backoff (1s, 2s, 4s). After 3 failures, the error is logged and the message is skipped.

---

## Security Considerations

- All shell executions logged to `~/.saia/logs/exec.log`
- No privilege escalation — daemon runs as `dexter` user, not root
- Discord allowed_channels restricts which channels accept commands
- API keys stored in `config.json`, never committed to git
- `.saia/` directory in `.gitignore`
- Optional: skill sandboxing (Docker exec for Phase 2+)

---

## Developer Onboarding

```bash
# Clone
git clone https://github.com/tylerdotai/saia.git ~/saia
cd ~/saia

# Install dependencies
bun install

# Build
go build -o build/saiad ./cmd/saiad
go build -o build/saia ./cmd/saia

# Configure
cp ~/.saia/config.json.example ~/.saia/config.json
# Edit config.json with your Discord bot token + MiniMax API key

# Run daemon
./build/saiad

# Or TUI mode
./build/saia
```

---

## What's NOT in Phase 1

- Multi-platform (Telegram, Signal, WhatsApp) — Phase 2
- MCP server (SAiA as an MCP server for other tools) — Phase 3
- Self-improving skills — Phase 3
- Web UI / dashboard — never (CLI + TUI only)
- A2A protocol — out of scope
- Docker image building + Docker Hub / GHCR.io publish — Phase 2+

---

## Naming Rationale

**SAiA** — Simple AI Agent. Honest. No mythology, no metaphors. It does what it says. Short enough to type, recognizable, not already taken.

---

## Gaps Identified (and Resolved)

| Gap | Resolution |
|---|---|
| No Makefile | Added `make` targets |
| No CI/CD | GitHub Actions workflow on PR + tag |
| Raw panic traces on Discord | ErrorResponse struct + graceful wrapping |
| No healthcheck | HTTP server on port 8080 with `/health` + `/health/ready` |
| No graceful shutdown | SIGTERM/SIGINT handler with in-flight drain |
| No logging strategy | JSON file logs + stdout, daily rotation, 7-day retention |
| No rate limit handling | Respect Discord headers + model API RPM cap |
| No env var overrides | `SAIA_*` prefix env vars override JSON config |
| No version flag | `--version` via `-ldflags` from git tag |
| No message retry | Exponential backoff (3 retries) on Discord send failures |
| No CONTRIBUTING.md | Placeholder added |
| Binary cross-compile | GitHub Actions builds amd64 + arm64 on tag |

---

## Status

**Phase 1 — Specification Complete**
- This document defines Phase 1 scope
- Build order established
- Schema designed
- Config format defined
- Quality gates documented
- Operations concerns addressed

**Ready to build.**
