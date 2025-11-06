You are editing an existing Go repo. Do NOT restructure.

Input:
- SERVICE: one of agent|dhcp4|dhcp6
- COMMAND: exact Kea command string (e.g., "status-get")
- PATH: "/control-agent" | "/dhcp4" | "/dhcp6"

Steps:
1) In {SERVICE}/, add or update wrapper on Service using `client.CallAndDecode` with PATH and COMMAND.
   - Function name: Convert hyphenated COMMAND to PascalCase (e.g., status-get -> StatusGet).
   - Signature returns `types.MultiResult[<RespType>]` (or appropriate) and `error`.

2) Define minimal request/response structs in {SERVICE}/types.go consistent with existing style. Only add fields we assert in tests.

3) Add unit test {SERVICE}/{command}_test.go:
   - Read golden JSON from {SERVICE}/testdata/{COMMAND}.json.
   - Unmarshal into `types.MultiResult[<RespType>]` and assert `Result==0` and key fields.

4) If no golden file exists, create a placeholder from docs (shape only), mark with `// TODO: refresh with live`.

5) Run `go test ./...` locally assumptions; keep two-space indent; no tabs.

Do not touch client package API. Do not add third-party deps.
