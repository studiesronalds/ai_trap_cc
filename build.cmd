@echo off
setlocal enabledelayedexpansion

set BINARY=aisim
set PKG=.\cmd\aisim
set DIST=dist

if exist "%DIST%" rd /s /q "%DIST%"
mkdir "%DIST%"

call :build linux   amd64
call :build linux   arm64
call :build darwin  amd64
call :build darwin  arm64
call :build windows amd64
call :build windows arm64

echo.
echo Done. Binaries in %DIST%\:
dir /b "%DIST%"
goto :eof

:build
set GOOS=%1
set GOARCH=%2
set OUT=%DIST%\%BINARY%-%GOOS%-%GOARCH%
if "%GOOS%"=="windows" set OUT=%OUT%.exe
echo Building %GOOS%/%GOARCH% ^> %OUT%
go build -o "%OUT%" %PKG%
goto :eof
