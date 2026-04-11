# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.0] - 2026-04-11

### Added

- Daily Du'a with real Arabic text, transliteration, translation, and reference — shown in all output formats (table, pretty, JSON, Slack, Discord)
- Iqama time offsets displayed as an extra column in table and pretty output
- Hijri holiday banners in table and pretty output
- Calculation methods 14–23: Moonsighting Committee, Dubai, JAKIM, Tunisia, Algeria, KEMENAG, Morocco, Comunidade Islamica de Lisboa, MUIS/Jordan, Custom
- CI workflow: multi-version Go matrix (1.23/1.24), lint, and cross-platform build jobs
- Conventional Commits linter for pull requests
- Branch protection setup script (`scripts/setup-branch-protection.sh`)
- `config set` support for all new keys: `city`, `country`, `timezone`, `output.color_enabled`, `output.no_emoji`, `features.hijri_holidays`, `features.traveler_mode`, `calendar.*`, `jumuah.*`, `ramadan.*`, `iqama.*`, `cache_enabled`, `update_check`, `api_timeout`
- Cosign keyless signing of release checksums

### Changed

- Default API switched to AlAdhan (api.aladhan.com/v1) — more reliable globally
- IP-detected location source label changed from `"ip"` to `"auto"`; legacy config values migrated on load
- Update checker now uses `golang.org/x/mod/semver` for correct pre-release ordering and detects Homebrew/Scoop/Go install paths to show the right upgrade command
- Setting `address` via `config set` now clears coordinates, and vice-versa

### Fixed

- ICS calendar URL parameters: `ramadan` → `ramadanMode`, `taraweehDuration` → `traweehDuration`

## [1.0.7] - 2026-02-13

### Added

- Install script for Linux/macOS (`install.sh`)
- Install script for Windows (`install.ps1`)
- Quick install instructions in README

## [1.0.6] - 2026-02-12

### Changed

- Updated GoReleaser config and release workflow to enable package uploads

## [1.0.5] - 2026-02-12

### Changed

- Updated module path and imports

## [1.0.4] - 2026-02-12

### Changed

- Updated release workflow to skip publishing step

## [1.0.3] - 2026-02-12

### Added

- GitHub Actions workflow for automated releases

## [1.0.2] - 2026-02-12

### Fixed

- GitHub owner in README installation command

## [1.0.1] - 2026-02-12

### Fixed

- GitHub owner in GoReleaser config

## [1.0.0] - 2026-02-12

### Added

- Initial release
- Prayer times fetching from pray.ahmedelywa.com API
- 23+ calculation method support
- Multiple output formats: table, pretty, JSON, Slack Block Kit, Discord embeds
- IP-based automatic location detection
- Calendar ICS generation and download
- Interactive setup wizard (`pray init`)
- File-based response caching
- Automatic update checker via GitHub Releases
- Cross-platform support: Linux, macOS, Windows (amd64, arm64, armv7)
- Configuration via YAML (`~/.config/pray/config.yaml`)
- Traveler/Qasr mode, Jumu'ah mode, Ramadan mode
- Iqama time offsets
- Hijri date display and holidays
- Qibla direction
- Daily Du'a
- `pray next` countdown to next prayer
- `pray diff` for comparing prayer times across dates
- Shell completion support (bash, zsh, fish, PowerShell)

## [0.0.0] - 2026-02-12

### Added

- Initial project scaffolding

[Unreleased]: https://github.com/AbdElrahmaN31/pray-cli/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/AbdElrahmaN31/pray-cli/compare/v1.0.7...v1.1.0
[1.0.7]: https://github.com/AbdElrahmaN31/pray-cli/compare/v1.0.6...v1.0.7
[1.0.6]: https://github.com/AbdElrahmaN31/pray-cli/compare/v1.0.5...v1.0.6
[1.0.5]: https://github.com/AbdElrahmaN31/pray-cli/compare/v1.0.4...v1.0.5
[1.0.4]: https://github.com/AbdElrahmaN31/pray-cli/compare/v1.0.3...v1.0.4
[1.0.3]: https://github.com/AbdElrahmaN31/pray-cli/compare/v1.0.2...v1.0.3
[1.0.2]: https://github.com/AbdElrahmaN31/pray-cli/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/AbdElrahmaN31/pray-cli/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/AbdElrahmaN31/pray-cli/compare/v0.0.0...v1.0.0
[0.0.0]: https://github.com/AbdElrahmaN31/pray-cli/releases/tag/v0.0.0
