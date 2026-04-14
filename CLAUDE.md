# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Pray CLI is a Go CLI tool for Islamic prayer times. It fetches data from the `pray.ahmedelywa.com` API, supports 23+ calculation methods, multiple output formats (table, pretty, JSON, Slack, Discord), calendar ICS generation, and IP-based location detection.

## Build & Development Commands

```bash
make build              # Build binary to bin/pray
make install            # Install globally via go install
make test               # Run all tests (go test -v ./...)
make test-coverage      # Tests with coverage report (coverage.html)
make lint               # golangci-lint (falls back to go vet)
make fmt                # go fmt ./...
make tidy               # go mod tidy
make run                # Run without installing
make run-args ARGS="--help"  # Run with arguments
make dev-setup          # Install deps + golangci-lint
```

Run a single test: `go test -v -run TestName ./internal/api/`

## Architecture

**Entry point**: `cmd/pray/main.go` sets version info (injected via ldflags) and calls `cmd.Execute()`.

**Command layer** (`cmd/pray/cmd/`): Cobra commands. `root.go` defines global flags, config loading via Viper, and the update checker. Each file is one command (today, next, countdown, diff, get, calendar, config, cache, methods, init, version, completion). The default command (no args) runs `today`.

**Internal packages** (`internal/`):
- `api/` - HTTP client with retry logic (`client.go`), cached wrapper (`cached_client.go`), request params builder, response types, and validation
- `config/` - Config struct, loader (reads/writes `~/.config/pray/config.yaml`), defaults, validation
- `location/` - IP-based geolocation detector using ip-api.com
- `output/` - Formatter interface with implementations: table, pretty, JSON, Slack Block Kit, Discord embeds
- `calendar/` - ICS URL generator, file downloader, subscription instructions
- `cache/` - File-based response caching (`~/.cache/pray/`)
- `ui/` - Spinner and interactive setup wizard
- `update/` - GitHub release version checker

**Public package** (`pkg/prayer/`): Reusable prayer time utilities and calculation method data.

## Key Patterns

- Version info is injected at build time via `-ldflags` (see Makefile `LDFLAGS`)
- Config is managed by Viper, stored as YAML at `~/.config/pray/config.yaml`
- API responses are cached to disk; bypass with `--no-cache`
- Output formatters implement a common interface in `internal/output/formatter.go`
- Releases are automated via GoReleaser (`.goreleaser.yml`) triggered by git tags

## Dependencies

- **cobra** / **viper** - CLI framework and config
- **fatih/color** - Terminal colors
- **olekukonko/tablewriter** - ASCII tables
- **gopkg.in/yaml.v3** - YAML parsing

## Git Workflow

### Commits
All commits must follow [Conventional Commits](https://www.conventionalcommits.org/):
  `type(scope)!: description`

- **Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`
- **Scope** is optional but recommended for targeted changes (e.g., `fix(api): ...`)
- **Breaking changes**: use `!` after type/scope (e.g., `feat!: drop Go 1.22 support`)
- Merge commits are exempt from linting

### Branches
- `main` — protected; all changes land via PRs
- `release/vX.Y.Z` — release preparation branches (e.g., `release/v1.1.0`)

### CI (runs on PRs and pushes to `main`)
1. **Tests** — `go test -v -race` on Go 1.23 and 1.24
2. **Lint** — golangci-lint + `gofmt` check + `go vet`
3. **Build** — cross-compile (linux/darwin/windows, amd64/arm64)
4. **Commitlint** — validates all PR commits follow conventional format

### Releases
Triggered by pushing a `v*` tag (e.g., `git tag v1.2.0 && git push origin v1.2.0`):
- GoReleaser builds binaries for Linux/macOS/Windows (amd64/arm64/arm)
- Checksums signed with cosign
- Changelog auto-generated (excludes docs/test/chore/ci commits)
- GitHub Release created with install instructions

## Module Path

`github.com/AbdElrahmaN31/pray-cli` (Go 1.23+)