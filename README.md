# captureCraft API

Go + Fiber API for CaptureCraft (screenshot polishing tool). Provides auth, projects/screenshots CRUD, comments, presigned upload stub, and public share links. Designed to run as a normal server or behind API Gateway/Lambda via the Fiber adapter.

## Stack
- Go 1.22
- Fiber v2
- JWT auth (HS256)
- In-memory store (swap with DB later)
- Optional AWS Lambda + API Gateway via `aws-lambda-go-api-proxy/fiber`

## Features
- Auth: register/login with hashed passwords + JWT issuance
- Projects: create/list/get/update/delete
- Screenshots: create/list/get/update/delete per project
- Comments: list/create per screenshot (owner-only for now)
- Upload helper: returns a mock presigned URL scaffold
- Public share link: tokenized access to a project + screenshots
- Health endpoint: `GET /health`

## Quick start (local server)
1) Install Go 1.22+
2) Fetch deps and run:
```
go mod tidy
go run ./cmd/server
```
3) Configure via env (optional):
```
PORT=8080
JWT_SECRET=change-me
TOKEN_TTL_HOURS=24
RUN_MODE=server
```

## Lambda mode
- Set `RUN_MODE=lambda` and build for Linux: `GOOS=linux GOARCH=amd64 go build -o bootstrap ./cmd/server`
- Package as a Zip or container for API Gateway/Lambda. The adapter `github.com/awslabs/aws-lambda-go-api-proxy/fiber` handles the bridge.

## API sketch (v1)
- `POST /api/v1/auth/register` → `{token, user}`
- `POST /api/v1/auth/login` → `{token, user}`
- `GET /api/v1/projects` (auth)
- `POST /api/v1/projects` (auth)
- `GET /api/v1/projects/:id` (auth)
- `PATCH /api/v1/projects/:id` (auth)
- `DELETE /api/v1/projects/:id` (auth)
- `GET /api/v1/projects/:projectId/screenshots` (auth)
- `POST /api/v1/projects/:projectId/screenshots` (auth)
- `GET /api/v1/screenshots/:id` (auth)
- `PATCH /api/v1/screenshots/:id` (auth)
- `DELETE /api/v1/screenshots/:id` (auth)
- `GET /api/v1/screenshots/:screenshotId/comments` (auth)
- `POST /api/v1/screenshots/:screenshotId/comments` (auth)
- `POST /api/v1/projects/:projectId/share` (auth) → `{token}`
- `GET /public/share/:token` (public) → `{project, screenshots}`
- `POST /api/v1/uploads/presign` (auth) → mock presigned URL payload

## Notes
- Storage is in-memory; swap in a real database by implementing `internal/storage.Store`.
- Upload presign endpoint is a stub—replace with S3, GCS, or R2 signer.
- Authorization is owner-only; adjust rules as collaboration features evolve.