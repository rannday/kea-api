@echo off
setlocal enabledelayedexpansion

REM Config
set CONTAINER_NAME=kea-int-test
set IMAGE_NAME=kea-lab:latest
set KEA_DOCKER_PATH=..\kea-docker

REM Flags
set RUN_INTEGRATION=false
set FORCE_REBUILD=false

REM Internal flags
set CLEANUP_CONTAINER=false

REM Parse CLI args
for %%x in (%*) do (
  if "%%x"=="-i" set RUN_INTEGRATION=true
  if "%%x"=="-r" set FORCE_REBUILD=true
)

echo.
echo === Running unit tests ===
REM Run unit tests with full coverage and output to a single file
go test -covermode=atomic ^
  -coverpkg=github.com/rannday/kea-api/client,github.com/rannday/kea-api/agent,github.com/rannday/kea-api/dhcp4,github.com/rannday/kea-api/dhcp6 ^
  -coverprofile=coverage.unit.out ^
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
    docker run -d --rm --name %CONTAINER_NAME% ^
      --cap-add=NET_ADMIN ^
      -p 8000:8000 ^
      -p 6767:67/udp ^
      -p 6547:547/udp ^
      -v "%KEA_DOCKER_PATH%\config:/etc/kea:ro" ^
      -v "%KEA_DOCKER_PATH%\logs:/var/log/kea" ^
      %IMAGE_NAME%

    set CLEANUP_CONTAINER=true
    
    echo Waiting for Kea container to become healthy...
    set "MAX_ATTEMPTS=15"
    set /a attempt=0
    :wait_for_healthy
    set /a attempt+=1
    for /f %%h in ('docker inspect -f "{{.State.Health.Status}}" %CONTAINER_NAME% 2^>nul') do set CONTAINER_HEALTH=%%h

    if /i "!CONTAINER_HEALTH!"=="healthy" (
      echo Kea container is healthy.
    ) else if "!attempt!" GEQ "!MAX_ATTEMPTS!" (
      echo ERROR: Kea container did not become healthy after %MAX_ATTEMPTS% attempts.
      
      echo .
      echo Dumping logs:
      docker ps -a --filter "name=%CONTAINER_NAME%"

      echo.
      echo === Healthcheck log ===
      docker inspect -f "{{range .State.Health.Log}}{{println .End \" \" .ExitCode \" \" .Output}}{{end}}" %CONTAINER_NAME%
      
      echo.
      echo === Container logs ===
      docker logs %CONTAINER_NAME%

      exit /b 1
    ) else (
      timeout /t 3 >nul
      goto wait_for_healthy
    )

  ) else (
      echo Container already running. Reusing it.
  )

  echo.
  echo === Running integration tests ===
  go test -covermode=atomic ^
    -coverpkg=github.com/rannday/kea-api/client,github.com/rannday/kea-api/agent,github.com/rannday/kea-api/dhcp4,github.com/rannday/kea-api/dhcp6 ^
    -coverprofile=coverage.integration.out ^
    -tags=integration ^
    github.com/rannday/kea-api/client ^
    github.com/rannday/kea-api/agent ^
    github.com/rannday/kea-api/dhcp4 ^
    github.com/rannday/kea-api/dhcp6

  go tool cover -func=coverage.integration.out > coverage.integration.summary

  if exist coverage.integration.out (
    go tool cover -html=coverage.integration.out -o coverage.integration.html
    echo.
    echo === Integration coverage report generated: coverage.integration.html ===
    start "" coverage.integration.html
  ) else (
    echo Integration coverage file not found. Skipping HTML generation.
  )

)

:cleanup
echo.
echo === Cleanup ===

if "%CLEANUP_CONTAINER%"=="true" (
  echo Stopping container %CONTAINER_NAME%...
  docker stop %CONTAINER_NAME% >nul 2>&1
)

del /q coverage.unit.out >nul 2>&1
if exist coverage.integration.out del /q coverage.integration.out >nul 2>&1

echo.
echo === Done ===
echo.

endlocal