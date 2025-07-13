@echo off
setlocal enabledelayedexpansion

REM Config
set CONTAINER_NAME=kea-int-test
set IMAGE_NAME=kea-custom:latest
set KEA_DOCKER_PATH=..\kea-docker

REM Flags
set RUN_INTEGRATION=false
set FORCE_REBUILD=false

REM Parse CLI args
for %%x in (%*) do (
  if "%%x"=="-i" set RUN_INTEGRATION=true
  if "%%x"=="-r" set FORCE_REBUILD=true
)

echo.
echo === Running unit tests ===
REM Run unit tests with full coverage and output to a single file
go test -covermode=atomic -coverpkg=./... -coverprofile=coverage.unit.out ^
  github.com/rannday/kea-api/client ^
  github.com/rannday/kea-api/agent ^
  github.com/rannday/kea-api/dhcp4 ^
  github.com/rannday/kea-api/dhcp6

go tool cover -func=coverage.unit.out > coverage.unit.summary

REM Generate HTML report from coverage.out
if exist coverage.unit.out (
  go tool cover -html=coverage.unit.out -o coverage.unit.html
  echo.
  echo === Coverage report generated: coverage.unit.html ===
  start "" coverage.unit.html
) else (
  echo Coverage file not found. Skipping HTML generation.
)

if "%RUN_INTEGRATION%"=="true" (
  echo.
  echo === Preparing integration test environment ===

  :: Ensure kea-docker repo exists
  if not exist "%KEA_DOCKER_PATH%\Dockerfile" (
    echo kea-docker not found. Cloning...
    git clone https://github.com/rannday/kea-docker.git %KEA_DOCKER_PATH%
    if errorlevel 1 (
        echo Failed to clone kea-docker repo. Aborting.
        exit /b 1
    )
  ) else (
    echo kea-docker repo found. Pulling latest changes...
    pushd %KEA_DOCKER_PATH%
    git pull
    if errorlevel 1 (
        echo git pull failed. Aborting.
        popd
        exit /b 1
    )
    popd
  )

  echo.
  echo === Checking Docker image: %IMAGE_NAME% ===
  for /f %%i in ('docker images -q %IMAGE_NAME%') do set IMAGE_EXISTS=%%i

  if "%FORCE_REBUILD%"=="true" (
    echo Forcing rebuild of Docker image...
    docker build --no-cache -t %IMAGE_NAME% %KEA_DOCKER_PATH%
    if errorlevel 1 (
        echo Docker build failed. Aborting.
        exit /b 1
    )
  ) else if not defined IMAGE_EXISTS (
    echo Building Docker image...
    docker build -t %IMAGE_NAME% %KEA_DOCKER_PATH%
    if errorlevel 1 (
        echo Docker build failed. Aborting.
        exit /b 1
    )
  ) else (
    echo Image already built. Skipping rebuild.
  )

  echo.
  echo === Starting Kea container ===
  for /f %%i in ('docker ps -q --filter "name=%CONTAINER_NAME%"') do set CONTAINER_RUNNING=%%i

  if not defined CONTAINER_RUNNING (
    docker run -d --rm --name %CONTAINER_NAME% -p 8000:8000 %IMAGE_NAME%
    set CLEANUP_CONTAINER=true
    echo Waiting for Kea to be ready...
    timeout /t 3 >nul
  ) else (
      echo Container already running. Reusing it.
  )

  echo.
  echo === Running integration tests with coverage ===
  go test -cover -coverpkg=./client -coverprofile=coverage.client.out -tags=integration github.com/rannday/kea-api/client
  go test -cover -coverpkg=./agent -coverprofile=coverage.agent.out -tags=integration github.com/rannday/kea-api/agent
  go test -cover -coverpkg=./dhcp4 -coverprofile=coverage.dhcp4.out -tags=integration github.com/rannday/kea-api/dhcp4
  go test -cover -coverpkg=./dhcp6 -coverprofile=coverage.dhcp6.out -tags=integration github.com/rannday/kea-api/dhcp6

  REM echo.
  REM echo === Coverage summary ===
  REM go tool cover -func=coverage.client.out
  REM go tool cover -func=coverage.agent.out
  REM go tool cover -func=coverage.dhcp4.out
  REM go tool cover -func=coverage.dhcp6.out

  echo.
  echo === Cleanup ===
  echo Deleting coverage files.
  del /q coverage.client.out coverage.agent.out coverage.dhcp4.out coverage.dhcp6.out >nul 2>&1
)

del /q coverage.unit.out >nul 2>&1

echo.
echo === Done ===
echo.

endlocal