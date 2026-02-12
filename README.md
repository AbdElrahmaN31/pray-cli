# 🕌 Pray CLI

A beautiful, feature-rich, and lightning-fast command-line tool for Islamic prayer times.

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/AbdElrahmaN31/pray-cli)](https://github.com/AbdElrahmaN31/pray-cli/releases)

## ✨ Features

### 🌍 Location & Detection
- **Auto-location detection** from IP address with intelligent fallback
- **Manual location** via city name, address, or GPS coordinates
- **Timezone awareness** with automatic DST handling

### ⏰ Prayer Times
- **Multiple calculation methods** (Egyptian, ISNA, MWL, Umm al-Qura, and 20+ more)
- **Real-time countdown** to next prayer with live updates
- **Location comparison** - compare prayer times between two cities
- **Traveler mode** for shortened prayers during travel
- **Jumu'ah (Friday) prayer** support
- **Ramadan mode** with Iftar, Suhoor, and Taraweeh timings
- **Iqama times** with customizable offsets

### 📅 Calendar Integration
- **ICS file generation** with prayer times for multiple months
- **Calendar subscription URLs** for automatic sync
- **Custom alarms** and event durations
- **Configurable events** - include only the prayers you want
- **Calendar colors** for better organization
- **Google Calendar, Apple Calendar, Outlook** compatibility

### 🎨 Display & Output
- **Multiple output formats**: Table, Pretty, JSON, Slack Block Kit, Discord Embeds
- **Beautiful colors and emojis** for enhanced readability
- **Hijri calendar** dates with flexible display options
- **Qibla direction** with compass bearing
- **Daily Du'a and Adhkar** integration
- **No-color mode** for non-terminal environments

### ⚙️ Configuration & Management
- **Persistent configuration** stored in `~/.config/pray/config.yaml`
- **Interactive setup wizard** for first-time users
- **Flexible config management** - set, get, validate, reset, import/export
- **Smart caching** for faster response times
- **Auto-update checks** to stay current
- **Verbose and quiet modes** for debugging or minimal output

## 📦 Installation

### Using Go Install (Recommended)

```bash
go install github.com/anashaat/pray-cli/cmd/pray@latest
```

### Binary Download

Download pre-built binaries for your platform:

**Linux (amd64)**
```bash
curl -L https://github.com/AbdElrahmaN31/pray-cli/releases/latest/download/pray_linux_amd64.tar.gz | tar xz
sudo mv pray /usr/local/bin/
```

**macOS (Apple Silicon)**
```bash
curl -L https://github.com/AbdElrahmaN31/pray-cli/releases/latest/download/pray_darwin_arm64.tar.gz | tar xz
sudo mv pray /usr/local/bin/
```

**macOS (Intel)**
```bash
curl -L https://github.com/AbdElrahmaN31/pray-cli/releases/latest/download/pray_darwin_amd64.tar.gz | tar xz
sudo mv pray /usr/local/bin/
```

**Windows**

Download the Windows binary from the [Releases](https://github.com/AbdElrahmaN31/pray-cli/releases) page.

### From Source

```bash
git clone https://github.com/AbdElrahmaN31/pray-cli.git
cd pray-cli
make install
```

### Verify Installation

```bash
pray version
```

## 🚀 Quick Start

### Interactive Setup

```bash
# Run the setup wizard
pray init
```

### Auto-detect Location

```bash
# Auto-detect your location and save it
pray config detect --save

# Show today's prayer times
pray
```

### Manual Location

```bash
# Using city name
pray -a "Cairo, Egypt"

# Using coordinates
pray --lat 30.0444 --lon 31.2357
```

## 📖 Usage

### Basic Commands

```bash
# Show today's prayer times (default command)
pray

# Show today's prayer times explicitly
pray today

# Show next prayer only
pray next

# Live countdown to next prayer (updates every second)
pray countdown

# Get prayer times for a specific date
pray get --date 2026-03-15

# Compare prayer times between two locations
pray diff "Cairo, Egypt" "London, UK"

# Show all available calculation methods
pray methods

# Show version information
pray version
```

### Output Formats

```bash
# Table format (default) - clean ASCII table
pray

# Pretty format with colors and emojis
pray -o pretty

# JSON output for programmatic use
pray -o json

# Slack Block Kit format
pray -o slack

# Discord embed format
pray -o discord

# Save output to file
pray -o json -f prayer-times.json
```

### Location Options

```bash
# Auto-detect location from your IP address
pray --auto
pray -A

# Use city or address
pray --address "Dubai, UAE"
pray -a "New York"

# Use GPS coordinates
pray --lat 30.0444 --lon 31.2357

# Coordinates with timezone
pray --lat 51.5074 --lon -0.1278
```

### Calculation Methods

```bash
# List all available methods
pray methods

# Use a specific method (1-23)
pray --method 5
pray -m 2  # ISNA method

# Popular methods:
# 1  - University of Islamic Sciences, Karachi
# 2  - Islamic Society of North America (ISNA)
# 3  - Muslim World League (MWL)
# 4  - Umm al-Qura, Makkah
# 5  - Egyptian General Authority (default)
# 12 - Diyanet, Turkey
```

### Calendar Features

```bash
# Download ICS calendar file
pray calendar get

# Download with custom filename
pray calendar get -f ramadan-prayers.ics

# Download with custom settings
pray calendar get --months 6 --duration 30 --alarm "5,10,15"

# Get calendar subscription URL
pray calendar url

# Show detailed subscription instructions
pray calendar subscribe

# Advanced calendar with all features
pray calendar get --ramadan --jumuah --events "0,2,4" -f custom.ics
```

### Configuration Management

```bash
# Interactive setup wizard (recommended for first-time users)
pray init

# Auto-detect and save your location
pray config detect --save

# Set individual configuration values
pray config set method 5
pray config set language ar
pray config set features.qibla true
pray config set calendar.duration 30

# Get a specific config value
pray config get method
pray config get location.address

# Show all configuration
pray config show

# Validate your configuration
pray config validate

# Reset to default configuration
pray config reset

# Edit config in your default editor
pray config edit

# Export configuration to file
pray config export backup-config.yaml

# Import configuration from file
pray config import backup-config.yaml

# Show detailed location information
pray config location
```

### Cache Management

```bash
# Show cache status and size
pray cache show

# Clear all cached data
pray cache clear

# Show cache directory path
pray cache path

# Bypass cache for fresh data
pray --no-cache
```

### Feature Flags

```bash
# Include Qibla direction
pray --qibla

# Include daily Du'a
pray --dua

# Hijri date display options
pray --hijri title     # Show in title
pray --hijri desc      # Show in description
pray --hijri both      # Show in both places
pray --hijri none      # Don't show

# Enable traveler mode (shortened prayers)
pray --traveler

# Enable Jumu'ah (Friday prayer)
pray --jumuah

# Enable Ramadan mode (Iftar, Suhoor, Taraweeh)
pray --ramadan

# Combine multiple features
pray --qibla --dua --ramadan --jumuah
```

### Advanced Usage

```bash
# Use different language (en or ar)
pray --lang ar
pray -l en

# Save current flags as default config
pray --address "Tokyo" --method 2 --qibla --save

# One-time use without saving to config
pray -a "Paris" --no-save

# Verbose output (show API calls and debug info)
pray --verbose
pray -v

# Quiet mode (errors only)
pray --quiet
pray -q

# Disable colored output
pray --no-color

# Full-featured command example
pray -a "Mecca" -m 4 --qibla --dua --ramadan --jumuah \
  --hijri both --lang ar -o pretty
```

## ⚙️ Configuration

Configuration is stored in `~/.config/pray/config.yaml` (Linux/macOS) or `%APPDATA%\pray\config.yaml` (Windows).

### Configuration File Structure

```yaml
# Location settings
location:
  address: "Cairo, Egypt"              # Human-readable address
  latitude: 30.0444                    # Latitude in decimal degrees
  longitude: 31.2357                   # Longitude in decimal degrees
  timezone: "Africa/Cairo"             # IANA timezone identifier
  detected_at: "2026-02-03T10:30:00Z" # Auto-detection timestamp
  source: "ip"                         # Source: ip/manual/gps

# Calculation method (1-23)
method: 5                              # Default: Egyptian General Authority

# Language (en or ar)
language: "en"

# Output preferences
output:
  format: "table"                      # Default: table, pretty, json, slack, discord
  color_enabled: true                  # Enable colored output
  no_emoji: false                      # Disable emojis

# Feature toggles
features:
  qibla: true                          # Include Qibla direction
  dua: true                            # Include daily Du'a/Adhkar
  hijri: "desc"                        # Hijri date: title/desc/both/none
  hijri_holidays: false                # Include Islamic holidays
  traveler_mode: false                 # Enable travel/Qasr mode

# Calendar settings
calendar:
  duration: 25                         # Event duration in minutes
  months: 3                            # Number of months to generate (1-12)
  alarm: "5,10,15"                     # Alarm offsets before prayer (minutes)
  events: "all"                        # Events: "all" or indices "0,2,4"
  color: "#1e90ff"                     # Calendar color (hex or name)

# Jumu'ah (Friday prayer) settings
jumuah:
  enabled: false                       # Enable Jumu'ah events
  duration: 60                         # Duration in minutes

# Ramadan settings
ramadan:
  enabled: false                       # Enable Ramadan mode
  iftar_duration: 30                   # Iftar event duration
  taraweeh_duration: 60                # Taraweeh event duration
  suhoor_duration: 30                  # Suhoor event duration

# Iqama (congregation time) settings
iqama:
  enabled: false                       # Enable Iqama times
  offsets: "15,0,10,10,5,10,0"        # Minutes after Adhan for each prayer

# Advanced settings
cache_enabled: true                    # Enable response caching
update_check: true                     # Check for CLI updates
api_timeout: 30                        # API timeout in seconds
```

### Calculation Methods

| ID | Method | Region/Organization |
|----|--------|---------------------|
| 1 | University of Islamic Sciences, Karachi | Pakistan, Bangladesh, India, Afghanistan |
| 2 | Islamic Society of North America (ISNA) | North America (USA, Canada) |
| 3 | Muslim World League (MWL) | Europe, Far East, parts of Americas |
| 4 | Umm al-Qura University | Saudi Arabia |
| 5 | Egyptian General Authority of Survey | Egypt, Syria, Lebanon, Malaysia (default) |
| 6 | Institute of Geophysics, University of Tehran | Iran, Afghanistan, parts of Azerbaijan |
| 7 | Gulf Region | Kuwait, Qatar, Bahrain, United Arab Emirates |
| 8 | Kuwait | Kuwait |
| 9 | Qatar | Qatar |
| 10 | Majlis Ugama Islam Singapura | Singapore |
| 11 | Union Organization Islamic de France | France |
| 12 | Diyanet İşleri Başkanlığı | Turkey |
| 13 | Spiritual Administration of Muslims of Russia | Russia |
| 14 | Moonsighting Committee Worldwide | N/A |
| 15 | Dubai | United Arab Emirates |
| 16 | Jabatan Kemajuan Islam Malaysia | Malaysia |
| 17 | Tunisia | Tunisia |
| 18 | Algeria | Algeria |
| 19 | KEMENAG | Indonesia |
| 20 | Morocco | Morocco |
| 21 | Comunidade Islamica de Lisboa | Portugal |
| 22 | Ministry of Awqaf, Islamic Affairs and Holy Places | Jordan |
| 23 | Presidency of Religious Affairs | Turkey (alternative) |

Run `pray methods` to see complete details and descriptions.

### Prayer Event Indices (for --events flag)

| Index | Prayer | Description |
|-------|--------|-------------|
| 0 | Fajr | Dawn prayer |
| 1 | Sunrise | Sunrise time (not a prayer) |
| 2 | Dhuhr | Noon prayer |
| 3 | Asr | Afternoon prayer |
| 4 | Maghrib | Sunset prayer |
| 5 | Isha | Night prayer |
| 6 | Midnight | Islamic midnight |

Example: `--events "0,2,4"` will include only Fajr, Dhuhr, and Maghrib.

## 📋 Examples

### Show Prayer Times

```bash
$ pray

┌──────────────────────────────────────────────────┐
│           Prayer Times - Cairo, Egypt            │
│                   03 Feb 2026                    │
│                15 Shaʿbān 1447                   │
├──────────────────────────────────────────────────┤
┌──────────┬───────┬──────────────────────────────┐
│  PRAYER  │ TIME  │           STATUS             │
├──────────┼───────┼──────────────────────────────┤
│ Fajr     │ 05:15 │ ✓ Passed                     │
│ Sunrise  │ 06:44 │ ✓ Passed                     │
│ Dhuhr    │ 12:09 │ ▶ Next (in 45 min)           │
│ Asr      │ 15:11 │                              │
│ Maghrib  │ 17:34 │                              │
│ Isha     │ 18:53 │                              │
│ Midnight │ 00:09 │                              │
└──────────┴───────┴──────────────────────────────┘
│ Method: Egyptian General Authority of Survey     │
└──────────────────────────────────────────────────┘
```

### Pretty Output

```bash
$ pray -o pretty

🕌 Prayer Times for Cairo, Egypt
📅 03 Feb 2026 | 15 Shaʿbān 1447

🌅 Fajr      05:15  ✓ Passed
🌄 Sunrise   06:44  ✓ Passed
☀️ Dhuhr     12:09  ▶ Next prayer in 45 minutes
🌤️ Asr       15:11
🌆 Maghrib   17:34
🌙 Isha      18:53
🌃 Midnight  00:09

🧭 Qibla Direction: SE (136.2°)
⚙️  Method: Egyptian General Authority of Survey
```

### JSON Output

```bash
$ pray -o json

{
  "date": {
    "gregorian": "03 Feb 2026",
    "hijri": {
      "day": "15",
      "month": {"number": 8, "en": "Shaʿbān"},
      "year": "1447"
    }
  },
  "location": {
    "latitude": 30.0507,
    "longitude": 31.2489,
    "timezone": "Africa/Cairo",
    "address": "Cairo, Egypt"
  },
  "timings": {
    "Fajr": "05:15",
    "Sunrise": "06:44",
    "Dhuhr": "12:09",
    "Asr": "15:11",
    "Maghrib": "17:34",
    "Isha": "18:53",
    "Midnight": "00:09"
  },
  "nextPrayer": {
    "name": "Dhuhr",
    "time": "12:09",
    "minutesUntil": 45
  }
}
```

## 🔧 Command Reference

### Main Commands

| Command | Description |
|---------|-------------|
| `pray` | Show today's prayer times (default command) |
| `pray today` | Show today's prayer times (explicit alias) |
| `pray next` | Show next prayer only with time remaining |
| `pray countdown` | Live countdown to next prayer (updates every second) |
| `pray get` | Fetch prayer times with custom date |
| `pray diff <loc1> <loc2>` | Compare prayer times between two locations |
| `pray methods` | List all available calculation methods |
| `pray init` | Interactive setup wizard |
| `pray version` | Show version, commit, and build information |
| `pray completion` | Generate shell completion scripts |

### Configuration Commands

| Command | Description |
|---------|-------------|
| `pray config detect [--save]` | Auto-detect location from IP address |
| `pray config show` | Display current configuration |
| `pray config set <key> <value>` | Set a configuration value |
| `pray config get <key>` | Get a configuration value |
| `pray config validate` | Validate current configuration |
| `pray config reset` | Reset configuration to defaults |
| `pray config edit` | Edit config file in $EDITOR |
| `pray config export <file>` | Export configuration to file |
| `pray config import <file>` | Import configuration from file |
| `pray config location` | Show detailed location information |

### Calendar Commands

| Command | Description |
|---------|-------------|
| `pray calendar get` | Download ICS calendar file |
| `pray calendar url` | Generate calendar subscription URL |
| `pray calendar subscribe` | Show detailed subscription instructions |

**Calendar Flags:**
- `-f, --file <path>` - Output file path
- `--months <n>` - Number of months (1-12, default: 3)
- `-d, --duration <n>` - Event duration in minutes (default: 25)
- `--alarm <offsets>` - Alarm offsets, e.g., "5,10,15"
- `--color <color>` - Calendar color (hex or name)
- `-e, --events <indices>` - Events to include, e.g., "0,2,4" or "all"

### Cache Commands

| Command | Description |
|---------|-------------|
| `pray cache show` | Display cache status, size, and file count |
| `pray cache clear` | Clear all cached prayer times data |
| `pray cache path` | Show cache directory path |

### Global Flags

#### Location Flags
| Flag | Description |
|------|-------------|
| `-a, --address <string>` | City or address (e.g., "Cairo, Egypt") |
| `--lat <float>` | Latitude in decimal degrees |
| `--lon <float>` | Longitude in decimal degrees |
| `-A, --auto` | Auto-detect location from IP address |

#### Calculation & Display Flags
| Flag | Description |
|------|-------------|
| `-m, --method <int>` | Calculation method ID (1-23, default: 5) |
| `-l, --lang <string>` | Language: en or ar (default: en) |
| `--qibla` | Include Qibla direction |
| `--dua` | Include daily Du'a/Adhkar |
| `--hijri <mode>` | Hijri date: title/desc/both/none |

#### Feature Flags
| Flag | Description |
|------|-------------|
| `--traveler` | Enable travel/Qasr mode (shortened prayers) |
| `--jumuah` | Add Jumu'ah (Friday) prayer events |
| `--ramadan` | Enable Ramadan mode (Iftar, Suhoor, Taraweeh) |

#### Output Flags
| Flag | Description |
|------|-------------|
| `-o, --output <format>` | Output format: table/pretty/json/slack/discord |
| `-f, --file <path>` | Save output to file |
| `--no-color` | Disable colored output |

#### Other Flags
| Flag | Description |
|------|-------------|
| `--save` | Save current flags as default config |
| `--no-save` | Don't save to config (one-time use) |
| `--no-cache` | Bypass cache, force fresh API data |
| `-v, --verbose` | Verbose output (show API calls, debug info) |
| `-q, --quiet` | Minimal output (errors only) |
| `--config <path>` | Custom config file path |
| `-h, --help` | Show help for any command |

## 🏗️ Project Structure

```
pray-cli/
├── cmd/
│   └── pray/              # Main CLI application
│       ├── main.go        # Entry point
│       └── cmd/           # Cobra commands
│           ├── root.go    # Root command and global flags
│           ├── today.go   # Default command (show today's times)
│           ├── next.go    # Next prayer command
│           ├── countdown.go  # Live countdown command
│           ├── diff.go    # Location comparison command
│           ├── get.go     # Fetch prayer times with date
│           ├── calendar.go   # Calendar operations
│           ├── config.go  # Configuration management
│           ├── cache.go   # Cache management
│           ├── methods.go # List calculation methods
│           ├── init.go    # Interactive setup wizard
│           ├── version.go # Version information
│           └── completion.go # Shell completions
│
├── internal/              # Private application code
│   ├── api/              # API client and types
│   │   ├── client.go     # HTTP client with retry logic
│   │   ├── cached_client.go # Cached API client
│   │   ├── params.go     # Request parameter builder
│   │   ├── types.go      # Response structures
│   │   └── validator.go  # Parameter validation
│   │
│   ├── config/           # Configuration management
│   │   ├── config.go     # Config structures
│   │   ├── loader.go     # Load/save configuration
│   │   ├── defaults.go   # Default values
│   │   └── validator.go  # Config validation
│   │
│   ├── location/         # Location detection
│   │   ├── detector.go   # IP-based geolocation
│   │   └── types.go      # Location structures
│   │
│   ├── output/           # Output formatters
│   │   ├── formatter.go  # Formatter interface
│   │   ├── table.go      # ASCII table output
│   │   ├── pretty.go     # Colored pretty output
│   │   ├── json.go       # JSON output
│   │   ├── slack.go      # Slack Block Kit format
│   │   └── discord.go    # Discord Embed format
│   │
│   ├── calendar/         # Calendar operations
│   │   ├── generator.go  # ICS URL generator
│   │   ├── downloader.go # ICS file downloader
│   │   └── subscriber.go # Subscription instructions
│   │
│   ├── cache/            # Caching system
│   │   └── cache.go      # Cache implementation
│   │
│   ├── ui/               # User interface components
│   │   ├── spinner.go    # Loading spinner
│   │   └── wizard.go     # Interactive setup wizard
│   │
│   └── update/           # Update checker
│       └── checker.go    # Check for new releases
│
├── pkg/                  # Public, reusable packages
│   └── prayer/
│       ├── times.go      # Prayer time utilities
│       └── methods.go    # Calculation methods data
│
├── bin/                  # Compiled binaries (git-ignored)
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── Makefile             # Build automation
├── .goreleaser.yml      # Release automation config
├── LICENSE              # MIT License
└── README.md            # This file
```

## 🌐 API & Data Sources

### Prayer Times API

This CLI uses the **[pray.ahmedelywa.com](https://pray.ahmedelywa.com)** API, which provides:
- Accurate prayer times using AlAdhan calculations
- ICS calendar file generation
- Support for 23+ calculation methods
- Hijri calendar integration
- Qibla direction calculation

**API Endpoints:**
- `GET /api/prayer-times.json` - JSON prayer times data
- `GET /api/prayer-times.ics` - ICS calendar file

### Location Detection

The CLI uses multiple IP geolocation services with intelligent fallback:
1. **Primary**: [ip-api.com](http://ip-api.com) (free, no key required)
2. **Fallback**: Manual input or GPS coordinates

### Caching

- Responses are cached locally to improve performance
- Cache location: `~/.cache/pray/` (Linux/macOS) or `%LOCALAPPDATA%\pray\` (Windows)
- Cache TTL: Based on API response headers (typically 1-24 hours)
- Can be bypassed with `--no-cache` flag

## 🛠️ Development

### Prerequisites

- Go 1.23 or higher
- Make (optional, but recommended)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/AbdElrahmaN31/pray-cli.git
cd pray-cli

# Download dependencies
make deps

# Build the binary
make build

# Install locally
make install

# Run tests
make test

# Run tests with coverage
make test-coverage
```

### Development Commands

```bash
# Format code
make fmt

# Lint code
make lint

# Tidy dependencies
make tidy

# Run without installing
make run

# Run with arguments
make run-args ARGS="--help"

# Build for all platforms
make build-all

# Clean build artifacts
make clean
```

### Cross-Compilation

```bash
# Linux
make build-linux-amd64
make build-linux-arm64

# macOS
make build-darwin-amd64  # Intel
make build-darwin-arm64  # Apple Silicon

# Windows
make build-windows-amd64
```

### Dependencies

```go
// Core dependencies
github.com/spf13/cobra    // CLI framework
github.com/spf13/viper    // Configuration management
github.com/fatih/color    // Colored terminal output
github.com/olekukonko/tablewriter // ASCII tables
gopkg.in/yaml.v3          // YAML parsing
```

See `go.mod` for complete dependency list.

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
make test-coverage

# View coverage in browser
open coverage.html
```

### Test Coverage

The project includes comprehensive tests for:
- ✅ API client and response parsing
- ✅ Location detection and validation
- ✅ Configuration management
- ✅ Output formatters
- ✅ Calendar generation
- ✅ Parameter validation

## 🚀 Release Process

Releases are automated using [GoReleaser](https://goreleaser.com/):

```bash
# Create a new tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# GoReleaser automatically:
# - Builds binaries for all platforms
# - Creates GitHub release
# - Generates changelog
# - Updates Homebrew tap
# - Creates checksums
```

**Supported Platforms:**
- Linux: amd64, arm64, arm (v7)
- macOS: amd64 (Intel), arm64 (Apple Silicon)
- Windows: amd64

## 🐛 Troubleshooting

### Common Issues

**Problem: "No location configured"**
```bash
# Solution: Run the setup wizard or detect location
pray init
# or
pray config detect --save
```

**Problem: "Failed to fetch prayer times"**
```bash
# Solution: Check your internet connection and try with verbose mode
pray -v

# Or bypass cache
pray --no-cache
```

**Problem: "Invalid config file"**
```bash
# Solution: Validate or reset your configuration
pray config validate
# or
pray config reset
```

**Problem: Colors not displaying correctly**
```bash
# Solution: Disable colors
pray --no-color

# Or set in config
pray config set output.color_enabled false
```

**Problem: Prayer times seem incorrect**
```bash
# Solution: Verify your location and method
pray config show

# Try a different calculation method
pray -m 2  # ISNA method
pray -m 3  # MWL method

# Check all methods
pray methods
```

### Debug Mode

Enable verbose output to see detailed information:

```bash
pray -v        # Verbose mode
pray --verbose # Show API calls, cache hits, timing info
```

### Configuration Location

**Linux/macOS:**
- Config: `~/.config/pray/config.yaml`
- Cache: `~/.cache/pray/`

**Windows:**
- Config: `%APPDATA%\pray\config.yaml`
- Cache: `%LOCALAPPDATA%\pray\`

To find exact paths:
```bash
pray config show  # Shows config file location
pray cache path   # Shows cache directory
```

## 🔐 Privacy & Security

- **Location Detection**: IP-based location detection is optional. You can always use manual addresses or coordinates.
- **Data Storage**: Only configuration and cached prayer times are stored locally. No personal data is sent to external servers.
- **API Communication**: All API requests are made over HTTPS.
- **No Tracking**: This CLI does not collect or transmit any analytics or usage data.

## 📊 Performance

- **Fast**: Prayer times fetched in <1 second (with good connection)
- **Cached**: Repeated requests use local cache for instant results
- **Parallel**: Location comparisons fetch data concurrently
- **Lightweight**: Binary size ~8-12 MB (static compilation)

## 🌍 Supported Languages

- **English** (en) - Default
- **Arabic** (ar) - العربية

```bash
# Use Arabic
pray --lang ar
pray -l ar

# Save as default
pray config set language ar
```

## 🎨 Shell Completions

Generate shell completion scripts for faster command entry:

```bash
# Bash
pray completion bash > /etc/bash_completion.d/pray

# Zsh
pray completion zsh > "${fpath[1]}/_pray"

# Fish
pray completion fish > ~/.config/fish/completions/pray.fish

# PowerShell
pray completion powershell > pray.ps1
```

## 📱 Integration Examples

### Cron Job (Daily Prayer Times)

```bash
# Add to crontab (run at 5:00 AM daily)
0 5 * * * /usr/local/bin/pray -o pretty > ~/prayer-times.txt
```

### Shell Alias

```bash
# Add to ~/.bashrc or ~/.zshrc
alias prayer='pray -o pretty --qibla --dua'
alias next-prayer='pray next'
```

### tmux/Terminal Display

```bash
# Add to tmux status bar
set -g status-right "#(pray next --no-color)"
```

### Slack Webhook

```bash
# Post prayer times to Slack
curl -X POST -H 'Content-type: application/json' \
  --data "$(pray -o slack)" \
  https://hooks.slack.com/services/YOUR/WEBHOOK/URL
```

### Discord Webhook

```bash
# Post prayer times to Discord
curl -X POST -H 'Content-type: application/json' \
  --data "$(pray -o discord)" \
  https://discord.com/api/webhooks/YOUR/WEBHOOK
```

## 📚 Use Cases

- **Personal Use**: Daily prayer time reminders and notifications
- **Mosques**: Generate monthly prayer calendars for community
- **Mobile Apps**: Backend API integration for prayer time apps
- **Websites**: Display prayer times on Islamic websites
- **Automation**: Integrate with smart home systems for prayer notifications
- **Travel**: Quickly find prayer times when traveling
- **Comparison**: Compare prayer times across different cities/countries

## 📄 License

MIT License - Copyright (c) 2026 Pray CLI Contributors

See [LICENSE](LICENSE) file for full details.

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. **Report Bugs**: Open an issue with details and reproduction steps
2. **Suggest Features**: Share your ideas in the issues section
3. **Submit PRs**: Fork, create a feature branch, and submit a pull request
4. **Improve Docs**: Help improve documentation and examples
5. **Test**: Test on different platforms and report issues

### Development Guidelines

- Follow Go best practices and idioms
- Write tests for new features
- Update documentation for changes
- Format code with `go fmt`
- Lint with `golangci-lint`
- Keep commits atomic and descriptive

## 🙏 Acknowledgments

- **[pray.ahmedelywa.com](https://pray.ahmedelywa.com)** - Prayer times API and calendar generation
- **[AlAdhan API](https://aladhan.com/)** - Calculation algorithms and data
- **[Cobra](https://github.com/spf13/cobra)** - Powerful CLI framework
- **[Viper](https://github.com/spf13/viper)** - Configuration management
- **[Fatih Color](https://github.com/fatih/color)** - Terminal colors
- **[TableWriter](https://github.com/olekukonko/tablewriter)** - ASCII table formatting
- **[ip-api.com](http://ip-api.com)** - IP geolocation service

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/AbdElrahmaN31/pray-cli/issues)
- **Discussions**: [GitHub Discussions](https://github.com/AbdElrahmaN31/pray-cli/discussions)
- **Email**: Open an issue for bug reports or feature requests

## 🔗 Links

- **Repository**: [github.com/anashaat/pray-cli](https://github.com/AbdElrahmaN31/pray-cli)
- **Releases**: [Releases Page](https://github.com/AbdElrahmaN31/pray-cli/releases)
- **Documentation**: This README and `pray --help`
- **API Documentation**: [pray.ahmedelywa.com](https://pray.ahmedelywa.com)

---

<div align="center">

**Made with ❤️ for the Muslim community**

If this tool benefits you, please ⭐ star the repository and share it with others!

</div>
