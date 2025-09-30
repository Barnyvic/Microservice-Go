# Quick Start Guide

This guide will help you get the Product Microservice up and running in minutes.

## Prerequisites Check

Before starting, ensure you have:

- [ ] Go 1.21+ installed (`go version`)
- [ ] Git installed
- [ ] Protocol Buffers compiler (`protoc --version`)

## Step 1: Install Protocol Buffer Plugins

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Make sure `$GOPATH/bin` is in your PATH:

**Windows (PowerShell):**
```powershell
$env:PATH += ";$env:USERPROFILE\go\bin"
```

**macOS/Linux:**
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## Step 2: Install Dependencies

```bash
go mod download
```

## Step 3: Generate Protocol Buffer Files

```bash
make proto
```

Or manually:
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/product.proto proto/subscription.proto
```

## Step 4: Run the Server

```bash
make run
```

Or:
```bash
go run cmd/server/main.go
```

You should see:
```
Database connection established
Running database migrations...
Database migrations completed successfully
gRPC server listening on port 50051
```

## Step 5: Test the API

### Option A: Using grpcurl (Recommended)

Install grpcurl:
- **macOS:** `brew install grpcurl`
- **Windows:** Download from [GitHub releases](https://github.com/fullstorydev/grpcurl/releases)
- **Linux:** Download from [GitHub releases](https://github.com/fullstorydev/grpcurl/releases)

List available services:
```bash
grpcurl -plaintext localhost:50051 list
```

Create a product:
```bash
grpcurl -plaintext -d '{
  "name": "Test Product",
  "description": "A test product",
  "price": 99.99,
  "product_type": "digital"
}' localhost:50051 product.ProductService/CreateProduct
```

### Option B: Using the Test Script

**Linux/macOS:**
```bash
chmod +x scripts/test_api.sh
./scripts/test_api.sh
```

**Windows:**
```cmd
scripts\test_api.bat
```

## Step 6: Run Tests

```bash
make test
```

Or:
```bash
go test -v ./...
```

## Common Issues

### Issue: "protoc-gen-go: program not found"

**Solution:** Install the plugins and add to PATH:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Issue: "cannot find package"

**Solution:** Run `go mod download` and `go mod tidy`

### Issue: Port 50051 already in use

**Solution:** Change the port:
```bash
export PORT=50052
go run cmd/server/main.go
```

### Issue: Database connection fails

**Solution:** For SQLite (default), ensure write permissions in the current directory. For PostgreSQL, verify connection settings.

## Next Steps

1. Read the full [README.md](README.md) for detailed documentation
2. Explore the [debug_exercise/](debug_exercise/) to learn about common issues
3. Check out the API examples in the README
4. Write your own tests
5. Extend the service with new features

## Project Structure Overview

```
├── cmd/server/main.go          # Start here - application entry point
├── proto/                      # gRPC service definitions
├── internal/
│   ├── models/                 # Database models
│   ├── repository/             # Data access layer
│   ├── service/                # Business logic
│   └── handler/                # gRPC handlers
└── debug_exercise/             # Learning materials
```

## Useful Commands

| Command | Description |
|---------|-------------|
| `make proto` | Generate protobuf files |
| `make build` | Build the server binary |
| `make run` | Run the server |
| `make test` | Run all tests |
| `make test-coverage` | Run tests with coverage report |
| `make clean` | Clean generated files |

## Learning Path

1. **Start with models** (`internal/models/`) - Understand the data structure
2. **Check repositories** (`internal/repository/`) - See how data is accessed
3. **Review services** (`internal/service/`) - Understand business logic
4. **Explore handlers** (`internal/handler/`) - See gRPC implementation
5. **Read tests** - Learn how everything works together

## Getting Help

- Check the [README.md](README.md) for detailed documentation
- Review the [debug_exercise/DEBUG_EXPLANATION.md](debug_exercise/DEBUG_EXPLANATION.md)
- Look at test files for usage examples
- Check GORM documentation: https://gorm.io/docs/
- Check gRPC documentation: https://grpc.io/docs/languages/go/

## Success Checklist

- [ ] Server starts without errors
- [ ] Can create a product via gRPC
- [ ] Can retrieve the created product
- [ ] Can create a subscription plan
- [ ] All tests pass
- [ ] Understand the project structure
- [ ] Read the debug exercise

Congratulations! You now have a working gRPC microservice. 🎉

