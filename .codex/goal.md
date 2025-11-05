Goal: Implement typed wrappers for **all** Kea Control Agent, DHCPv4, DHCPv6 API commands using the **existing** structure in this repo:
- Keep `client` helpers central (do not duplicate).
- Public names use PascalCase with `*Get`, `*Write`, etc. Example: `StatusGet`, `ConfigGet`, `ListCommands`.
- Service paths: agent `/control-agent`, dhcp4 `/dhcp4`, dhcp6 `/dhcp6`.
- Strong types live in each service’s `types.go` and shared bits in `types/api_shared.go`.
- Two-space indentation everywhere. No third-party logging libs. Windows + VSCode friendly.
- Add unit tests beside code using golden JSON in `testdata/`. Integration tests behind `//go:build integration`.

Deliverables per command:
1) Wrapper fn on the Service (e.g. `func (s *Service) StatusGet(ctx context.Context) (...)`).
2) Request/response types (only the fields that actually appear).
3) Unit test with golden decode.
4) Update examples if user-facing.

If docs vs live differ, prefer **live JSON** we already captured in tests. Otherwise follow https://kea.readthedocs.io/en/stable/api.html.
