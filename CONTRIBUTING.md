# Contributing to Pray CLI

Thank you for your interest in contributing! This guide covers development setup, project structure, testing, and release processes.

## Ways to Contribute

1. **Report Bugs** -- open an issue with details and reproduction steps
2. **Suggest Features** -- share ideas in the [issues](https://github.com/AbdElrahmaN31/pray-cli/issues) section
3. **Submit PRs** -- fork, create a feature branch, and submit a pull request
4. **Improve Docs** -- help improve documentation and examples
5. **Test** -- test on different platforms and report issues

## Development Setup

### Prerequisites

- Go 1.23 or higher
- Make (optional, but recommended)

### Getting Started

```bash
git clone https://github.com/AbdElrahmaN31/pray-cli.git
cd pray-cli
make deps
make build
```

### Build & Run Commands

```bash
make build           # Build binary to bin/pray
make install         # Install globally via go install
make run             # Run without installing
make run-args ARGS="--help"  # Run with arguments
make fmt             # Format code (go fmt)
make lint            # Lint (golangci-lint, falls back to go vet)
make tidy            # go mod tidy
make clean           # Remove build artifacts
```

### Cross-Compilation

```bash
make build-linux-amd64
make build-linux-arm64
make build-darwin-amd64   # macOS Intel
make build-darwin-arm64   # macOS Apple Silicon
make build-windows-amd64
make build-all            # All platforms
```

## Testing

```bash
make test              # Run all tests
make test-coverage     # Tests with coverage report (opens coverage.html)
go test -v -run TestName ./internal/api/   # Run a single test
```

The project includes tests for:
- API client and response parsing
- Location detection and validation
- Configuration management
- Output formatters
- Calendar generation
- Parameter validation

## Project Structure

```
pray-cli/
├── cmd/pray/                  # CLI application
│   ├── main.go                # Entry point (version injected via ldflags)
│   └── cmd/                   # Cobra commands
│       ├── root.go            # Root command, global flags, config loading
│       ├── today.go           # Default command
│       ├── next.go            # Next prayer
│       ├── countdown.go       # Live countdown
│       ├── diff.go            # Location comparison
│       ├── get.go             # Fetch with custom date
│       ├── calendar.go        # Calendar operations
│       ├── config.go          # Configuration management
│       ├── cache.go           # Cache management
│       ├── methods.go         # List calculation methods
│       ├── init.go            # Interactive setup wizard
│       ├── version.go         # Version info
│       └── completion.go      # Shell completions
│
├── internal/                  # Private packages
│   ├── api/                   # HTTP client, caching, types, validation
│   ├── config/                # Config struct, loader, defaults, validation
│   ├── location/              # IP-based geolocation
│   ├── output/                # Formatter interface + implementations
│   ├── calendar/              # ICS generation, download, subscription
│   ├── cache/                 # File-based response caching
│   ├── ui/                    # Spinner, interactive wizard
│   └── update/                # GitHub release version checker
│
├── pkg/prayer/                # Public reusable package (utilities, method data)
├── Makefile                   # Build automation
├── .goreleaser.yml            # Release automation
└── go.mod / go.sum            # Go modules
```

## Key Architecture Patterns

- **Version injection**: Build-time `-ldflags` (see `Makefile` `LDFLAGS`)
- **Config**: Managed by Viper, stored as YAML at `~/.config/pray/config.yaml`
- **Caching**: API responses cached to `~/.cache/pray/`; bypass with `--no-cache`
- **Output**: Formatters implement a common interface (`internal/output/formatter.go`)
- **Default command**: Running `pray` with no args executes the `today` command

## Dependencies

```go
github.com/spf13/cobra           // CLI framework
github.com/spf13/viper           // Configuration management
github.com/fatih/color           // Colored terminal output
github.com/olekukonko/tablewriter // ASCII tables
gopkg.in/yaml.v3                 // YAML parsing
```

See `go.mod` for the complete list.

## Release Process

Releases are automated using [GoReleaser](https://goreleaser.com/), triggered by git tags:

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

GoReleaser then builds binaries for all platforms, creates the GitHub release, and generates the changelog.

**Supported platforms:** Linux (amd64, arm64, armv7), macOS (amd64, arm64), Windows (amd64)

## Guidelines

- Follow Go best practices and idioms
- Write tests for new features
- Update documentation for user-facing changes
- Format code with `go fmt`
- Lint with `golangci-lint`
- Keep commits atomic and descriptive
