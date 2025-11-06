# ISC Kea API — Agent Guidelines
## Overview
Typed Go wrappers for the ISC Kea **Control Agent**, **DHCPv4**, and **DHCPv6** APIs. Real curl outputs in `docs/kea/calls/` are the source of truth for response shapes and tests.

---

## Project Layout

| Path | Purpose |
|------|---------|
| `client/` | Core HTTP client and JSON helpers (`Call`, `CallAndDecode`, etc.). |
| `agent/`, `dhcp4/`, `dhcp6/` | Service wrappers built on `client`. Each exposes a `Service` and typed command funcs. |
| `types/` | Shared envelope/result types (`ResultItem`, `MultiResult`, `ResultCode`). |
| `internal/testenv/` | Mock servers + Docker helpers for integration tests. |
| `examples/` | Runnable examples of common API flows. |
| `docs/kea/calls/` | Real curl requests/responses used as reference and for golden JSON. |

---

## Development & Testing

- Run all tests:

    go test ./...

- Integration tests (requires Docker):

    go test -tags=integration ./agent ./dhcp4 ./dhcp6

- Set environment variables:

    KEA_API_URL=http://localhost:8000
    KEA_API_USERNAME=kea-api
    KEA_API_PASSWORD=kea

- Windows quick run:

    .\test.ps1

- Smoke test:

    go run ./examples/basic

---

## Code Style

- Indentation: **two spaces** (no tabs).
- Public names: **PascalCase** (`StatusGet`, `ConfigGet`, `ListCommands`, …).
- Locals: **camelCase**.
- Max line length ~100 chars.
- No third-party logging libs.
- One exported `Service` per package; path constants per service (`/control-agent`, `/dhcp4`, `/dhcp6`).

---

## Adding Commands (with Codex or manually)

1. Capture/confirm real curl output in `docs/kea/calls/` (see `shared/STATUSGET.md` etc.).
2. Copy JSON to `{service}/testdata/<command>.json`.
3. Add/minimize request/response structs in `{service}/types.go` based on that JSON.
4. Add wrapper in `{service}/{service}.go` calling `client.CallAndDecode` with the service path and command.
5. Add a decode test beside it using the golden JSON.
6. `go test ./...` must pass.

Codex helpers: `.codex/goal.md` and `.codex/recipes/add-command.md`.

---

## Testing Guidelines

- Unit tests live next to code with `_test.go`.
- Integration tests use `//go:build integration` and `internal/testenv`.
- Prefer golden JSON from live Kea over docs; widen/narrow structs to match reality.
- Keep assertions minimal but meaningful (e.g., `result == 0`, key fields present).

---

## Commit & PR Rules

- Commits: short, lowercase (e.g., `add dhcp6 status-get`).
- Squash WIP before PR.
- PR must describe user-facing change, testing performed, follow-ups.
- Link issues; include console output/screenshots for integration runs when relevant.

---

## Security & Configuration

- Never commit credentials.
- Use `KEA_API_URL`, `KEA_API_USERNAME`, `KEA_API_PASSWORD`.
- Clean up integration containers:

    docker stop kea-int-test

- Redact sensitive data in captured JSON.

---

## Style Summary

- Two-space indent
- No tabs
- No external logging deps
- Real JSON is the source of truth
- Typed wrappers (avoid `map[string]any` outside tests)

---

## Codex Automation

We use Codex CLI non-interactively to add/extend commands.

- CachyOS/Linux/macOS:
  ```bash
  scripts/codex-run.sh agent
  scripts/codex-run.sh dhcp4 lease4-get

---

_Last updated: 2025-11-05_
