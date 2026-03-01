# Job Board

A full-stack job board application built with Go and React.

## Tech Stack

**Backend**
- Go — REST API
- Chi — HTTP router
- PostgreSQL — database
- JWT — authentication (HttpOnly cookies)
- bcrypt — password hashing
- Docker — containerization

**Frontend**
- React + TypeScript
- Vite + Bun
- TanStack Query — server state management
- TanStack Router — client-side routing
- Tailwind CSS — styling

## Features

- Browse job listings
- View job details
- Register as an employer or job seeker
- Login / logout with secure HttpOnly cookies
- Post jobs (employers only, protected route)
- JWT access tokens (15 min) + refresh tokens (7 days)
- Role-based access control

## API Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/jobs` | No | List all jobs |
| GET | `/jobs/{id}` | No | Get job by ID |
| POST | `/jobs` | Employer | Create a job |
| PUT | `/jobs/{id}` | Employer | Update a job |
| DELETE | `/jobs/{id}` | Employer | Delete a job |
| POST | `/auth/register` | No | Register |
| POST | `/auth/login` | No | Login |
| POST | `/auth/logout` | No | Logout |
| POST | `/auth/refresh` | No | Refresh tokens |
| GET | `/auth/me` | No | Get current user |
| GET | `/health` | No | Health check |

## Running Locally

**Prerequisites:** Docker and Docker Compose

1. Clone the repository
2. Start the backend and database:
```bash
docker-compose up --build
```
3. Install frontend dependencies and start the dev server:
```bash
cd frontend
bun install
bun run dev
```
4. Open `http://localhost:5173`

## Project Structure

```
.
├── cmd/server/         # Entry point
├── internal/
│   ├── auth/           # Auth handlers, JWT, bcrypt
│   ├── database/       # DB connection
│   ├── handlers/       # Job handlers
│   ├── middleware/      # Auth middleware, logging
│   ├── models/         # DB models and queries
│   └── user/           # User repository
├── frontend/
│   └── src/
│       ├── routes/     # Page components
│       └── components/ # Shared components
├── Dockerfile
└── docker-compose.yml
```
