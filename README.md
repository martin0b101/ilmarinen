# Ilmarinen

**Forging reliable execution.**

Ilmarinen is a background job processing system written in Go, backed by PostgreSQL.

It allows applications to submit jobs through an HTTP API and execute them asynchronously using a configurable pool of concurrent workers.

The goal of this project is to build a reliable and extensible job processing engine with support for retries, failure recovery, and graceful shutdown.

## Features

- Asynchronous job execution
- PostgreSQL-backed job queue
- Configurable worker pool
- Concurrent job processing
- Multiple job types and handlers
- Automatic retries with exponential backoff
- Job lease management and failure recovery
- At-least-once execution semantics
- Graceful shutdown
- Structured logging

## Architecture

The system consists of several main components:

- **HTTP API:** Accepts job submissions and provides job status information.
- **Job Service:** Handles job creation and business rules.
- **Queue Service:** Manages job scheduling, claiming, and acknowledgement.
- **Worker Pool:** Coordinates concurrent job execution.
- **Workers:** Execute jobs using registered handlers.
- **PostgreSQL:** Provides persistent job storage.

## Supported Job Types

Initial planned job types:

| Type | Description |
|---|---|
| `email` | Simulates sending an email |
| `report` | Simulates generating a report |

## Getting Started

### Prerequisites

- Go
- PostgreSQL
- Docker (optional)

### Installation

```bash
git clone https://github.com/martin0b101/ilmarinen.git
cd ilmarinen
```

### Run

```bash
go run ./cmd/ilmarinen
```

## Development Status

Ilmarinen is currently under active development.

The project is being developed incrementally, starting with the application foundation and progressing toward a reliable background job processing system.