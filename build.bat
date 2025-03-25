@echo off
setlocal enabledelayedexpansion

go mod download
go mod verify
set CGO_ENABLED=0
set "PLATFORMS=windows/amd64 windows/arm64 linux/amd64 linux/arm64"
set "OUT_DIR=./bin"
set "MAIN_PATH=.\cmd\BaPs\BaPs.go"
set "OUTPUT_NAME=BaPs"

for %%p in (%PLATFORMS%) do (
    for /f "tokens=1,2 delims=/" %%a in ("%%p") do (
        set "GOOS=%%a"
        set "GOARCH=%%b"

        echo Compiling for GOOS=!GOOS! GOARCH=!GOARCH!...

        if not exist "!OUT_DIR!" mkdir "!OUT_DIR!"

        if "!GOOS!"=="windows" (
            go build -ldflags="-s -w" -o "!OUT_DIR!/!OUTPUT_NAME!_!GOOS!_!GOARCH!.exe" %MAIN_PATH%
        ) else (
            go build -ldflags="-s -w" -o "!OUT_DIR!/!OUTPUT_NAME!_!GOOS!_!GOARCH!" %MAIN_PATH%
        )
    )
)

endlocal
