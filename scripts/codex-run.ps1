Param(
  [Parameter(Mandatory=$true)][ValidateSet("agent","dhcp4","dhcp6")] [string]$Service,
  [string]$Command
)

$ErrorActionPreference = "Stop"
$root = (Resolve-Path "$PSScriptRoot\..").Path

switch ($Service) {
  "agent" { $Path="/control-agent"; $Manifest=Join-Path $root "agent\commands.manifest.json" }
  "dhcp4" { $Path="/dhcp4";         $Manifest=Join-Path $root "dhcp4\commands.manifest.json" }
  "dhcp6" { $Path="/dhcp6";         $Manifest=Join-Path $root "dhcp6\commands.manifest.json" }
}

$branch = "codex/$Service/" + (Get-Date -Format "yyyyMMdd-HHmmss")
git switch -c $branch | Out-Null

function Run-One($cmd) {
  Write-Host "`n=== $Service :: $cmd ===`n" -ForegroundColor Cyan

  $Prompt = @"
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
"@

  $Prompt = $Prompt.Replace("{SERVICE}", $Service).Replace("{PATH}", $Path).Replace("{COMMAND}", $cmd)

  codex exec -C $root --sandbox workspace-write -a on-failure -- $Prompt

  go fmt ./...
  go test ./... | Out-Host

  git add .
  git commit -m "add $Service $cmd" | Out-Null
}

if ($Command) {
  Run-One $Command
} else {
  $Commands = Get-Content $Manifest | ConvertFrom-Json
  foreach ($c in $Commands) { Run-One $c }
}

Write-Host "`nDone on branch: $branch" -ForegroundColor Green
