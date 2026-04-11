# 🕌 Pray CLI

A beautiful, feature-rich, and lightning-fast command-line tool for Islamic prayer times.

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/AbdElrahmaN31/pray-cli)](https://github.com/AbdElrahmaN31/pray-cli/releases)
[![zread](https://img.shields.io/badge/Ask_Zread-_.svg?style=plastic&color=00b0aa&labelColor=000000&logo=data%3Aimage%2Fsvg%2Bxml%3Bbase64%2CPHN2ZyB3aWR0aD0iMTYiIGhlaWdodD0iMTYiIHZpZXdCb3g9IjAgMCAxNiAxNiIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHBhdGggZD0iTTQuOTYxNTYgMS42MDAxSDIuMjQxNTZDMS44ODgxIDEuNjAwMSAxLjYwMTU2IDEuODg2NjQgMS42MDE1NiAyLjI0MDFWNC45NjAxQzEuNjAxNTYgNS4zMTM1NiAxLjg4ODEgNS42MDAxIDIuMjQxNTYgNS42MDAxSDQuOTYxNTZDNS4zMTUwMiA1LjYwMDEgNS42MDE1NiA1LjMxMzU2IDUuNjAxNTYgNC45NjAxVjIuMjQwMUM1LjYwMTU2IDEuODg2NjQgNS4zMTUwMiAxLjYwMDEgNC45NjE1NiAxLjYwMDFaIiBmaWxsPSIjZmZmIi8%2BCjxwYXRoIGQ9Ik00Ljk2MTU2IDEwLjM5OTlIMi4yNDE1NkMxLjg4ODEgMTAuMzk5OSAxLjYwMTU2IDEwLjY4NjQgMS42MDE1NiAxMS4wMzk5VjEzLjc1OTlDMS42MDE1NiAxNC4xMTM0IDEuODg4MSAxNC4zOTk5IDIuMjQxNTYgMTQuMzk5OUg0Ljk2MTU2QzUuMzE1MDIgMTQuMzk5OSA1LjYwMTU2IDE0LjExMzQgNS42MDE1NiAxMy43NTk5VjExLjAzOTlDNS42MDE1NiAxMC42ODY0IDUuMzE1MDIgMTAuMzk5OSA0Ljk2MTU2IDEwLjM5OTlaIiBmaWxsPSIjZmZmIi8%2BCjxwYXRoIGQ9Ik0xMy43NTg0IDEuNjAwMUgxMS4wMzg0QzEwLjY4NSAxLjYwMDEgMTAuMzk4NCAxLjg4NjY0IDEwLjM5ODQgMi4yNDAxVjQuOTYwMUMxMC4zOTg0IDUuMzEzNTYgMTAuNjg1IDUuNjAwMSAxMS4wMzg0IDUuNjAwMUgxMy43NTg0QzE0LjExMTkgNS42MDAxIDE0LjM5ODQgNS4zMTM1NiAxNC4zOTg0IDQuOTYwMVYyLjI0MDFDMTQuMzk4NCAxLjg4NjY0IDE0LjExMTkgMS42MDAxIDEzLjc1ODQgMS42MDAxWiIgZmlsbD0iI2ZmZiIvPgo8cGF0aCBkPSJNNCAxMkwxMiA0TDQgMTJaIiBmaWxsPSIjZmZmIi8%2BCjxwYXRoIGQ9Ik00IDEyTDEyIDQiIHN0cm9rZT0iI2ZmZiIgc3Ryb2tlLXdpZHRoPSIxLjUiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIvPgo8L3N2Zz4K&logoColor=ffffff)](https://zread.ai/AbdElrahmaN31/pray-cli)

## Features

- **Auto-location detection** from IP address with manual fallback (city name, address, or GPS coordinates)
- **23+ calculation methods** including Egyptian, ISNA, MWL, Umm al-Qura, Diyanet, and more
- **Real-time countdown** to next prayer with live terminal updates
- **Location comparison** -- compare prayer times between two cities side by side
- **Calendar integration** -- ICS generation, subscription URLs, custom alarms (Google, Apple, Outlook)
- **Multiple output formats** -- Table, Pretty, JSON, Slack Block Kit, Discord Embeds
- **Hijri calendar** dates, **Qibla direction**, and daily **Du'a/Adhkar**
- **Ramadan mode** with Iftar, Suhoor, and Taraweeh timings
- **Traveler mode** for shortened prayers, **Jumu'ah** support, **Iqama offsets**
- **Smart caching**, auto-update checks, interactive setup wizard
- **Bilingual** -- English and Arabic

## Installation

### Quick Install (Recommended)

**Linux/macOS:**
```bash
curl -sSL https://raw.githubusercontent.com/AbdElrahmaN31/pray-cli/main/install.sh | sh
```

**Windows (PowerShell):**
```powershell
iwr -useb https://raw.githubusercontent.com/AbdElrahmaN31/pray-cli/main/install.ps1 | iex
```

### Using Go

```bash
go install github.com/AbdElrahmaN31/pray-cli/cmd/pray@latest
```

### Manual Download

Download the latest binary for your platform from the [Releases Page](https://github.com/AbdElrahmaN31/pray-cli/releases/latest).

```bash
# Example: Linux amd64 (replace VERSION, e.g., 1.0.6)
VERSION=1.0.6
curl -L https://github.com/AbdElrahmaN31/pray-cli/releases/download/v${VERSION}/pray-cli_${VERSION}_linux_amd64.tar.gz | tar xz
sudo mv pray /usr/local/bin/
```

Binaries are available for Linux (amd64/arm64), macOS (Intel/Apple Silicon), and Windows (amd64).

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

## Quick Start

```bash
# Interactive setup wizard
pray init

# Or auto-detect your location and save it
pray config detect --save

# Show today's prayer times
pray

# Using a specific city
pray -a "Cairo, Egypt"
```

## Output Examples

### Table (default)

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

### Pretty

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

### JSON

```bash
$ pray -o json
```

```json
{
  "date": { "gregorian": "03 Feb 2026", "hijri": { "day": "15", "month": {"number": 8, "en": "Shaʿbān"}, "year": "1447" } },
  "location": { "latitude": 30.0507, "longitude": 31.2489, "timezone": "Africa/Cairo", "address": "Cairo, Egypt" },
  "timings": { "Fajr": "05:15", "Sunrise": "06:44", "Dhuhr": "12:09", "Asr": "15:11", "Maghrib": "17:34", "Isha": "18:53", "Midnight": "00:09" },
  "nextPrayer": { "name": "Dhuhr", "time": "12:09", "minutesUntil": 45 }
}
```

Also supports `slack` (Block Kit) and `discord` (Embed) formats.

## Command Reference

### Commands

| Command                   | Description                                          |
|---------------------------|------------------------------------------------------|
| `pray`                    | Show today's prayer times (default)                  |
| `pray today`              | Show today's prayer times (explicit)                 |
| `pray next`               | Show next prayer with time remaining                 |
| `pray countdown`          | Live countdown to next prayer (updates every second) |
| `pray get --date <date>`  | Fetch prayer times for a specific date               |
| `pray diff <loc1> <loc2>` | Compare prayer times between two locations           |
| `pray methods`            | List all available calculation methods               |
| `pray init`               | Interactive setup wizard                             |
| `pray version`            | Show version and build information                   |
| `pray completion <shell>` | Generate shell completion (bash/zsh/fish/powershell) |

### Configuration Commands

| Command                         | Description                          |
|---------------------------------|--------------------------------------|
| `pray config detect [--save]`   | Auto-detect location from IP address |
| `pray config show`              | Display current configuration        |
| `pray config set <key> <value>` | Set a configuration value            |
| `pray config get <key>`         | Get a configuration value            |
| `pray config validate`          | Validate current configuration       |
| `pray config reset`             | Reset configuration to defaults      |
| `pray config edit`              | Edit config file in $EDITOR          |
| `pray config export <file>`     | Export configuration to file         |
| `pray config import <file>`     | Import configuration from file       |
| `pray config location`          | Show detailed location information   |

### Calendar Commands

| Command                   | Description                             |
|---------------------------|-----------------------------------------|
| `pray calendar get`       | Download ICS calendar file              |
| `pray calendar url`       | Generate calendar subscription URL      |
| `pray calendar subscribe` | Show detailed subscription instructions |

**Calendar flags:** `-f, --file` (output path), `--months` (1-12), `-d, --duration` (minutes), `--alarm` (offsets e.g. "5,10,15"), `--color` (hex/name), `-e, --events` (indices e.g. "0,2,4" or "all")

### Cache Commands

| Command            | Description                                |
|--------------------|--------------------------------------------|
| `pray cache show`  | Display cache status, size, and file count |
| `pray cache clear` | Clear all cached prayer times data         |
| `pray cache path`  | Show cache directory path                  |

### Global Flags

| Flag                       | Description                                    |
|----------------------------|------------------------------------------------|
| `-a, --address <string>`   | City or address (e.g., "Cairo, Egypt")         |
| `--lat <float>`            | Latitude in decimal degrees                    |
| `--lon <float>`            | Longitude in decimal degrees                   |
| `-A, --auto`               | Auto-detect location from IP address           |
| `-m, --method <int>`       | Calculation method ID (1-23, default: 5)       |
| `-l, --lang <string>`      | Language: `en` or `ar` (default: `en`)         |
| `-o, --output <format>`    | Output: table/pretty/json/slack/discord         |
| `-f, --file <path>`        | Save output to file                            |
| `--qibla`                  | Include Qibla direction                        |
| `--dua`                    | Include daily Du'a/Adhkar                      |
| `--hijri <mode>`           | Hijri date: title/desc/both/none               |
| `--traveler`               | Enable travel/Qasr mode (shortened prayers)    |
| `--jumuah`                 | Add Jumu'ah (Friday) prayer events             |
| `--ramadan`                | Enable Ramadan mode (Iftar, Suhoor, Taraweeh)  |
| `--save`                   | Save current flags as default config           |
| `--no-save`                | Don't save to config (one-time use)            |
| `--no-cache`               | Bypass cache, force fresh API data             |
| `--no-color`               | Disable colored output                         |
| `-v, --verbose`            | Show API calls and debug info                  |
| `-q, --quiet`              | Minimal output (errors only)                   |
| `--config <path>`          | Custom config file path                        |

## Configuration

Configuration is stored in `~/.config/pray/config.yaml` (Linux/macOS) or `%APPDATA%\pray\config.yaml` (Windows).

```yaml
location:
  address: "Cairo, Egypt"
  latitude: 30.0444
  longitude: 31.2357
  timezone: "Africa/Cairo"
  source: "ip"                         # ip/manual/gps

method: 5                              # Calculation method (1-23)
language: "en"                         # en or ar

output:
  format: "table"                      # table/pretty/json/slack/discord
  color_enabled: true
  no_emoji: false

features:
  qibla: true
  dua: true
  hijri: "desc"                        # title/desc/both/none
  traveler_mode: false
  hijri_holidays: false

calendar:
  duration: 25                         # Event duration (minutes)
  months: 3                            # Months to generate (1-12)
  alarm: "5,10,15"                     # Alarm offsets (minutes before prayer)
  events: "all"                        # "all" or indices like "0,2,4"
  color: "#1e90ff"

jumuah:
  enabled: false
  duration: 60

ramadan:
  enabled: false
  iftar_duration: 30
  taraweeh_duration: 60
  suhoor_duration: 30

iqama:
  enabled: false
  offsets: "15,0,10,10,5,10,0"         # Minutes after Adhan per prayer

cache_enabled: true
update_check: true
api_timeout: 30
```

### Popular Calculation Methods

| ID | Method                            | Region                            |
|----|-----------------------------------|-----------------------------------|
| 2  | ISNA                              | North America                     |
| 3  | Muslim World League (MWL)         | Europe, Far East                  |
| 4  | Umm al-Qura University            | Saudi Arabia                      |
| 5  | Egyptian General Authority (default) | Egypt, Syria, Lebanon, Malaysia |
| 12 | Diyanet                           | Turkey                            |

Run `pray methods` for the full list of 23+ methods.

## Integration Examples

### Cron Job

```bash
# Daily prayer times at 5:00 AM
0 5 * * * /usr/local/bin/pray -o pretty > ~/prayer-times.txt
```

### Shell Alias

```bash
# Add to ~/.bashrc or ~/.zshrc
alias prayer='pray -o pretty --qibla --dua'
alias next-prayer='pray next'
```

### tmux Status Bar

```bash
set -g status-right "#(pray next --no-color)"
```

### Slack Webhook

```bash
curl -X POST -H 'Content-type: application/json' \
  --data "$(pray -o slack)" \
  https://hooks.slack.com/services/YOUR/WEBHOOK/URL
```

### Discord Webhook

```bash
curl -X POST -H 'Content-type: application/json' \
  --data "$(pray -o discord)" \
  https://discord.com/api/webhooks/YOUR/WEBHOOK
```

## API & Data Sources

- **Prayer Times**: [pray.ahmedelywa.com](https://pray.ahmedelywa.com) -- AlAdhan-based calculations, ICS generation, Hijri calendar, Qibla direction
- **Location Detection**: [ip-api.com](http://ip-api.com) (free, no key) with manual fallback
- **Caching**: Responses cached at `~/.cache/pray/` (Linux/macOS) or `%LOCALAPPDATA%\pray\` (Windows). Bypass with `--no-cache`.

## Troubleshooting

**"No location configured"**
```bash
pray init                    # Setup wizard
pray config detect --save    # Or auto-detect
```

**"Failed to fetch prayer times"**
```bash
pray -v          # Verbose mode to diagnose
pray --no-cache  # Bypass cache
```

**"Invalid config file"**
```bash
pray config validate   # Check what's wrong
pray config reset      # Reset to defaults
```

**Prayer times seem incorrect**
```bash
pray config show   # Verify location and method
pray -m 2          # Try ISNA
pray methods       # See all methods
```

**Colors not displaying** -- use `pray --no-color` or `pray config set output.color_enabled false`

**File locations:**
- Config: `~/.config/pray/config.yaml` (Linux/macOS) | `%APPDATA%\pray\config.yaml` (Windows)
- Cache: `~/.cache/pray/` (Linux/macOS) | `%LOCALAPPDATA%\pray\` (Windows)

Run `pray config show` and `pray cache path` to see exact paths.

## Privacy

- **Location detection** is optional -- you can always provide manual coordinates
- **No tracking** -- no analytics, telemetry, or usage data is collected
- **Local storage only** -- config and cache are stored on your machine
- **HTTPS** for all API communication

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, project structure, testing, and guidelines.

Quick links:
- [Report a bug](https://github.com/AbdElrahmaN31/pray-cli/issues)
- [Request a feature](https://github.com/AbdElrahmaN31/pray-cli/issues)
- [Start a discussion](https://github.com/AbdElrahmaN31/pray-cli/discussions)

## Acknowledgments

- [pray.ahmedelywa.com](https://pray.ahmedelywa.com) -- Prayer times API
- [AlAdhan](https://aladhan.com/) -- Calculation algorithms
- [Cobra](https://github.com/spf13/cobra) & [Viper](https://github.com/spf13/viper) -- CLI framework and config
- [Fatih Color](https://github.com/fatih/color) -- Terminal colors
- [TableWriter](https://github.com/olekukonko/tablewriter) -- ASCII tables
- [ip-api.com](http://ip-api.com) -- IP geolocation

## License

MIT License - Copyright (c) 2026 Pray CLI Contributors. See [LICENSE](LICENSE).

---

<div align="center">

**Made with ❤️ for the Muslim community**

If this tool benefits you, please ⭐ star the repository and share it with others!

</div>
