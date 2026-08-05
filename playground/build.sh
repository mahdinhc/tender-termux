#!/bin/bash
echo "Building Tender WebAssembly (tender.wasm)..."
GOOS=js GOARCH=wasm go build -o tender.wasm main.go
if [ $? -eq 0 ]; then
    echo "WASM build successful: tender.wasm created!"
else
    echo "WASM build failed."
fi
