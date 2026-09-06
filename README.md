# GoTodo

GoTodo is a full-stack todo application with an Ionic/Angular client and a Go backend. It supports account signup with email confirmation, login/logout, todo management, password reset, password change, email change, and account deletion.

## Stack

- Client: Angular 22, Ionic 9, zoneless change detection, Angular signal forms, `httpResource`, RxJS
- Server: Go 1.27, chi, SQLBoiler, scs sessions, PostgreSQL
- Emails: MJML 5 templates compiled into Go template files

## Requirements

- Node.js ^22.22.3, ^24.15.0, or 26+ and npm
- Go 1.27.1 or newer
- Docker for the local PostgreSQL and Inbucket services

## Client

```bash
cd client
npm install
npm start
```

The Angular dev server uses `proxy.conf.json` to forward API requests to the backend on `http://localhost:8080`.

Useful client commands:

```bash
npm run lint
npm run build
npm run serve-dist
```

`npm run build` generates ignored production build metadata, builds the app into `client/dist/app`, and runs `bread-compressor` on the browser output.

## Server

Start PostgreSQL and Inbucket:

```bash
cd server
docker compose up -d
```

Apply the database migrations before the first server start and after pulling
changes that add migrations:

```bash
go run gotodo.rasc.ch/cmd/migrate up
```

The default backend configuration in `server/app.env` uses:

- Database: `gotodo`
- User/password: `gotodo` / `gotodo`
- PostgreSQL address: `127.0.0.1:5432`
- HTTP address: `localhost:8080`
- SMTP address: `localhost:2500`
- App URL for emails: `http://localhost:8100/`

Run the backend:

```bash
cd server
go run gotodo.rasc.ch/cmd/web
```

Run backend tests:

```bash
cd server
go test ./...
```

The Docker-backed browser test builds and exercises the complete client/server
flow:

```bash
task test-integration
```

## Emails

```bash
cd emails
npm install
npm start
```

`npm start` compiles `emails/src/*.mjml` to `emails/output` and then runs `cmd/mailgen` to update `server/mails/*.tmpl`.

## Development Notes

- API routes are served under `/v1`.
- The server reads the rightmost `X-Forwarded-For` address for request logging
  and per-client rate limiting. This assumes exactly one trusted reverse proxy
  is the only network path to the backend.
- Generated TypeScript API types live in `client/src/app/api/types.ts`.
- SQLBoiler models are checked into `server/internal/models`.
- Database migrations live in `server/migrations`.

## Verification

Current health checks:

```bash
cd client && npm run lint
cd client && npm run build
cd emails && npm start
cd server && go test ./...
```
