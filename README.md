# Log Watcher

Log Watcher is a tool that monitors `journalctl` logs and uses a local Ollama instance to suggest fixes for errors.

## Installation

### Debian / Ubuntu

1. Download the latest `.deb` package from the [Releases](https://github.com/delirehberi/log_watcher/releases) page.
2. Install the package:
   ```bash
   sudo dpkg -i log_watcher_*.deb
   sudo apt-get install -f # Fix dependencies if needed
   ```

### Nix

```bash
nix build .#deb
```

## Usage

### Systemd Service

The package includes a systemd user service. To enable and start it:

```bash
systemctl --user enable --now log_watcher
```

Check the status:

```bash
systemctl --user status log_watcher
```

### Configuration

The application uses environment variables for configuration. You can create a `~/.config/log_watcher/.env` file or set variables in the systemd service override.

### Variables

- `LLM_PROVIDER`: Select the backend provider. Options: `ollama` (default), `openai`, `gemini`, `claude`.
- `MODEL`: The model name to use (e.g., `gpt-4o`, `gemini-pro`, `claude-3-opus`).
- `OLLAMA_URL`: URL for Ollama (default: `http://localhost:11434/api/generate`).
- `OPENAI_API_KEY`: Required if provider is `openai`.
- `GEMINI_API_KEY`: Required if provider is `gemini`.
- `ANTHROPIC_API_KEY`: Required if provider is `claude`.

A Go application that watches the system journal for errors and asks Ollama for fix suggestions, displaying them via desktop notifications.

## Prerequisites

-   Go 1.22+
-   `libnotify` (for notifications)
-   Ollama (running locally)

## Configuration

The application is configured via environment variables:

| Variable | Description | Default |
| :--- | :--- | :--- |
| `OLLAMA_URL` | URL of the Ollama API | `http://localhost:11434/api/generate` |
| `MODEL` | Name of the Ollama model to use | `llama3` |

## Build and Run

### Using Go

```bash
# Build
make build

# Run
make run
```

### Using Nix

```bash
# Enter development shell
make nix-shell

# Build package
make nix-build
```

## Usage

Start the application:

```bash
./log_watcher
```

Trigger an error (e.g., using `logger`) to test the watcher:

```bash
logger -p user.err "This is a test error for log watcher"
```

## Disclaimer

This tool is co-developed with the assistance of AI. Please review the generated suggestions carefully before applying any fixes to your system.
