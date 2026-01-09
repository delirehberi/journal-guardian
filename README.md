[![Release](https://github.com/delirehberi/journal-guardian/actions/workflows/release.yml/badge.svg?event=release)](https://github.com/delirehberi/journal-guardian/actions/workflows/release.yml)

# Log Watcher / Journal Guardian

Log Watcher is a cross-platform tool that monitors system logs (`journalctl` on Linux, `log stream` on macOS) and uses LLMs (Ollama, OpenAI, Gemini, Claude) to automatically suggest fixes for system errors.

## Features

- **Cross-Platform**: Works on Linux (Systemd) and macOS (Unified Logging System).
- **Multi-Provider**: Support for Ollama (local), OpenAI, Google Gemini, and Anthropic Claude.
- **Configurable Sources**: Monitor multiple log files, journalctl, or macOS logs simultaneously.
- **Smart Notifications**: Sends desktop notifications with AI-generated fix suggestions when errors are detected.

## Installation

### Linux (Debian/Ubuntu)

1. Download the latest `.deb` package from the [Releases](https://github.com/delirehberi/journal-guardian/releases) page.
2. Install the package:
   ```bash
   sudo dpkg -i log_watcher_*.deb
   sudo apt-get install -f # Fix dependencies if needed
   ```

### macOS

1. Download the `log_watcher_darwin_*.tar.gz` for your architecture (amd64 for Intel, arm64 for Apple Silicon) from the [Releases](https://github.com/delirehberi/journal-guardian/releases) page.
2. Extract and run:
   ```bash
   tar -xzvf log_watcher_darwin_*.tar.gz
   cd log_watcher
   ./log_watcher
   ```
   *(Optional) You can add it to your PATH or create a LaunchAgent for background execution.*

### From Source (Nix)

If you have [Nix](https://nixos.org/) installed:

```bash
nix build .
./result/bin/log_watcher
```

## Configuration

The application is configured via environment variables. You can set these in your shell or use a `.env` file if running manually.

| Variable | Description | Default |
| :--- | :--- | :--- |
| `LLM_PROVIDER` | Backend provider: `ollama`, `openai`, `gemini`, `claude` | `ollama` |
| `MODEL` | Model name (e.g., `gpt-4`, `llama3`, `gemini-pro`) | `llama3` |
| `OLLAMA_URL` | URL for Ollama API | `http://localhost:11434/api/generate` |
| `OPENAI_API_KEY` | API Key for OpenAI | - |
| `GEMINI_API_KEY` | API Key for Google Gemini | - |
| `ANTHROPIC_API_KEY` | API Key for Anthropic Claude | - |
| `GEMINI_API_KEY` | API Key for Google Gemini | - |
| `ANTHROPIC_API_KEY` | API Key for Anthropic Claude | - |

### Configuration File (Optional)

You can specify log sources in a `config.json` file. The application searches for this file in the following order:
1. Current directory: `./config.json`
2. User config directory:
   - Linux: `~/.config/log_watcher/config.json`
   - macOS: `~/Library/Application Support/log_watcher/config.json`
3. System config directory: `/etc/log_watcher/config.json`

**Example `config.json`:**
```json
{
  "sources": [
    {
      "type": "journalctl",
      "params": {}
    },
    {
      "type": "file",
      "params": {
        "path": "/var/log/nginx/error.log"
      }
    },
    {
      "type": "file",
      "params": {
        "path": "/home/user/app.log"
      }
    }
  ]
}
```

If no configuration file is found, the application falls back to the default OS behavior (monitoring `journalctl` on Linux or Unified Logger on macOS).

## Usage

Start the application:

```bash
./log_watcher
```

### Testing

**Linux:**
```bash
logger -p user.err "This is a test error for log watcher"
```

**macOS:**
```bash
log write "This is a test error for log watcher" --level error
```

The tool should detect the error, query the configured LLM, and send a desktop notification with a suggestion.

## Systemd Service (Linux)

The Debian package installs a user-level systemd service.

```bash
systemctl --user enable --now log_watcher
systemctl --user status log_watcher
```

## Development

### Prerequisites

- Go 1.22+
- Nix (optional, for reproducible builds)
- `libnotify` (Linux only)

### Build Commands

```bash
# Standard Go Build
make build

# Nix Build
make nix-build

# Enter Nix Shell
make nix-shell
```

## Disclaimer

This tool uses AI to generate system administration suggestions. **Always review the suggestions carefully before running any commands.** The developers are not responsible for any damage caused by applying AI-generated fixes.
