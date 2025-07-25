# Mata

A lightweight network traffic duplicator that operates at the transport layer as a transparent proxy.

## Features

- **Transparent TCP proxy** - Works with any TCP-based protocol (HTTP, HTTPS, gRPC, etc.)
- **Traffic duplication** - Send identical traffic to multiple targets simultaneously  
- **Transport layer operation** - Protocol-agnostic byte stream duplication
- **Configurable target selection** - Flexible target management and routing
- **Graceful shutdown** - Clean resource management with signal handling

## Quick Start

### Build and Run

```bash
# Build the binary
make build

# Echo mode for testing
./bin/mata -source :8080 -echo

# Single target forwarding
./bin/mata -source :8080 -targets localhost:9000

# Multi-target duplication
./bin/mata -source :8080 -targets app1:80,app2:80,analytics:80
```

### Docker Compose Stack

```bash
# Start complete stack (Nginx + Mata + services)
make docker-up

# Test the stack
curl http://localhost:8000

# View logs
make docker-logs

# Stop stack
make docker-down
```

## Use Cases

- **A/B Testing** - Send same requests to different service versions
- **Migration Testing** - Forward to old system while testing new system
- **Load Testing** - Duplicate production traffic to staging environments  
- **Analytics** - Send copies to monitoring/logging services
- **Development** - Test with real traffic patterns

## Architecture

Mata uses a decorator pattern with interface-based design:
- `ConnectionHandler`: Base connection handling interface
- `DuplicatingHandler`: Decorates handlers with traffic duplication
- `TargetSelector`: Manages target selection and routing

## Development

```bash
# Run tests
make test

# Lint code
make lint

# Format code  
make fmt
```

See [examples/](examples/) for usage examples and Docker configurations.