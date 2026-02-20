# Goverland Core Storage

<a href="https://github.com/goverland-labs/goverland-core-storage?tab=License-1-ov-file" rel="nofollow"><img src="https://img.shields.io/github/license/goverland-labs/goverland-core-storage" alt="GPL 3.0" style="max-width:100%;"></a>
![unit-tests](https://github.com/goverland-labs/goverland-core-storage/workflows/unit-tests/badge.svg)
![golangci-lint](https://github.com/goverland-labs/goverland-core-storage/workflows/golangci-lint/badge.svg)

Core data storage service for the [Goverland](https://goverland.xyz) platform. It consumes DAO governance data from Snapshot via NATS, persists it in PostgreSQL, and exposes a gRPC API for other services.

## Architecture

- **gRPC API** — serves DAOs, proposals, votes, delegates, ENS names, and stats
- **NATS consumers** — ingest events from the Snapshot data source
- **PostgreSQL** — primary data store (via GORM)
- **Background workers** — periodic tasks like top proposals caching, token price updates, delegate calculations

## Project Structure

```
main.go                  # Entry point
internal/
  app.go                 # Application bootstrap
  config/                # Environment-based configuration
  dao/                   # DAO domain (repo, service, server, consumer)
  proposal/              # Proposal domain
  vote/                  # Vote domain
  delegate/              # Delegate domain
  ensresolver/           # ENS name resolution
  stats/                 # Platform statistics
  discord/               # Discord integration
  events/                # Internal event definitions
  pubsub/                # NATS pub/sub helpers
  metrics/               # Prometheus metrics
  logger/                # Zerolog logger setup
protocol/
  storagepb/             # Protobuf definitions and generated Go code
pkg/
  grpcsrv/               # gRPC server utilities
  health/                # Health check endpoint
  middleware/             # gRPC middleware
  prometheus/            # Prometheus HTTP handler
  sdk/zerion/            # Zerion API client
resources/               # SQL migration files
```

## Build & Run

```bash
go build ./...
go test ./...
golangci-lint run
```

## Configuration

The service is configured via environment variables (parsed with [caarlos0/env](https://github.com/caarlos0/env)):

| Variable | Description |
|---|---|
| `LOG_LEVEL` | Zerolog level (default: `info`) |
| `POSTGRES_DSN` | PostgreSQL connection string |
| `NATS_URL` | NATS server URL |
| `API_GRPC_SERVER_BIND` | gRPC listen address (default: `:11000`) |
| `PROMETHEUS_LISTEN` | Prometheus metrics address |
| `HEALTH_LISTEN` | Health check address |

## Contribution Rules

[CONTRIBUTING.md](CONTRIBUTING.md)

## Changelog

[CHANGELOG.md](CHANGELOG.md)
