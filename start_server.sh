#!/bin/bash

echo "========================================="
echo "  Building Product Microservice..."
echo "========================================="

export CGO_ENABLED=1

go build -o bin/server.exe cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "✓ Build successful!"
    echo ""
    echo "========================================="
    echo "  Starting Server..."
    echo "========================================="
    echo ""
    
    ./bin/server.exe
else
    echo " Build failed!"
    exit 1
fi

