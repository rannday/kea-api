@echo off
setlocal enabledelayedexpansion

:: Config
set CONTAINER_NAME=kea-int-test
set IMAGE_NAME=kea-custom:latest
set KEA_DOCKER_PATH=..\kea-docker
set UPDATE_DOCKER_REPO=false

:: Check for -i flag
set RUN_INTEGRATION=false
for %%x in (%*) do (
    if "%%x"=="-i" set RUN_INTEGRATION=true
)

echo.
echo === Running unit tests ===
go test github.com/rannday/kea-api/client ^
        github.com/rannday/kea-api/agent ^
        github.com/rannday/kea-api/dhcp4 ^
        github.com/rannday/kea-api/dhcp6

if "%RUN_INTEGRATION%"=="true" (
    echo.
    echo === Preparing integration test environment ===

    :: Check if kea-docker repo is cloned
    if not exist "%KEA_DOCKER_PATH%\Dockerfile" (
        echo kea-docker not found. Cloning...
        git clone https://github.com/rannday/kea-docker.git %KEA_DOCKER_PATH%
        if errorlevel 1 (
            echo Failed to clone kea-docker repo. Aborting.
            exit /b 1
        )
    ) else (
        echo kea-docker repo found.
        if "%UPDATE_DOCKER_REPO%"=="true" (
            echo Updating kea-docker...
            pushd %KEA_DOCKER_PATH%
            git pull
            popd
        )
    )

    echo Checking for Docker image "%IMAGE_NAME%"...
    for /f %%i in ('docker images -q %IMAGE_NAME%') do set IMAGE_EXISTS=%%i

    if not defined IMAGE_EXISTS (
        echo Building Docker image from %KEA_DOCKER_PATH%...
        docker build -t %IMAGE_NAME% %KEA_DOCKER_PATH%
        if errorlevel 1 (
            echo Docker build failed. Aborting.
            exit /b 1
        )
    ) else (
        echo Image already built.
    )

    echo Checking if container "%CONTAINER_NAME%" is running...
    for /f %%i in ('docker ps -q --filter "name=%CONTAINER_NAME%"') do set CONTAINER_RUNNING=%%i

    if not defined CONTAINER_RUNNING (
        echo Starting Kea container...
        docker run -d --rm --name %CONTAINER_NAME% -p 8000:8000 %IMAGE_NAME%
        set CLEANUP_CONTAINER=true
        echo Waiting for Kea to be ready...
        timeout /t 3 >nul
    ) else (
        echo Container already running. Reusing it.
    )

    echo.
    echo === Running integration tests ===
    go test -tags=integration github.com/rannday/kea-api/client

    if defined CLEANUP_CONTAINER (
        echo.
        echo Stopping Kea container...
        docker stop %CONTAINER_NAME% >nul
    )
)

endlocal
