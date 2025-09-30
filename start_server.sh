#!/bin/bash

echo "========================================="
echo "  Building Product Microservice..."
echo "========================================="

# Enable CGO for SQLite support
export CGO_ENABLED=1

# Build the server
go build -o bin/server.exe cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "✓ Build successful!"
    echo ""
    echo "========================================="
    echo "  Starting Server..."
    echo "========================================="
    echo ""
    
    # Run the server
    ./bin/server.exe
else
    echo "✗ Build failed!"
    exit 1
fi

