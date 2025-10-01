# Debug Exercise: Fixing Common GORM and UUID Issues

## Problem Statement

The `broken_code.go` file contains a Product microservice implementation with multiple critical issues related to GORM usage, UUID handling, and database relationships. When running this code, you'll encounter several problems that prevent the application from working correctly.

---

## Issues Identified

### Issue 1: Using String Instead of UUID Type

**Problem:**
```go
type Product struct {
    ID string  
    ...
}
```

**Symptoms:**
- Product IDs are empty strings after creation
- Cannot reliably query products by ID
- No UUID validation

**Documentation Consulted:**
- **GORM Data Types**: https://gorm.io/docs/data_types.html
- **google/uuid Package**: https://pkg.go.dev/github.com/google/uuid

**Solution:**
```go
import "github.com/google/uuid"

type Product struct {
    ID uuid.UUID `gorm:"type:uuid;primary_key"`
    ...
}
```

**Why This Works:**
- `uuid.UUID` is a proper type for UUIDs with built-in validation
- GORM recognizes this type and creates appropriate database columns
- The `gorm:"type:uuid"` tag ensures correct database type (especially for PostgreSQL)

---

### Issue 2: Missing GORM Convention Fields

**Problem:**
```go
type Product struct {
    ID   uuid.UUID
    Name string
}
```

**Symptoms:**
- No automatic timestamp tracking
- Cannot use soft deletes
- No audit trail for when records were created/modified

**Documentation Consulted:**
- **GORM Models**: https://gorm.io/docs/models.html
- **GORM Delete**: https://gorm.io/docs/delete.html#Soft-Delete

**Solution:**
```go
import (
    "time"
    "gorm.io/gorm"
)

type Product struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key"`
    Name      string    `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

**Why This Works:**
- GORM automatically populates `CreatedAt` on record creation
- GORM automatically updates `UpdatedAt` on record modification
- `gorm.DeletedAt` enables soft deletes (records marked as deleted but not removed)
- The `index` tag on DeletedAt improves query performance

---

### Issue 3: Missing UUID Auto-Generation

**Problem:**
```go
product := Product{
    Name: "Test Product",
}
db.Create(&product)
```

**Symptoms:**
- Products created with nil/zero UUIDs
- Duplicate key violations if multiple records created
- Cannot reference products by ID

**Documentation Consulted:**
- **GORM Hooks**: https://gorm.io/docs/hooks.html
- **UUID Generation**: https://pkg.go.dev/github.com/google/uuid#New

**Solution:**
```go
func (p *Product) BeforeCreate(tx *gorm.DB) error {
    if p.ID == uuid.Nil {
        p.ID = uuid.New()
    }
    return nil
}
```

**Why This Works:**
- `BeforeCreate` is a GORM lifecycle hook called before inserting a record
- Automatically generates a new UUID if one isn't provided
- Ensures every product has a unique identifier
- The `uuid.Nil` check allows manual UUID assignment if needed

---

### Issue 4: Missing Foreign Key Constraints

**Problem:**
```go
type SubscriptionPlan struct {
    ProductID string 
}

type Product struct {
    SubscriptionPlans []SubscriptionPlan  
}
```

**Symptoms:**
- Can create subscription plans with invalid ProductIDs
- No referential integrity
- Deleting products leaves orphaned subscription plans
- Associations don't load automatically

**Documentation Consulted:**
- **GORM Belongs To**: https://gorm.io/docs/belongs_to.html
- **GORM Has Many**: https://gorm.io/docs/has_many.html
- **GORM Constraints**: https://gorm.io/docs/constraints.html

**Solution:**
```go
type Product struct {
    ID                uuid.UUID           `gorm:"type:uuid;primary_key"`
    SubscriptionPlans []SubscriptionPlan  `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

type SubscriptionPlan struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key"`
    ProductID uuid.UUID `gorm:"type:uuid;not null;index"`
    Product   Product   `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE"`
}
```

**Why This Works:**
- `foreignKey:ProductID` defines which field is the foreign key
- `references:ID` specifies which field in Product it references
- `constraint:OnDelete:CASCADE` ensures subscription plans are deleted when product is deleted
- `index` on ProductID improves query performance
- Maintains referential integrity at the database level

---

### Issue 5: Associations Not Loading

**Problem:**
```go
var product Product
db.First(&product, "name = ?", "Test Product")
fmt.Println(len(product.SubscriptionPlans))  
```

**Symptoms:**
- SubscriptionPlans slice is always empty
- Related data not fetched from database
- Need to make separate queries for associations

**Documentation Consulted:**
- **GORM Preloading**: https://gorm.io/docs/preload.html

**Solution:**
```go
var product Product
db.Preload("SubscriptionPlans").First(&product, "name = ?", "Test Product")
fmt.Println(len(product.SubscriptionPlans)) 
```

**Why This Works:**
- `Preload("SubscriptionPlans")` tells GORM to fetch associated records
- GORM executes an additional query to load the relationship
- Can chain multiple Preloads for nested associations
- Eager loading prevents N+1 query problems

---

### Issue 6: No Error Handling

**Problem:**
```go
db.Create(&product)  
db.First(&product)   
```

**Symptoms:**
- Silent failures
- Difficult to debug issues
- Data corruption possible
- Application continues with invalid state

**Documentation Consulted:**
- **GORM Error Handling**: https://gorm.io/docs/error_handling.html
- **Go Error Handling Best Practices**: https://go.dev/blog/error-handling-and-go

**Solution:**
```go
if err := db.Create(&product).Error; err != nil {
    return fmt.Errorf("failed to create product: %w", err)
}

if err := db.First(&product, id).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return fmt.Errorf("product not found: %w", err)
    }
    return fmt.Errorf("database error: %w", err)
}
```

**Why This Works:**
- `.Error` property contains any error from the operation
- Error wrapping with `%w` preserves error chain
- Can distinguish between "not found" and other database errors
- Provides context for debugging

---

## Complete Fixed Implementation

See the actual implementation in:
- `internal/models/product.go` - Correct Product model
- `internal/models/subscription_plan.go` - Correct SubscriptionPlan model
- `internal/repository/product_repository.go` - Proper error handling and Preload usage

---

## Key Learnings

### 1. Always Use Proper Types
- Use `uuid.UUID` instead of `string` for UUIDs
- Use `time.Time` for timestamps
- Use `gorm.DeletedAt` for soft deletes

### 2. Follow GORM Conventions
- Include `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt` fields
- Use struct tags for constraints (`not null`, `index`, etc.)
- Implement hooks for automatic field population

### 3. Define Relationships Properly
- Use `foreignKey` and `references` tags
- Add `constraint:OnDelete:CASCADE` for referential integrity
- Use `Preload` to fetch associations

### 4. Always Handle Errors
- Check `.Error` on all database operations
- Distinguish between different error types
- Provide context in error messages

### 5. Use Database Constraints
- Foreign keys prevent orphaned records
- Cascade deletes maintain data integrity
- Indexes improve query performance

---

## Testing the Fix

To verify the fixes work correctly:

```bash
# Run the tests
go test ./internal/models/... -v
go test ./internal/repository/... -v

# Check that:
# 1. Products are created with valid UUIDs
# 2. Timestamps are automatically set
# 3. Associations load correctly with Preload
# 4. Cascade delete works (deleting product deletes plans)
# 5. Soft delete works (records marked as deleted, not removed)
```

---

## Documentation References Summary

1. **GORM Models**: https://gorm.io/docs/models.html
2. **GORM Data Types**: https://gorm.io/docs/data_types.html
3. **GORM Hooks**: https://gorm.io/docs/hooks.html
4. **GORM Belongs To**: https://gorm.io/docs/belongs_to.html
5. **GORM Has Many**: https://gorm.io/docs/has_many.html
6. **GORM Preload**: https://gorm.io/docs/preload.html
7. **GORM Constraints**: https://gorm.io/docs/constraints.html
8. **GORM Error Handling**: https://gorm.io/docs/error_handling.html
9. **google/uuid Package**: https://pkg.go.dev/github.com/google/uuid

Each of these documentation sources was essential in understanding and fixing the specific issues in the broken code.

