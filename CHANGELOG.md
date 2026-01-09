# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- **Configurable Log Sources**: Added support for monitoring multiple log sources simultaneously via `config.json`.
- **File Watcher**: New `FileLogWatcher` implementation using `nxadm/tail` to tail log files.
- **Fan-In Architecture**: `main.go` now aggregates logs from multiple concurrent watchers (Files, Journalctl, etc.).
- **Configuration Search**: The application now searches for `config.json` in:
    - Current directory (`./`)
    - User config directory (`~/.config/log_watcher/` or `~/Library/Application Support/log_watcher/`)
    - System config directory (`/etc/log_watcher/`)
- **Documentation**: Updated `README.md` with configuration instructions and example usage.

### Changed
- **Default Behavior**: If no configuration is found, the application correctly falls back to the OS-specific default watcher (Journalctl on Linux, Unified Logger on macOS).
- **Dependencies**: Added `github.com/nxadm/tail` for robust file tailing.
