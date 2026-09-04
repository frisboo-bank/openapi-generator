# OpenAPI Generator

A tool for generating OpenAPI specifications from frisboo-core-banking APIs.

## Structure

- `api/` - Go backend (OpenAPI generator service)
- `frontend/` - React/TypeScript frontend application (Rsbuild)

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
