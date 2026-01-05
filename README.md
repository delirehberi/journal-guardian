# Log Watcher

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
