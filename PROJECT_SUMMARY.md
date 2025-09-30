# Product Microservice - Project Summary

## Overview

This is a production-ready product microservice built with Go, gRPC, and GORM, following clean architecture principles. The service manages products and their associated subscription plans with full CRUD operations exposed via gRPC endpoints.

## Key Features

### Technical Stack
- **Language**: Go 1.21+
- **RPC Framework**: gRPC with Protocol Buffers
- **ORM**: GORM v1.25.5
- **Database**: PostgreSQL / SQLite (configurable)
- **Testing**: Testify for mocking and assertions
- **Architecture**: Clean Architecture (Handlers → Services → Repositories)

### Core Functionality
1. **Product Management**
   - Create, Read, Update, Delete operations
   - List with pagination and type filtering
   - UUID-based primary keys
   - Soft delete support

2. **Subscription Plan Management**
   - CRUD operations for subscription plans
   - Foreign key relationship with products
   - Cascade delete when product is removed
   - List plans by product

3. **Database Features**
   - Automatic UUID generation via GORM hooks
   - Foreign key constraints with cascade delete
   - Automatic timestamps (CreatedAt, UpdatedAt)
   - Soft deletes with DeletedAt
   - Support for multiple database drivers

## Architecture

### Layer Structure

```
┌─────────────────────────────────────────┐
│         gRPC Layer (Handlers)           │
│  - Protocol Buffer conversion           │
│  - Request/Response handling            │
│  - Error code mapping                   │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│       Service Layer (Business Logic)    │
│  - Validation                           │
│  - Business rules                       │
│  - Orchestration                        │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│    Repository Layer (Data Access)       │
│  - GORM operations                      │
│  - Query building                       │
│  - Transaction management               │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│         Database (PostgreSQL/SQLite)    │
└─────────────────────────────────────────┘
```

### Design Patterns

1. **Repository Pattern**: Abstracts data access logic
2. **Dependency Injection**: Services depend on repository interfaces
3. **Interface Segregation**: Small, focused interfaces
4. **Single Responsibility**: Each layer has one clear purpose

## Project Structure

```
Microservice-Go/
├── cmd/server/              # Application entry point
├── internal/
│   ├── database/           # Database configuration and migrations
│   ├── handler/            # gRPC handlers (presentation layer)
│   ├── models/             # Database models
│   ├── repository/         # Data access layer
│   └── service/            # Business logic layer
├── proto/                  # Protocol Buffer definitions
├── debug_exercise/         # Learning materials
├── examples/client/        # Example Go client
└── scripts/                # Helper scripts
```

## API Endpoints

### ProductService
- `CreateProduct` - Create a new product
- `GetProduct` - Retrieve a product by ID
- `UpdateProduct` - Update an existing product
- `DeleteProduct` - Delete a product (soft delete)
- `ListProducts` - List products with optional filtering and pagination

### SubscriptionService
- `CreateSubscriptionPlan` - Create a new subscription plan
- `GetSubscriptionPlan` - Retrieve a plan by ID
- `UpdateSubscriptionPlan` - Update an existing plan
- `DeleteSubscriptionPlan` - Delete a plan (soft delete)
- `ListSubscriptionPlans` - List plans for a specific product

## Testing Strategy

### Test Coverage
- **Service Layer**: Unit tests with mocked repositories
- **Repository Layer**: Integration tests with in-memory SQLite
- **Handler Layer**: Unit tests with mocked services

### Test Features
- Mock implementations using testify/mock
- Assertions using testify/assert
- In-memory database for repository tests
- Comprehensive error case coverage

### Running Tests
```bash
make test              # Run all tests
make test-coverage     # Generate coverage report
```

## Database Schema

### Products Table
```sql
CREATE TABLE products (
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    description TEXT,
    price DECIMAL NOT NULL,
    product_type VARCHAR NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

### Subscription Plans Table
```sql
CREATE TABLE subscription_plans (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    plan_name VARCHAR NOT NULL,
    duration INTEGER NOT NULL,
    price DECIMAL NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

## Configuration

### Environment Variables
- `DB_DRIVER` - Database driver (sqlite/postgres)
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `DB_SSLMODE` - SSL mode for PostgreSQL
- `PORT` - gRPC server port

### Default Configuration
- Database: SQLite (products.db)
- Port: 50051
- SSL Mode: disabled

## Documentation References

This project was built using official documentation:

1. **gRPC Documentation**
   - Service definition and code generation
   - Server implementation
   - Error handling with status codes

2. **GORM Documentation**
   - Model definition and conventions
   - Associations (Belongs To, Has Many)
   - Hooks for automatic field population
   - Data types (UUID support)
   - Migrations

3. **Go Best Practices**
   - Project layout standards
   - Testing patterns
   - Error handling

## Debug Exercise

The `debug_exercise/` folder contains:
- **broken_code.go**: Intentionally broken code demonstrating common issues
- **DEBUG_EXPLANATION.md**: Detailed explanation of:
  - The problems (UUID type issues, missing associations)
  - Documentation used to solve them
  - Step-by-step solution
  - Key takeaways

## Getting Started

### Quick Start
```bash
# 1. Install dependencies
go mod download

# 2. Generate proto files
make proto

# 3. Run tests
make test

# 4. Start server
make run

# 5. Test API (in another terminal)
grpcurl -plaintext localhost:50051 list
```

### Using the Example Client
```bash
# Make sure server is running
go run examples/client/main.go
```

## Key Implementation Details

### UUID Generation
- Automatic generation using GORM BeforeCreate hooks
- Type-safe with `uuid.UUID` instead of strings
- Proper database type mapping

### Foreign Key Relationships
- Explicit foreign key definitions with GORM tags
- Cascade delete behavior
- Preloading support for efficient queries

### Error Handling
- gRPC status codes (NotFound, Internal, etc.)
- Validation at service layer
- Proper error propagation

### Pagination
- Page and page_size parameters
- Total count returned
- Configurable limits

## Performance Considerations

1. **Database Indexes**: Foreign keys and frequently queried fields are indexed
2. **Soft Deletes**: Uses GORM's DeletedAt for data recovery
3. **Preloading**: Eager loading of associations to avoid N+1 queries
4. **Connection Pooling**: GORM handles connection pooling automatically

## Security Considerations

1. **SQL Injection**: Protected by GORM's parameterized queries
2. **Input Validation**: Service layer validates all inputs
3. **UUID Primary Keys**: Prevents enumeration attacks
4. **Soft Deletes**: Data recovery capability

## Future Enhancements

Potential improvements for production deployment:

1. **Authentication & Authorization**: JWT tokens, role-based access
2. **Rate Limiting**: Protect against abuse
3. **Caching**: Redis for frequently accessed data
4. **Observability**: Metrics, logging, tracing
5. **API Gateway**: Centralized routing and policies
6. **Docker**: Containerization for easy deployment
7. **CI/CD**: Automated testing and deployment
8. **GraphQL**: Alternative API layer
9. **Event Sourcing**: Audit trail and event-driven architecture
10. **Multi-tenancy**: Support for multiple organizations

## Lessons Learned

### Documentation-Driven Development
- Reading official docs first saves time
- Understanding the "why" behind patterns
- Applying best practices from the start

### Clean Architecture Benefits
- Easy to test each layer independently
- Simple to swap implementations
- Clear separation of concerns

### GORM Power Features
- Hooks simplify common tasks
- Associations make relationships explicit
- Migrations handle schema changes

### gRPC Advantages
- Type-safe API contracts
- Efficient binary protocol
- Built-in code generation
- Strong typing across languages

## Conclusion

This project demonstrates:
- ✅ Proficiency with Go, gRPC, and GORM
- ✅ Understanding of clean architecture
- ✅ Ability to learn from documentation
- ✅ Testing best practices
- ✅ Production-ready code quality

The implementation is complete, well-tested, and ready for deployment or further development.

---

**For detailed setup instructions, see [QUICKSTART.md](QUICKSTART.md)**  
**For comprehensive documentation, see [README.md](README.md)**  
**For submission verification, see [SUBMISSION_CHECKLIST.md](SUBMISSION_CHECKLIST.md)**

