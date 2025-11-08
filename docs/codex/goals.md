# Project Goals
Implement typed wrappers for **all** Kea Control Agent, DHCPv4, DHCPv6, and DDNS API commands using the **existing** structure in this repo:
- Keep `client` helpers central (do not duplicate).
- Service paths: agent `/control-agent`, dhcp4 `/dhcp4`, dhcp6 `/dhcp6`.
- Strong types live in each service’s `types.go` and shared bits in `types/api_shared.go`.

## Deliverables per command:
1) Wrapper fn on the Service (e.g. `func (s *Service) StatusGet(ctx context.Context) (...)`).
2) Request/response types (only the fields that actually appear).
3) Unit test with golden decode.
4) Update examples if user-facing.

## Adding Commands (with Codex or manually)
1. Capture/confirm real curl output in `docs/kea/calls/` (see `shared/STATUSGET.md` etc.).
2. Copy JSON to `{service}/testdata/<command>.json`.
3. Add/minimize request/response structs in `{service}/types.go` based on that JSON.
4. Add wrapper in `{service}/{service}.go` calling `client.CallAndDecode` with the service path and command.
5. Add a decode test beside it using the golden JSON.
6. `go test ./...` must pass.
