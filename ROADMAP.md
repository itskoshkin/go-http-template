# Roadmap & Review

Review snapshot of this template and a prioritized backlog. This file is
planning-only — it documents intended work, it does not change any code.

## This template's role

`go-http-template` is the **minimal honest base** of the template family:
Gin + Viper + custom logger + embedded static, with graceful shutdown and
constructor-based DI. Everything else (REST, Fx, TgBot, SVC) descends from
the same conventions.

## Status

- Readiness: **~85%** — does exactly what it advertises, cleanly.
- Build: **`go build ./...` passes**.
- Shared utilities (`utils/text`, `utils/colors`, `logger`) are in sync with
  the rest of the family (they differ only by module import path — no drift).

## Repo-specific backlog

- [ ] **README Go version mismatch** — README says "Go 1.21+", `go.mod` says
      `go 1.25`. Align them.
- [ ] **`/healthz` endpoint** — trivial liveness probe for Docker/k8s.
- [ ] **First test** — one `httptest` example on a handler, so the template
      teaches the testing habit.

## Cross-cutting backlog (whole family)

These gaps exist in every repo and are worth solving once, consistently:

- [ ] **Tests** — at least one example per layer (service with a mocked repo,
      `httptest` controller, `testcontainers-go` repository). Biggest gap.
- [ ] **CI** (`.github/workflows`): `go build`, `go vet`, `golangci-lint`,
      `go test -race`, `govulncheck`.
- [ ] **`.golangci.yml`** linter config.
- [ ] **Makefile / Taskfile**: `build / run / lint / test / docker`.
- [ ] **`LICENSE`**.
- [ ] **Dependabot / Renovate** for dependency updates.

## Planned sibling templates (family-level)

- `go-userbot-template` — Telegram userbot on MTProto (gotd/td) for parsing
  where a bot account can't reach.
- `go-cli-template` — cobra + viper + goreleaser base for future CLIs
  (github-, youtube-, opn-/pf-).
- `go-grpc-template` — gRPC + grpc-gateway + buf.
- `go-worker-template` — HTTP-less consumer (Kafka / NATS / Redis Streams).
