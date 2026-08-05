@echo off
echo Building Tender WebAssembly (tender.wasm)...
set GOOS=js
set GOARCH=wasm
go build -o tender.wasm main.go
if %errorlevel% equ 0 (
    echo WASM build successful: tender.wasm created!
) else (
    echo WASM build failed with error code %errorlevel%
)
