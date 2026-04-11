# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/AbdElrahmaN31/pray-cli/compare/v1.0.7...HEAD
[1.0.7]: https://github.com/AbdElrahmaN31/pray-cli/compare/v1.0.6...v1.0.7
[1.0.6]: https://github.com/AbdElrahmaN31/pray-cli/compare/v1.0.5...v1.0.6
[1.0.5]: https://github.com/AbdElrahmaN31/pray-cli/compare/v1.0.4...v1.0.5
[1.0.4]: https://github.com/AbdElrahmaN31/pray-cli/compare/v1.0.3...v1.0.4
[1.0.3]: https://github.com/AbdElrahmaN31/pray-cli/compare/v1.0.2...v1.0.3
[1.0.2]: https://github.com/AbdElrahmaN31/pray-cli/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/AbdElrahmaN31/pray-cli/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/AbdElrahmaN31/pray-cli/compare/v0.0.0...v1.0.0
[0.0.0]: https://github.com/AbdElrahmaN31/pray-cli/releases/tag/v0.0.0
