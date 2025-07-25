# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Mata is a lightweight network traffic duplicator that operates at the transport layer as a transparent proxy. It intercepts network traffic on a source port and duplicates it to multiple target destinations while maintaining transparency to the original application.

## Development Commands

**Build and Run:**
```bash
make build          # Build mata binary to bin/mata
make run ARGS="..."  # Run with arguments (e.g., ARGS="-source :8080 -echo")
```

**Testing (TDD Approach):**
```bash
make test           # Run all tests with coverage report
make test-pkg PKG=pkg/proxy  # Test specific package
```

**Code Quality:**
```bash
make lint           # Run golangci-lint
make fmt            # Format code with go fmt
```

**Setup:**
```bash
make install-deps   # Install linting tools
make init           # Initialize and download dependencies
```

## Architecture

**Core Design Pattern:** Decorator pattern with interface-based design

**Key Interfaces:**
- `ConnectionHandler` (pkg/proxy): Base interface for handling connections
- `DuplicatingHandler` (pkg/duplicator): Decorator that adds traffic duplication
- `TargetSelector` (pkg/target): Manages target selection and routing

**Directory Structure:**
```
cmd/mata/           # CLI entry point and argument parsing
pkg/proxy/          # Core proxy logic and ConnectionHandler interface
pkg/duplicator/     # Traffic duplication decorators
pkg/target/         # Target selection and management
internal/           # Private application code
test/               # Integration tests
examples/           # Usage examples and demonstrations
```

## Development Conventions

- **TDD**: Write tests first, use table-driven test patterns
- **Interface Design**: Prefer small, focused interfaces
- **Error Handling**: Use explicit error returns, wrap errors with context
- **Concurrency**: Leverage goroutines for connection handling
- **Logging**: Use structured logging (consider `slog` or `logrus`)

## Testing Strategy

- Unit tests for each package with >80% coverage
- Integration tests in `test/` directory
- Benchmark tests for performance-critical paths
- Example tests demonstrating usage patterns