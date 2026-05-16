@echo off
echo Fetching dependencies...
go mod tidy
if %errorlevel% neq 0 (
    echo ERROR: go mod tidy failed
    exit /b 1
)
echo.
echo Building CLI tool...
go build -o make_colors.exe ./cmd/make_colors
if %errorlevel% neq 0 (
    echo ERROR: build failed
    exit /b 1
)
echo.
echo Done! Run: make_colors.exe -t
