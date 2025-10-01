## Todo API (Go, Gin, Postgres)

A simple JWT-authenticated Todo REST API built with Go, Gin, and Postgres. Includes endpoints to manage users, todo lists, and items. Deployable via Procfile/Railway.

### Features
- Sign up and sign in with JWT auth
- CRUD for todo lists and items
- Layered architecture: handler → service → repository

### Tech Stack
- Go, Gin, sqlx, Postgres, Viper, Logrus, jwt-go

### Project Layout
- `cmd/main.go`: app bootstrap and server start
- `pkg/handler`: HTTP routes, middleware, responses
- `pkg/service`: business logic, validation
- `pkg/repository`: Postgres implementations
- `configs/config.yml`: server and DB config
- `schema/`: SQL migrations

### Configuration
Environment variables override config when present:
- `PORT`: server port (default from `configs/config.yml`)
- `DB_PASSWORD`: Postgres password
- `DATABASE_URL`: full DSN, if provided it overrides individual DB fields

Edit `configs/config.yml` for local development DB params.

### Run Locally
Prerequisites: Go, Postgres running and accessible per config.

```bash
go run ./cmd/main.go
```

Or using the Procfile locally:

```bash
# Requires a Procfile runner (e.g., forego/heroku local)
heroku local
```

### API
- `POST /auth/sign-up` — body: `{ "name", "username", "password" }`
- `POST /auth/sign-in` — body: `{ "username", "password" }` → `{ token }`
- `GET /api/lists` — Authorization: `Bearer <token>`
- Other list/item endpoints are mounted under `/api` with the same auth.

### Tests
Run all tests:

```bash
go test ./...
```

### Deployment
- `Procfile` and `railway.toml` included
- Set `PORT`, `DATABASE_URL` (or individual DB vars) in your platform

### License
MIT


