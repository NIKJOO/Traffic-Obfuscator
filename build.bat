@echo off
echo Building Winscribe Traffic Decoy...

:: Download dependencies
echo Downloading dependencies...
go mod download
go mod tidy

:: Build for Windows
echo Building for Windows...
go build -v -ldflags "-s -w" -o winscribe-decoy.exe .

echo Build complete! Binary: winscribe-decoy.exe