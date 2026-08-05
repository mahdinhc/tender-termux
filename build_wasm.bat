@echo off
echo Building Tender WebAssembly (playground/tender.wasm)...
set GOARCH=wasm
set GOOS=js
go build -o playground/tender.wasm playground/main.go
if %errorlevel% equ 0 (
    echo WASM build successful: playground/tender.wasm created!
) else (
    echo WASM build failed with error code %errorlevel%
)
