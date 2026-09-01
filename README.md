# Network Health Monitoring System

A two-stage network monitoring project that began as a lightweight Python diagnostic prototype and evolved into a concurrent, production-oriented Go service.

The system checks local and remote network health, records structured results, retries temporary failures, and can notify administrators through Telegram when a monitored service becomes unavailable.

> [Read the full project report](./Report.pdf)

## Project Overview

The project was implemented in two phases:

1. **Python prototype** - validates the core diagnostic logic through local system checks, TCP connectivity tests, terminal output, and JSON reporting.
2. **Go implementation** - adds concurrent monitoring, configuration-driven targets, HTTP/HTTPS checks, DNS resolution, retry logic, continuous execution, JSON reports, and Telegram notifications.

The Go implementation is the primary version of the project.

## Features

### Python Prototype

- Detects the active local IP address
- Reads the default gateway on Linux
- Reports available disk space and CPU load averages
- Tests TCP connectivity to configured endpoints
- Displays color-coded terminal results
- Saves diagnostic results in JSON format

### Go Monitoring Service

- Checks every configured domain over HTTP and HTTPS
- Resolves domain names to IPv4 and IPv6 addresses
- Performs additional reachability checks against resolved addresses
- Runs domain checks concurrently using goroutines and `sync.WaitGroup`
- Protects shared results with a mutex
- Retries failed checks with exponential backoff
- Sends Telegram alerts after all retries fail
- Runs continuously at a configurable interval
- Creates timestamped JSON health reports
- Includes a reusable TCP connectivity checker

## Repository Structure

```text
.
├── Golang/
│   ├── configs/
│   │   └── config.json
│   ├── internal/
│   │   ├── checker/
│   │   │   ├── dns.go
│   │   │   ├── http.go
│   │   │   └── tcp.go
│   │   ├── report/
│   │   │   └── report.go
│   │   ├── retry/
│   │   │   └── retry.go
│   │   └── telegram/
│   │       └── telegram.go
│   ├── go.mod
│   ├── go.sum
│   └── main.go
├── Python/
│   └── network_checker.py
├── Report.pdf
├── LICENSE
└── README.md
```

## Requirements

### Go implementation

- Go 1.25 or a compatible version
- Internet access
- A Telegram bot and chat ID if notifications are enabled

### Python prototype

- Python 3
- Linux, because the prototype reads `/proc/loadavg` and uses the `ip route` command
- No third-party Python packages

## Go Configuration

The Go service reads its settings from `Golang/configs/config.json`:

```json
{
  "telegram_bot_token": "YOUR_TELEGRAM_BOT_TOKEN",
  "telegram_chat_id": "YOUR_TELEGRAM_CHAT_ID",
  "interval_seconds": 60,
  "retry_count": 3,
  "domains": [
    "google.com",
    "github.com",
    "example.com"
  ]
}
```

Configuration fields:

| Field | Description |
|---|---|
| `telegram_bot_token` | Token used to call the Telegram Bot API |
| `telegram_chat_id` | Destination chat for failure notifications |
| `interval_seconds` | Delay between monitoring cycles |
| `retry_count` | Maximum attempts for a failed HTTP check |
| `domains` | Domain names monitored during each cycle |

> **Security:** Never commit real bot tokens, passwords, API keys, or other credentials to Git. Use placeholders in public repositories and rotate any credential that has already been exposed.

## Running the Go Service

Run the following commands from the repository root:

```bash
cd Golang
go mod download
go run .
```

The working directory must be `Golang` because the application loads the relative path `configs/config.json`.

The service continues running until it is stopped with `Ctrl+C`. After each monitoring cycle, it creates a report such as:

```text
health_report_20260601_183113.json
```

Each report entry uses the following structure:

```json
{
  "timestamp": "2026-06-01T18:31:13+03:30",
  "target": "https://example.com",
  "status": "OK",
  "details": "200 OK"
}
```

## Monitoring Workflow

For every configured domain, the Go service:

1. Creates HTTP and HTTPS targets.
2. Checks each URL and validates that it returns HTTP status `200`.
3. Retries failures using increasing delays.
4. Records the final result.
5. Sends a Telegram notification if all retries fail.
6. Resolves the domain to its IP addresses and performs additional checks.
7. Waits for all concurrent domain checks to finish.
8. Writes the collected results to a timestamped JSON report.

## Current Scope and Roadmap

This repository is an educational DevOps and network-programming project. Planned improvements include:

- Moving Telegram credentials from the JSON file to environment variables
- Connecting the existing TCP checker to the main monitoring loop
- Improving IPv6 address handling
- Adding unit and integration tests
- Adding structured logging and graceful shutdown
- Exporting Prometheus metrics
- Providing Docker and Kubernetes deployment manifests
- Adding configurable alert-recovery notifications

## Author

**Ehsan Borhani Nia**

- GitHub: [Ehsan-BrhNia](https://github.com/Ehsan-BrhNia)

## License

This project is available under the [MIT License](./LICENSE).
