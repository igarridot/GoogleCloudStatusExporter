# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Prometheus exporter for Google Cloud Platform status incidents written in Go. The application fetches incident data from the GCP status page API and exposes metrics in Prometheus format on port 9118.

## Architecture

The codebase consists of 5 main Go files in the `src/` directory:

- **main.go**: Entry point that starts the metric server
- **server.go**: HTTP server setup, command-line flag parsing, and configuration
- **exporter.go**: Core Prometheus collector implementation with filtering logic
- **request.go**: HTTP client for fetching GCP status data and JSON unmarshaling
- **request_test.go**: Unit tests for the request functionality

The application uses interfaces (`HttpClient`, `DoMethod`) for dependency injection and testability.

## Key Components

- **gcpStatusCollector**: Implements Prometheus collector interface
- **incident/update/product structs**: JSON data models for GCP status API
- **Filtering system**: Supports filtering by zones and products via environment variables
- **Severity mapping**: Converts GCP severity levels to numeric values (0-3)

## Build Commands

```bash
# Build for different platforms
make build-osx     # Build macOS binary to bin/
make build-linux   # Build Linux binary to bin/

# Docker operations
make build         # Build Docker image
make push          # Push image to Docker Hub

# Development
make install-requirements  # Update Go dependencies
make tests                 # Run unit tests
make run-local            # Run Docker container locally
make stop-local           # Stop local container
```

## Running the Application

### Local binary:
```bash
./gcp-status-exporter-linux --web.listen-address ":9119" --exporter.collect-resolved-incidents --exporter.save-last-update
```

### Configuration flags:
- `--web.listen-address`: Server address (default :9118)
- `--web.metrics-path`: Metrics endpoint path (default /metrics)
- `--exporter.collect-resolved-incidents`: Include resolved incidents
- `--exporter.save-last-update`: Add last_update label to metrics
- `--exporter.incidents-zones`: Filter by geographic zones (comma-separated)
- `--exporter.filtered-products`: Filter by product names (comma-separated)

## Development Workflow

1. Make changes in `src/` directory
2. Run `make tests` to verify unit tests pass
3. Use `make build-linux` or `make build-osx` for local testing
4. The Go module is located in `src/go.mod` (not root level)

## Testing

- Unit tests are in `src/request_test.go`
- Run tests with `make tests` (runs `go test` in src directory)
- Test fixtures are available in `fixtures/` directory

## Metrics Format

The exporter provides a single metric `gcp_incidents` with severity-based values:
- resolved: 0
- low: 1  
- medium: 2
- high: 3

Labels include: id, status, product, description, uri, and optionally last_update.
- Every feature should be developed used TDD strategy