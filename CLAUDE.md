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
- `internal/middleware/` — logging middleware + auth middleware (JWT validation, role-based access)
- `internal/auth/` — JWT auth, refresh tokens, auth handlers, token repository
- `internal/user/` — user repository (CreateUser, GetUserByEmail, GetUserByID)
- `internal/storage/` — legacy in-memory Bookmark slice (unused)

### Database schema

Three tables: `users`, `jobs`, `applications`. Jobs reference `users(employer_id)`. Applications reference both jobs and users. Migration 002 adds refresh tokens for JWT auth.

### Current state

- Jobs CRUD (`GET /jobs`, `GET /jobs/{id}`, `POST /jobs`, `PUT /jobs/{id}`, `DELETE /jobs/{id}`) is fully wired and hits PostgreSQL.
- Bookmark routes exist in `handlers/handlers.go` but are commented out in `main.go` (legacy in-memory implementation).
- Auth is fully implemented and wired: register, login, refresh on `/auth/*`; write job routes protected by `AuthMiddleware` + `RequireRole("employer")`.
- The `Job` model in `models/models.go` does not include `employer_id` or `status` — those fields exist in the DB schema but are not yet surfaced in the Go struct.
- `GET /health` endpoint wired at `handlers/job.go` — returns 200, used by Docker healthcheck.
- Validation uses `github.com/go-playground/validator/v10` struct tags on all auth request types.
- Docker: `Dockerfile` (multi-stage build), `docker-compose.yml` (app + postgres), `.dockerignore` all present at project root.

### Docker commands
```bash
docker compose up --build   # start everything, rebuild app image
docker compose down         # stop everything
docker compose down -v      # stop and wipe DB volume
docker compose logs -f app  # follow app logs

# Apply migrations (run after first up or after down -v)
docker compose exec -T db psql -U postgres -d jobboard < migrations/000001_init_schema.up.sql
docker compose exec -T db psql -U postgres -d jobboard < migrations/000002_create_refresh_tokens.up.sql
```

---

## Developer Roadmap

**Goal:** Get hired as a full-stack developer (React/TypeScript + Go)
**Timeline:** 8-10 weeks
**Flagship Project:** This Job Board Platform

### Roadmap Progress

#### Week 1 — Go Fundamentals ✅ (completed)
#### Week 2 — HTTP & REST APIs ✅ (completed)
#### Week 3 — PostgreSQL + Database ✅ (completed)

#### Week 4 — Authentication & Security ✅ (completed)
- [x] Password hashing with bcrypt (`auth/service.go`)
- [x] JWT generation + validation (`auth/jwt.go`)
- [x] Refresh token generation + hashing (`auth/refresh.go`)
- [x] User repository — CreateUser, GetUserByEmail, GetUserByID (`user/repository.go`)
- [x] Auth HTTP handlers — register, login, refresh (`auth/handler.go`)
- [x] Refresh token DB repository — save/lookup/revoke (`auth/repository.go`)
- [x] Auth middleware — protect routes (`middleware/auth.go`)
- [x] Role-based access — employer vs seeker (`middleware/auth.go`)
- [x] Input validation on all auth endpoints
- [x] Wire auth routes into `main.go`

#### Week 5 — Testing & Code Quality ✅ (completed)
- [x] Unit tests — JWT, password hashing, refresh token hashing (`internal/auth/auth_test.go`)
- [x] Middleware tests — AuthMiddleware and RequireRole (`internal/middleware/middleware_test.go`)
- [x] Refactor validation to use `github.com/go-playground/validator/v10`
- [ ] Handler tests — register, login, refresh (deferred to Week 6, needs test DB)
- [ ] Job handler tests — CRUD with test DB (deferred to Week 6)

#### Week 6 — Docker & Deployment ✅ (completed)
- [x] `GET /health` endpoint
- [x] Multi-stage `Dockerfile`
- [x] `.dockerignore`
- [x] `docker-compose.yml` — app + postgres with healthcheck
- [x] Migrations applied, full stack running

#### Week 7-8 — Frontend React/TypeScript (IN PROGRESS)
#### Week 9-10 — Portfolio Polish & Job Prep (not started)

### Planned Extras (Week 9-10)
- File upload for resumes
- Email notifications
- Rate limiting
- Redis caching
- WebSocket notifications
