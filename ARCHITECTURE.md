# Architecture Documentation

## System Architecture

### High-Level Overview

```
┌─────────────────────────────────────────────────────────────┐
│                        gRPC Clients                         │
│  (grpcurl, Go client, other language clients)               │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ gRPC Protocol (HTTP/2)
                         │
┌────────────────────────▼────────────────────────────────────┐
│                    gRPC Server (Port 50051)                 │
│                  (Reflection Enabled)                       │
└────────────────────────┬────────────────────────────────────┘
                         │
         ┌───────────────┴───────────────┐
         │                               │
┌────────▼──────────┐          ┌────────▼──────────┐
│ ProductService    │          │ SubscriptionService│
│ Handler           │          │ Handler            │
└────────┬──────────┘          └────────┬───────────┘
         │                               │
         │ Calls                         │ Calls
         │                               │
┌────────▼──────────┐          ┌────────▼───────────┐
│ ProductService    │          │ SubscriptionService│
│ (Business Logic)  │◄─────────┤ (Business Logic)   │
└────────┬──────────┘          └────────┬───────────┘
         │                               │
         │ Uses                          │ Uses
         │                               │
┌────────▼──────────┐          ┌────────▼───────────┐
│ ProductRepository │          │ SubscriptionRepo   │
│ (Data Access)     │          │ (Data Access)      │
└────────┬──────────┘          └────────┬───────────┘
         │                               │
         └───────────────┬───────────────┘
                         │
                         │ GORM
                         │
┌────────────────────────▼────────────────────────────────────┐
│                    Database (PostgreSQL/SQLite)             │
│  ┌──────────────┐              ┌──────────────────┐         │
│  │   products   │──────────────│ subscription_plans│        │
│  │              │  1:N         │                   │        │
│  └──────────────┘              └──────────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

## Layer Responsibilities

### 1. Handler Layer (Presentation)

**Location**: `internal/handler/`

**Responsibilities**:
- Receive gRPC requests
- Convert Protocol Buffer messages to domain models
- Call appropriate service methods
- Convert domain models back to Protocol Buffer messages
- Map errors to gRPC status codes
- Return gRPC responses

**Example Flow**:
```
gRPC Request → Handler.CreateProduct()
             → Convert PB to domain model
             → Call service.CreateProduct()
             → Convert domain model to PB
             → Return gRPC Response
```

**Files**:
- `product_handler.go` - Handles Product operations
- `subscription_handler.go` - Handles Subscription operations

### 2. Service Layer (Business Logic)

**Location**: `internal/service/`

**Responsibilities**:
- Validate input data
- Enforce business rules
- Orchestrate operations across repositories
- Handle business logic errors
- Return domain models

**Example Validations**:
- Product name cannot be empty
- Price cannot be negative
- Duration must be positive
- Product must exist before creating subscription plan

**Files**:
- `product_service.go` - Product business logic
- `subscription_service.go` - Subscription business logic

### 3. Repository Layer (Data Access)

**Location**: `internal/repository/`

**Responsibilities**:
- Execute database queries
- Handle GORM operations
- Manage transactions
- Convert database records to domain models
- Handle database-specific errors

**Operations**:
- CRUD operations
- Filtering and pagination
- Preloading associations
- Soft deletes

**Files**:
- `product_repository.go` - Product data access
- `subscription_repository.go` - Subscription data access

### 4. Model Layer (Domain)

**Location**: `internal/models/`

**Responsibilities**:
- Define database schema
- Specify GORM tags and constraints
- Define relationships
- Implement hooks (BeforeCreate, etc.)

**Files**:
- `product.go` - Product model
- `subscription_plan.go` - SubscriptionPlan model

## Data Flow

### Creating a Product

```
1. Client sends CreateProductRequest via gRPC
   ↓
2. ProductHandler.CreateProduct() receives request
   ↓
3. Handler extracts data from Protocol Buffer
   ↓
4. Handler calls ProductService.CreateProduct()
   ↓
5. Service validates input (name, price, type)
   ↓
6. Service creates Product model
   ↓
7. Service calls ProductRepository.Create()
   ↓
8. Repository executes GORM Create operation
   ↓
9. GORM BeforeCreate hook generates UUID
   ↓
10. Database inserts record
    ↓
11. Repository returns Product model
    ↓
12. Service returns Product model
    ↓
13. Handler converts to Protocol Buffer
    ↓
14. Handler returns CreateProductResponse
    ↓
15. Client receives response
```

### Creating a Subscription Plan (with validation)

```
1. Client sends CreateSubscriptionPlanRequest
   ↓
2. SubscriptionHandler receives request
   ↓
3. Handler calls SubscriptionService.CreateSubscriptionPlan()
   ↓
4. Service validates input
   ↓
5. Service calls ProductRepository.GetByID() to verify product exists
   ↓
6. If product not found → return error
   ↓
7. If product exists → create SubscriptionPlan model
   ↓
8. Service calls SubscriptionRepository.Create()
   ↓
9. Repository executes GORM Create with foreign key
   ↓
10. Database inserts record with FK constraint
    ↓
11. Returns through layers to client
```

## Database Relationships

### Entity Relationship Diagram

```
┌─────────────────────────────────────┐
│            Product                  │
├─────────────────────────────────────┤
│ id: UUID (PK)                       │
│ name: VARCHAR                       │
│ description: TEXT                   │
│ price: DECIMAL                      │
│ product_type: VARCHAR               │
│ created_at: TIMESTAMP               │
│ updated_at: TIMESTAMP               │
│ deleted_at: TIMESTAMP (nullable)    │
└──────────────┬──────────────────────┘
               │
               │ 1:N (One-to-Many)
               │ ON DELETE CASCADE
               │
┌──────────────▼──────────────────────┐
│       SubscriptionPlan              │
├─────────────────────────────────────┤
│ id: UUID (PK)                       │
│ product_id: UUID (FK) → Product.id  │
│ plan_name: VARCHAR                  │
│ duration: INTEGER                   │
│ price: DECIMAL                      │
│ created_at: TIMESTAMP               │
│ updated_at: TIMESTAMP               │
│ deleted_at: TIMESTAMP (nullable)    │
└─────────────────────────────────────┘
```

### Relationship Details

**Product → SubscriptionPlan**:
- Type: One-to-Many
- Foreign Key: `subscription_plans.product_id` → `products.id`
- Cascade: ON DELETE CASCADE
- GORM Tag: `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`

**Benefits**:
- Referential integrity enforced at database level
- Automatic cleanup of orphaned subscription plans
- Efficient queries with preloading

## Interface Design

### Repository Interfaces

```go
type ProductRepository interface {
    Create(product *models.Product) error
    GetByID(id uuid.UUID) (*models.Product, error)
    Update(product *models.Product) error
    Delete(id uuid.UUID) error
    List(productType string, page, pageSize int) ([]models.Product, int64, error)
}

type SubscriptionRepository interface {
    Create(plan *models.SubscriptionPlan) error
    GetByID(id uuid.UUID) (*models.SubscriptionPlan, error)
    Update(plan *models.SubscriptionPlan) error
    Delete(id uuid.UUID) error
    ListByProductID(productID uuid.UUID) ([]models.SubscriptionPlan, error)
}
```

### Service Interfaces

```go
type ProductService interface {
    CreateProduct(name, description string, price float64, productType string) (*models.Product, error)
    GetProduct(id string) (*models.Product, error)
    UpdateProduct(id, name, description string, price float64, productType string) (*models.Product, error)
    DeleteProduct(id string) error
    ListProducts(productType string, page, pageSize int) ([]models.Product, int64, error)
}

type SubscriptionService interface {
    CreateSubscriptionPlan(productID, planName string, duration int, price float64) (*models.SubscriptionPlan, error)
    GetSubscriptionPlan(id string) (*models.SubscriptionPlan, error)
    UpdateSubscriptionPlan(id, productID, planName string, duration int, price float64) (*models.SubscriptionPlan, error)
    DeleteSubscriptionPlan(id string) error
    ListSubscriptionPlans(productID string) ([]models.SubscriptionPlan, error)
}
```

## Error Handling Strategy

### Error Flow

```
Database Error
    ↓
Repository catches and returns Go error
    ↓
Service receives error and may wrap with context
    ↓
Handler receives error and maps to gRPC status code
    ↓
Client receives gRPC error with appropriate code
```

### Error Mapping

| Error Type | gRPC Code | Example |
|------------|-----------|---------|
| Not Found | `codes.NotFound` | Product/Plan doesn't exist |
| Invalid Input | `codes.InvalidArgument` | Empty name, negative price |
| Internal Error | `codes.Internal` | Database connection failed |
| Already Exists | `codes.AlreadyExists` | Duplicate entry |

## Testing Architecture

### Test Pyramid

```
        ┌─────────────┐
        │   Handler   │  ← Few tests (integration-like)
        │    Tests    │     Mock services
        └─────────────┘
       ┌───────────────┐
       │   Service     │  ← More tests (unit)
       │    Tests      │     Mock repositories
       └───────────────┘
      ┌─────────────────┐
      │   Repository    │  ← Most tests (integration)
      │     Tests       │     In-memory database
      └─────────────────┘
```

### Test Strategy by Layer

**Handler Tests**:
- Mock service layer
- Test Protocol Buffer conversion
- Test error code mapping
- Verify correct service methods are called

**Service Tests**:
- Mock repository layer
- Test business logic validation
- Test error handling
- Test orchestration logic

**Repository Tests**:
- Use in-memory SQLite
- Test actual database operations
- Test GORM features (preload, soft delete)
- Test pagination and filtering

## Deployment Architecture

### Single Instance Deployment

```
┌─────────────────────────────────────┐
│         Load Balancer               │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│      gRPC Server Instance           │
│  ┌────────────────────────────┐     │
│  │  Product Microservice      │     │
│  │  (Port 50051)              │     │
│  └────────────────────────────┘     │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│      PostgreSQL Database            │
└─────────────────────────────────────┘
```

### Scalable Deployment

```
┌─────────────────────────────────────┐
│         Load Balancer               │
└──────────┬──────────┬───────────────┘
           │          │
    ┌──────▼───┐  ┌──▼──────┐
    │ Instance │  │ Instance│  ... (N instances)
    │    1     │  │    2    │
    └──────┬───┘  └──┬──────┘
           │         │
           └────┬────┘
                │
    ┌───────────▼──────────────┐
    │  PostgreSQL (Primary)    │
    │  with Read Replicas      │
    └──────────────────────────┘
```

## Security Considerations

### Current Implementation
- ✅ SQL Injection protection (GORM parameterized queries)
- ✅ UUID primary keys (prevent enumeration)
- ✅ Input validation at service layer
- ✅ Soft deletes (data recovery)

### Production Additions Needed
- 🔲 Authentication (JWT tokens)
- 🔲 Authorization (role-based access)
- 🔲 Rate limiting
- 🔲 TLS/SSL for gRPC
- 🔲 API key management
- 🔲 Audit logging

## Performance Considerations

### Optimizations Implemented
- Database indexes on foreign keys
- Pagination to limit result sets
- Preloading to avoid N+1 queries
- Connection pooling (GORM default)

### Future Optimizations
- Caching layer (Redis)
- Database read replicas
- Query optimization
- Connection pool tuning
- Batch operations

## Monitoring and Observability

### Recommended Additions

```
┌─────────────────────────────────────┐
│         Monitoring Stack            │
├─────────────────────────────────────┤
│  Metrics:    Prometheus             │
│  Logging:    ELK Stack              │
│  Tracing:    Jaeger/Zipkin          │
│  Alerting:   AlertManager           │
└─────────────────────────────────────┘
```

### Key Metrics to Track
- Request rate (requests/second)
- Error rate (errors/total requests)
- Latency (p50, p95, p99)
- Database connection pool usage
- Active gRPC connections

---

This architecture provides a solid foundation for a production microservice with clear separation of concerns, testability, and scalability.

