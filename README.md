# Product Microservice with gRPC

A production-ready product microservice built with Go, gRPC, and GORM following clean architecture principles.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Running the Server](#running-the-server)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Documentation Used](#documentation-used)
- [Project Structure](#project-structure)

## Features

- **gRPC API** for Product and Subscription management
- **Clean Architecture** with separation of concerns (handlers, services, repositories)
- **GORM** for database operations with support for PostgreSQL and SQLite
- **UUID** primary keys with automatic generation
- **Foreign key relationships** with cascade delete
- **Comprehensive testing** (unit tests for services, repositories, and handlers)
- **Database migrations** with GORM AutoMigrate
- **Pagination** support for list operations
- **Type filtering** for products

## Architecture

The project follows clean architecture principles with clear separation of concerns:

```
┌─────────────────┐
│   gRPC Layer    │  (Handlers)
│  Product/Sub    │
└────────┬────────┘
         │
┌────────▼────────┐
│  Service Layer  │  (Business Logic)
│  Validation     │
└────────┬────────┘
         │
┌────────▼────────┐
│ Repository Layer│  (Data Access)
│      GORM       │
└────────┬────────┘
         │
┌────────▼────────┐
│    Database     │  (PostgreSQL/SQLite)
└─────────────────┘
```

### Layers:

- **Handlers**: gRPC request/response handling and protocol buffer conversion
- **Services**: Business logic, validation, and orchestration
- **Repositories**: Data access layer with GORM
- **Models**: Database entities with GORM tags

## Prerequisites

- Go 1.21 or higher
- Protocol Buffers compiler (`protoc`)
- PostgreSQL (optional, SQLite is default)
- `grpcurl` for testing (optional)

### Install Protocol Buffer Compiler

**Windows:**

```bash
# Download from https://github.com/protocolbuffers/protobuf/releases
# Add to PATH
```

**macOS:**

```bash
brew install protobuf
```

**Linux:**

```bash
sudo apt install -y protobuf-compiler
```

### Install Go Plugins for Protocol Buffers

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Make sure `$GOPATH/bin` is in your PATH.

## Installation

1. **Clone the repository:**

```bash
git clone <repository-url>
cd Microservice-Go
```

2. **Install dependencies:**

```bash
go mod download
```

3. **Generate Protocol Buffer files:**

```bash
make proto
```

Or manually:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/product.proto proto/subscription.proto
```

## Running the Server

### Using SQLite (Default)

```bash
make run
```

Or:

```bash
go run cmd/server/main.go
```

The server will start on port `50051` by default.

### Using PostgreSQL

Set environment variables:

```bash
export DB_DRIVER=postgres
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=yourpassword
export DB_NAME=products_db
export DB_SSLMODE=disable

go run cmd/server/main.go
```

### Environment Variables

| Variable      | Default       | Description                              |
| ------------- | ------------- | ---------------------------------------- |
| `DB_DRIVER`   | `sqlite`      | Database driver (`sqlite` or `postgres`) |
| `DB_HOST`     | `localhost`   | Database host                            |
| `DB_PORT`     | `5432`        | Database port                            |
| `DB_USER`     | `postgres`    | Database user                            |
| `DB_PASSWORD` | `postgres`    | Database password                        |
| `DB_NAME`     | `products.db` | Database name                            |
| `DB_SSLMODE`  | `disable`     | SSL mode for PostgreSQL                  |
| `PORT`        | `50051`       | gRPC server port                         |

## API Documentation

### Product Service

#### CreateProduct

```bash
grpcurl -plaintext -d '{
  "name": "Premium Software",
  "description": "Enterprise software solution",
  "price": 299.99,
  "product_type": "digital"
}' localhost:50051 product.ProductService/CreateProduct
```

#### GetProduct

```bash
grpcurl -plaintext -d '{
  "id": "your-product-uuid"
}' localhost:50051 product.ProductService/GetProduct
```

#### UpdateProduct

```bash
grpcurl -plaintext -d '{
  "id": "your-product-uuid",
  "name": "Updated Product Name",
  "description": "Updated description",
  "price": 349.99,
  "product_type": "digital"
}' localhost:50051 product.ProductService/UpdateProduct
```

#### DeleteProduct

```bash
grpcurl -plaintext -d '{
  "id": "your-product-uuid"
}' localhost:50051 product.ProductService/DeleteProduct
```

#### ListProducts

```bash
# List all products
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10
}' localhost:50051 product.ProductService/ListProducts

# Filter by product type
grpcurl -plaintext -d '{
  "product_type": "digital",
  "page": 1,
  "page_size": 10
}' localhost:50051 product.ProductService/ListProducts
```

### Subscription Service

#### CreateSubscriptionPlan

```bash
grpcurl -plaintext -d '{
  "product_id": "your-product-uuid",
  "plan_name": "Monthly Plan",
  "duration": 30,
  "price": 29.99
}' localhost:50051 subscription.SubscriptionService/CreateSubscriptionPlan
```

#### GetSubscriptionPlan

```bash
grpcurl -plaintext -d '{
  "id": "your-plan-uuid"
}' localhost:50051 subscription.SubscriptionService/GetSubscriptionPlan
```

#### UpdateSubscriptionPlan

```bash
grpcurl -plaintext -d '{
  "id": "your-plan-uuid",
  "product_id": "your-product-uuid",
  "plan_name": "Annual Plan",
  "duration": 365,
  "price": 299.99
}' localhost:50051 subscription.SubscriptionService/UpdateSubscriptionPlan
```

#### DeleteSubscriptionPlan

```bash
grpcurl -plaintext -d '{
  "id": "your-plan-uuid"
}' localhost:50051 subscription.SubscriptionService/DeleteSubscriptionPlan
```

#### ListSubscriptionPlans

```bash
grpcurl -plaintext -d '{
  "product_id": "your-product-uuid"
}' localhost:50051 subscription.SubscriptionService/ListSubscriptionPlans
```

### List Available Services

```bash
grpcurl -plaintext localhost:50051 list
```

### Describe a Service

```bash
grpcurl -plaintext localhost:50051 describe product.ProductService
```

### Example Go Client

A complete example client is provided in `examples/client/main.go`. To run it:

```bash
# Make sure the server is running first
go run examples/client/main.go
```

This example demonstrates:

- Creating products and subscription plans
- Retrieving and updating records
- Listing with filters and pagination
- Deleting records with cascade behavior

## Testing

### Run All Tests

```bash
make test
```

Or:

```bash
go test -v ./...
```

### Run Tests with Coverage

```bash
make test-coverage
```

This generates `coverage.html` that you can open in a browser.

### Test Structure

- **Service Tests** (`internal/service/*_test.go`): Test business logic with mocked repositories
- **Repository Tests** (`internal/repository/*_test.go`): Test data access with in-memory SQLite
- **Handler Tests** (`internal/handler/*_test.go`): Test gRPC handlers with mocked services

## Documentation Used

This section explains the official documentation consulted and how it helped in building this microservice.

### 1. gRPC Go Quick Start

**URL:** https://grpc.io/docs/languages/go/quickstart/

**How It Helped:**

- Learned how to define `.proto` files with service definitions
- Understood how to generate Go code from proto files using `protoc`
- Learned about `protoc-gen-go` and `protoc-gen-go-grpc` plugins
- Understood the structure of gRPC servers and how to register services

**Application:**

- Created `proto/product.proto` and `proto/subscription.proto` with proper service definitions
- Set up Makefile commands for proto generation
- Implemented server registration in `cmd/server/main.go`

### 2. gRPC Go Basics Tutorial

**URL:** https://grpc.io/docs/languages/go/basics/

**How It Helped:**

- Learned about unary RPC patterns (request-response)
- Understood how to implement server-side handlers
- Learned about context usage in gRPC
- Understood error handling with `status` package and gRPC codes

**Application:**

- Implemented handlers in `internal/handler/` with proper error codes
- Used `codes.NotFound`, `codes.Internal` for appropriate error scenarios
- Implemented context-aware handlers

### 3. GORM Documentation - Declaring Models

**URL:** https://gorm.io/docs/models.html

**How It Helped:**

- Learned about GORM conventions (ID, CreatedAt, UpdatedAt, DeletedAt)
- Understood struct tags for database constraints
- Learned about soft deletes with `gorm.DeletedAt`

**Application:**

- Created models in `internal/models/` with proper GORM tags
- Implemented soft delete functionality
- Used `gorm:"not null"`, `gorm:"index"` tags appropriately

### 4. GORM Associations - Belongs To & Has Many

**URL:** https://gorm.io/docs/belongs_to.html and https://gorm.io/docs/has_many.html

**How It Helped:**

- Learned how to define one-to-many relationships
- Understood foreign key configuration with `foreignKey` and `references` tags
- Learned about cascade delete with `constraint:OnDelete:CASCADE`
- Understood how to use `Preload` for eager loading

**Application:**

- Implemented Product → SubscriptionPlan one-to-many relationship
- Used proper foreign key constraints in models
- Implemented `Preload("SubscriptionPlans")` in repository layer

### 5. GORM Hooks

**URL:** https://gorm.io/docs/hooks.html

**How It Helped:**

- Learned about lifecycle hooks (BeforeCreate, BeforeUpdate, etc.)
- Understood how to use hooks for automatic field population
- Learned best practices for hook implementation

**Application:**

- Implemented `BeforeCreate` hooks in both Product and SubscriptionPlan models
- Automatic UUID generation before record creation
- Ensured consistency across all database operations

### 6. GORM Data Types

**URL:** https://gorm.io/docs/data_types.html

**How It Helped:**

- Learned about UUID support in GORM
- Understood how to use `uuid.UUID` type with proper tags
- Learned about type-specific database configurations

**Application:**

- Used `uuid.UUID` type with `gorm:"type:uuid"` tag
- Proper UUID handling in PostgreSQL and SQLite
- Type-safe foreign key relationships

### 7. GORM Migrations

**URL:** https://gorm.io/docs/migration.html

**How It Helped:**

- Learned about `AutoMigrate` for automatic schema creation
- Understood migration best practices
- Learned about index creation and constraints

**Application:**

- Implemented `RunMigrations` function in `internal/database/database.go`
- Automatic table creation with proper constraints
- Index creation for foreign keys and frequently queried fields

### 8. Go Project Layout

**URL:** https://github.com/golang-standards/project-layout

**How It Helped:**

- Understood standard Go project structure
- Learned about `cmd/`, `internal/`, `pkg/` directories
- Best practices for organizing Go applications

**Application:**

- Organized code into `cmd/server/` for main application
- Used `internal/` for private application code
- Separated concerns into handlers, services, repositories, models

### 9. Protocol Buffers Language Guide

**URL:** https://protobuf.dev/programming-guides/proto3/

**How It Helped:**

- Learned proto3 syntax and conventions
- Understood message definitions and field types
- Learned about importing well-known types (timestamp)

**Application:**

- Created well-structured proto files with proper field numbering
- Used `google.protobuf.Timestamp` for time fields
- Proper use of repeated fields for lists

### 10. Testify Documentation

**URL:** https://github.com/stretchr/testify

**How It Helped:**

- Learned about mock creation with `testify/mock`
- Understood assertion helpers with `testify/assert`
- Best practices for writing testable Go code

**Application:**

- Created mock implementations for repositories and services
- Used `assert` package for clean test assertions
- Implemented comprehensive test coverage

## Project Structure

```
Microservice-Go/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── database/
│   │   └── database.go             # Database connection and migrations
│   ├── handler/
│   │   ├── product_handler.go      # gRPC Product handler
│   │   ├── product_handler_test.go
│   │   ├── subscription_handler.go # gRPC Subscription handler
│   │   └── subscription_handler_test.go
│   ├── models/
│   │   ├── product.go              # Product model
│   │   └── subscription_plan.go    # SubscriptionPlan model
│   ├── repository/
│   │   ├── product_repository.go   # Product data access
│   │   ├── product_repository_test.go
│   │   ├── subscription_repository.go
│   │   └── subscription_repository_test.go
│   └── service/
│       ├── product_service.go      # Product business logic
│       ├── product_service_test.go
│       ├── subscription_service.go
│       └── subscription_service_test.go
├── proto/
│   ├── product.proto               # Product service definition
│   ├── product.pb.go               # Generated (not in repo)
│   ├── product_grpc.pb.go          # Generated (not in repo)
│   ├── subscription.proto          # Subscription service definition
│   ├── subscription.pb.go          # Generated (not in repo)
│   └── subscription_grpc.pb.go     # Generated (not in repo)
├── debug_exercise/
│   ├── broken_code.go              # Intentionally broken code
│   └── DEBUG_EXPLANATION.md        # Debugging write-up
├── go.mod
├── go.sum
├── Makefile
├── .gitignore
└── README.md
```

## Key Design Decisions

### 1. Clean Architecture

- **Why**: Separation of concerns, testability, maintainability
- **How**: Layered architecture with clear boundaries between handlers, services, and repositories

### 2. UUID Primary Keys

- **Why**: Distributed system friendly, no auto-increment conflicts, better for microservices
- **How**: Used `github.com/google/uuid` with GORM hooks for automatic generation

### 3. Repository Pattern

- **Why**: Abstracts data access, makes testing easier, allows switching databases
- **How**: Interface-based repositories with GORM implementations

### 4. Soft Deletes

- **Why**: Data recovery, audit trails, safer operations
- **How**: GORM's `DeletedAt` field with automatic handling

### 5. Pagination

- **Why**: Performance with large datasets, better API design
- **How**: Page and page_size parameters in List operations

## Common Issues and Solutions

### Issue: Proto files not generating

**Solution:** Ensure `protoc` and Go plugins are installed and in PATH

```bash
make install-proto
make proto
```

### Issue: Database connection fails

**Solution:** Check environment variables and database availability

```bash
# For SQLite, ensure write permissions
# For PostgreSQL, verify connection string
```

### Issue: Tests fail with "database locked"

**Solution:** Tests use in-memory SQLite, ensure no concurrent access

## Future Enhancements

- [ ] Add authentication and authorization
- [ ] Implement rate limiting
- [ ] Add caching layer (Redis)
- [ ] Implement event sourcing
- [ ] Add observability (metrics, tracing)
- [ ] Docker containerization
- [ ] Kubernetes deployment manifests
- [ ] CI/CD pipeline
- [ ] API gateway integration
- [ ] GraphQL layer

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run tests and ensure they pass
6. Submit a pull request

## License

MIT License - feel free to use this project for learning and development.

## Contact

For questions or feedback, please open an issue in the repository.
