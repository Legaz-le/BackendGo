# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run the server
go run cmd/server/main.go

# Build
go build ./cmd/server/

# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/handlers/...

# Apply migrations (uses golang-migrate or psql directly)
psql $DATABASE_URL -f migrations/000001_init_schema.up.sql
```

Requires a `.env` file with `DATABASE_URL` set to a PostgreSQL connection string.

## Architecture

This is a Go REST API for a job board. The stack is:
- **Router**: `go-chi/chi/v5`
- **Database**: PostgreSQL via `jackc/pgx/v5` (connection pool)
- **Migrations**: Plain SQL files in `migrations/`

### Request flow

```
HTTP request → chi router (main.go) → middleware → handler (internal/handlers/) → model function (internal/models/) → pgxpool query
```

### Package layout

- `cmd/server/main.go` — entry point, wires up router, DB connection, and middleware
- `internal/database/db.go` — global `pgxpool.Pool` singleton; `GetDB()` used throughout
- `internal/models/` — contains both struct definitions (`models.go`) and DB query functions (`job.go`). Model functions call `database.GetDB()` directly (no repository abstraction for jobs).
- `internal/handlers/` — HTTP handlers; `handlers.go` has the older in-memory Bookmark handlers (currently commented out in routes), `job.go` has the active job handlers
- `internal/middleware/` — logging middleware
- `internal/auth/` — JWT auth (in progress, not yet wired into routes)
- `internal/user/` — user repository (in progress)
- `internal/storage/` — legacy in-memory Bookmark slice (unused)

### Database schema

Three tables: `users`, `jobs`, `applications`. Jobs reference `users(employer_id)`. Applications reference both jobs and users. Migration 002 adds refresh tokens for JWT auth.

### Current state

- Jobs CRUD (`GET /jobs`, `GET /jobs/{id}`, `POST /jobs`, `PUT /jobs/{id}`, `DELETE /jobs/{id}`) is fully wired and hits PostgreSQL.
- Bookmark routes exist in `handlers/handlers.go` but are commented out in `main.go` (legacy in-memory implementation).
- Auth (`internal/auth/`) is implemented but not yet registered on any routes.
- The `Job` model in `models/models.go` does not include `employer_id` or `status` — those fields exist in the DB schema but are not yet surfaced in the Go struct.
