param (
  [Alias("i")]
  [switch]$Integration,

  [Alias("r")]
  [switch]$Rebuild
)

$ContainerName = "kea-int-test"
$ImageName = "kea-custom:latest"
$KeaDockerPath = "..\kea-docker"

Write-Host "`n=== Running unit tests ==="
go test -cover `
  github.com/rannday/kea-api/client `
  github.com/rannday/kea-api/agent `
  github.com/rannday/kea-api/dhcp4 `
  github.com/rannday/kea-api/dhcp6

if ($Integration) {
  Write-Host "`n=== Preparing integration test environment ==="

  if (-Not (Test-Path "$KeaDockerPath\Dockerfile")) {
    Write-Host "kea-docker not found. Cloning..."
    git clone https://github.com/rannday/kea-docker.git $KeaDockerPath
    if ($LASTEXITCODE -ne 0) {
      Write-Error "Failed to clone kea-docker repo. Aborting."
      exit 1
    }
  } else {
    Write-Host "kea-docker repo found. Pulling latest changes..."
    Push-Location $KeaDockerPath
    git pull
    if ($LASTEXITCODE -ne 0) {
      Write-Error "git pull failed. Aborting."
      Pop-Location
      exit 1
    }
    Pop-Location
  }

  Write-Host "`n=== Checking Docker image: $ImageName ==="
  $imageExists = docker images -q $ImageName
  if ($Rebuild) {
    Write-Host "Forcing rebuild of Docker image..."
    docker build --no-cache -t $ImageName $KeaDockerPath
  } elseif (-not $imageExists) {
    Write-Host "Building Docker image..."
    docker build -t $ImageName $KeaDockerPath
  } else {
    Write-Host "Image already built. Skipping rebuild."
  }

  Write-Host "`n=== Starting Kea container ==="
  $containerRunning = docker ps -q --filter "name=$ContainerName"
  if (-not $containerRunning) {
    docker run -d --rm --name $ContainerName -p 8000:8000 $ImageName
    Start-Sleep -Seconds 3
  } else {
    Write-Host "Container already running. Reusing it."
  }

  Write-Host "`n=== Running integration tests with coverage ==="
  go test -tags=integration `
    -coverpkg=./client `
    ./internal/tests/client_test
  
  go test -tags=integration `
    -coverpkg=./agent `
    ./internal/tests/agent_test

  go test -tags=integration `
    -coverpkg=./dhcp4 `
    ./internal/tests/dhcp4_test

  go test -tags=integration `
    -coverpkg=./dhcp6 `
    ./internal/tests/dhcp6_test

  Write-Host "`n=== Done ==="
}
