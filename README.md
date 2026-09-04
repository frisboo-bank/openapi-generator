# OpenAPI Generator

A tool for generating OpenAPI specifications from frisboo-core-banking APIs.

## Structure

- `api/` - Go backend (OpenAPI generator service)
- `frontend/` - React/TypeScript frontend application (Rsbuild)

## Build

The root `Makefile` is the single build entry point for the whole repo
(Go backend in `api/`, React/TypeScript frontend in `frontend/`).

```bash
make help        # list all available targets
make install     # install all dependencies (frontend + go modules)
make build       # build backend + frontend
make test        # run all tests
make lint        # lint backend + frontend
```

Per-sub-project commands:

```bash
make run            # run the API server (backend)
make dev-frontend   # start the frontend dev server
```

## Getting Started

### Frontend

```bash
cd frontend
pnpm install
pnpm run dev
```

### Backend (Go)

```bash
cd api
go run ./cmd/main.go
```

## License

Internal use only.
