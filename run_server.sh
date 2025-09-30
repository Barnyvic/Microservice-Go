#!/bin/bash
cd "$(dirname "$0")"
echo "Starting gRPC server..."
go run cmd/server/main.go

