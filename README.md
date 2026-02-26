<h1 align="center">GOLANGSTER</h1>

<p align="center">
  A Go linter that checks log messages for style and security issues.<br>
  Compatible with <code>go vet</code> and <code>golangci-lint</code>.
</p>

<p align="center">
  <img src="assets/golangster.png" alt="golangster" width="480">
</p>

<p align="center">
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License: MIT"></a>
  <a href="https://pkg.go.dev/github.com/idakhno/golangster"><img src="https://pkg.go.dev/badge/github.com/idakhno/golangster.svg" alt="Go Reference"></a>
  <img src="https://img.shields.io/badge/go-1.25+-00ADD8?logo=go" alt="Go 1.25+">
  <a href="https://github.com/idakhno/golangster/actions/workflows/ci.yml"><img src="https://github.com/idakhno/golangster/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
</p>

## Rules

| Rule | Description | Example (bad → good) |
|------|-------------|----------------------|
| **lowercase** | Log message must start with a lowercase letter | `"Starting server"` → `"starting server"` |
| **english** | Log message must be in English only | `"Запуск сервера"` → `"starting server"` |
| **special-chars** | No emoji, `!`, `?`, or `...` in log messages | `"started!🚀"` → `"started"` |
| **sensitive** | No sensitive data keywords (password, token, etc.) | `"user password: " + pwd` → remove or mask |

## Supported loggers

- `log` (standard library)
- `log/slog` (standard library, Go 1.21+)
- `go.uber.org/zap` (`*zap.Logger`, `*zap.SugaredLogger`)

## Installation

```bash
go install github.com/idakhno/golangster/cmd/golangster@latest
```

## Usage

### As a standalone `go vet` tool

```bash
# Build
go build -o golangster ./cmd/golangster/

# Run on a package
go vet -vettool=./golangster ./...

# Run with specific rules disabled
go vet -vettool=./golangster -sensitive=false ./...
```

### As a golangci-lint plugin

**Step 1.** Build the plugin (requires CGO):

```bash
CGO_ENABLED=1 go build -buildmode=plugin -o golangster.so plugin/plugin.go
```

**Step 2.** Configure `.golangci.yml`:

```yaml
version: "2"

linters:
  default: none
  enable:
    - golangster

linters-settings:
  custom:
    golangster:
      path: ./golangster.so
      description: Checks log messages for style and security issues
      original-url: github.com/idakhno/golangster
      settings:
        rules:
          lowercase: true
          english_only: true
          no_special_chars: true
          no_sensitive: true
        sensitive_keywords:
          - password
          - token
          - secret
          - api_key
```

> **Note:** The plugin and `golangci-lint` binary must be built with the same Go version and dependency versions.
> Check with: `go version -m $(which golangci-lint)`

**Step 3.** Run:

```bash
golangci-lint run
```

## Flags (standalone mode)

| Flag | Default | Description |
|------|---------|-------------|
| `-lowercase` | `true` | Check that messages start with lowercase |
| `-english` | `true` | Check that messages are in English only |
| `-special-chars` | `true` | Check for emoji and special characters |
| `-sensitive` | `true` | Check for sensitive data keywords |

## Configuration (plugin mode)

Settings are passed via `.golangci.yml` under `linters-settings.custom.golangster.settings`:

```yaml
settings:
  rules:
    lowercase: true
    english_only: true
    no_special_chars: true
    no_sensitive: true
  sensitive_keywords:
    - password
    - token
    - myCustomSecret
```

## Examples

```go
import "log/slog"

// BAD — triggers all 4 rules
slog.Info("Starting server!")           // uppercase + special char
slog.Error("Ошибка подключения")        // non-English
slog.Debug("api_key=" + apiKey)         // sensitive data
slog.Warn("loading...")                 // ellipsis

// GOOD
slog.Info("starting server on port 8080")
slog.Error("connection failed")
slog.Debug("request processed")
slog.Warn("cache miss, fetching from db")
```

## Development

```bash
# Run all tests
go test -race ./...

# Run only unit tests for rules
go test -v ./pkg/analyzer/rules/...

# Run integration tests (analysistest)
go test -v ./pkg/analyzer/...

# Build binary
go build -o golangster ./cmd/golangster/
```

## Project structure

```
golangster/
├── cmd/golangster/      # Standalone binary (go vet -vettool)
├── pkg/analyzer/
│   ├── analyzer.go      # Main analyzer (analysis.Analyzer)
│   ├── config.go        # Configuration
│   ├── detector.go      # Log call detection via AST + type checker
│   └── rules/
│       ├── lowercase.go  # Rule 1: lowercase start + SuggestedFix
│       ├── english.go    # Rule 2: English only
│       ├── special_chars.go  # Rule 3: no emoji/special chars
│       └── sensitive.go  # Rule 4: no sensitive data
├── plugin/plugin.go     # golangci-lint plugin entry point
└── testdata/src/        # analysistest testdata with // want annotations
```
