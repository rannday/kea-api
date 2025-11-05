#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "usage: $0 agent|dhcp4|dhcp6 [command]"
  exit 1
fi

service="$1"; shift || true
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

case "$service" in
  agent) path="/control-agent"; manifest="$root/agent/commands.manifest.json" ;;
  dhcp4) path="/dhcp4";         manifest="$root/dhcp4/commands.manifest.json" ;;
  dhcp6) path="/dhcp6";         manifest="$root/dhcp6/commands.manifest.json" ;;
  *) echo "bad service"; exit 1 ;;
esac

branch="codex/$service/$(date +%Y%m%d-%H%M%S)"
git switch -c "$branch"

run_one() {
  local cmd="$1"

  printf '\n=== %s :: %s ===\n' "$service" "$cmd"

  prompt="$(cat <<'EOF'
Use the repository's existing structure and style (two-space indent, no tabs).
Do NOT refactor or move files. No new deps.

Goal: Implement a typed wrapper for a single Kea command in this repo using existing helpers.
Steps (strict):
1) In {SERVICE}/, add a method on Service named by PascalCasing the command (e.g. "status-get" -> StatusGet) that calls client.CallAndDecode with PATH and COMMAND.
2) Define minimal request/response structs in {SERVICE}/types.go that match the REAL JSON already stored in docs/kea/calls/ (if present). Only include fields we assert.
3) Create/refresh {SERVICE}/testdata/{COMMAND}.json using docs/kea/calls if missing; otherwise leave as-is.
4) Add a _test.go that decodes the golden JSON into the repo's envelope type and asserts result==0 and a couple key fields.
5) Keep two-space indent. Don't touch client package API.

Constraints:
- SERVICE={SERVICE}
- PATH={PATH}
- COMMAND={COMMAND}

If something already exists, widen types minimally and keep names stable.
EOF
)"

  # Fill variables in the prompt safely
  prompt="${prompt//\{SERVICE\}/$service}"
  prompt="${prompt//\{PATH\}/$path}"
  prompt="${prompt//\{COMMAND\}/$cmd}"

  codex exec \
    -C "$root" \
    --sandbox workspace-write \
    -a on-failure \
    "$prompt"

  go fmt ./...
  if ! go test ./...; then
    echo "tests failed after $cmd"
    exit 1
  fi

  git add .
  git commit -m "add $service $cmd"
}

if [[ $# -ge 1 ]]; then
  run_one "$1"
else
  jq -r '.[]' "$manifest" | while read -r cmd; do
    run_one "$cmd"
  done
fi

echo "Done on branch: $branch"
