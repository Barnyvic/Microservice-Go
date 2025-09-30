# Debug Exercise: UUID and GORM Associations

## The Problem

The broken code in `broken_code.go` has several critical issues related to UUID handling and GORM associations:

### 1. **Incorrect UUID Type**
```go
// BROKEN
type Product struct {
    ID string  // Wrong!
    ...
}
```

**Error Encountered:**
- Foreign key relationships don't work properly
- GORM's `Preload` fails to establish associations
- Database queries using UUID comparisons are inefficient
- No automatic UUID generation on record creation

### 2. **Missing GORM Relationship Tags**
```go
// BROKEN
type SubscriptionPlan struct {
    ProductID string  // No relationship defined
    ...
}
```

**Error Encountered:**
- GORM doesn't recognize the relationship between Product and SubscriptionPlan
- `Preload("SubscriptionPlans")` fails because GORM doesn't know about this association
- No foreign key constraints in the database

### 3. **Manual UUID Generation**
```go
// BROKEN
product := Product{
    ID: uuid.New().String(),  // Manual and error-prone
    ...
}
```

**Error Encountered:**
- Easy to forget to set ID before creating records
- Inconsistent UUID generation across the codebase
- No validation that UUID is set

## Documentation Used

### 1. **GORM Data Types Documentation**
**URL:** https://gorm.io/docs/data_types.html

**What I Learned:**
- GORM supports `uuid.UUID` type natively when using appropriate database drivers
- Need to specify `gorm:"type:uuid"` tag for PostgreSQL
- UUID fields should use `uuid.UUID` type, not `string`

**How It Helped:**
Changed from:
```go
ID string
```
To:
```go
ID uuid.UUID `gorm:"type:uuid;primary_key"`
```

### 2. **GORM Associations (Belongs To / Has Many)**
**URL:** https://gorm.io/docs/belongs_to.html and https://gorm.io/docs/has_many.html

**What I Learned:**
- GORM requires explicit relationship definitions using struct tags
- Foreign keys need to be properly typed (same type as referenced primary key)
- Use `foreignKey` and `references` tags to define relationships
- Can use `constraint:OnDelete:CASCADE` for cascading deletes

**How It Helped:**
Added proper relationship definitions:
```go
type Product struct {
    ...
    SubscriptionPlans []SubscriptionPlan `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

type SubscriptionPlan struct {
    ProductID uuid.UUID `gorm:"type:uuid;not null;index"`
    Product   Product   `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE"`
}
```

### 3. **GORM Hooks (BeforeCreate)**
**URL:** https://gorm.io/docs/hooks.html

**What I Learned:**
- GORM provides lifecycle hooks like `BeforeCreate`, `BeforeUpdate`, etc.
- Hooks can be used to automatically generate UUIDs before creating records
- Hooks ensure consistency across all record creations

**How It Helped:**
Implemented automatic UUID generation:
```go
func (p *Product) BeforeCreate(tx *gorm.DB) error {
    if p.ID == uuid.Nil {
        p.ID = uuid.New()
    }
    return nil
}
```

### 4. **GORM Preload (Eager Loading)**
**URL:** https://gorm.io/docs/preload.html

**What I Learned:**
- `Preload` requires properly defined associations
- Association names in `Preload()` must match struct field names
- Can preload nested associations

**How It Helped:**
Understood why preloading was failing and how to fix it by properly defining associations.

## The Solution

### Fixed Code Structure:

```go
type Product struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key"`
    Name        string    `gorm:"not null"`
    Description string    `gorm:"type:text"`
    Price       float64   `gorm:"not null"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
    
    // Properly defined relationship
    SubscriptionPlans []SubscriptionPlan `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
    if p.ID == uuid.Nil {
        p.ID = uuid.New()
    }
    return nil
}

type SubscriptionPlan struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key"`
    ProductID uuid.UUID `gorm:"type:uuid;not null;index"`
    PlanName  string    `gorm:"not null"`
    Duration  int       `gorm:"not null"`
    Price     float64   `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
    
    // Properly defined relationship
    Product Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE"`
}

func (s *SubscriptionPlan) BeforeCreate(tx *gorm.DB) error {
    if s.ID == uuid.Nil {
        s.ID = uuid.New()
    }
    return nil
}
```

## Key Takeaways

1. **Always use proper types**: Use `uuid.UUID` instead of `string` for UUID fields
2. **Define relationships explicitly**: GORM needs struct tags to understand associations
3. **Use hooks for automation**: BeforeCreate hooks ensure UUIDs are always generated
4. **Read the documentation**: GORM's documentation is comprehensive and provides clear examples
5. **Test relationships**: Always test that preloading and associations work as expected

## Testing the Fix

After implementing the fixes, the following operations work correctly:

```go
// Create product (UUID auto-generated)
product := &Product{
    Name:  "Test Product",
    Price: 99.99,
}
db.Create(product)

// Create subscription plan with proper foreign key
plan := &SubscriptionPlan{
    ProductID: product.ID,  // Proper UUID type
    PlanName:  "Monthly",
    Duration:  30,
    Price:     29.99,
}
db.Create(plan)

// Preload works correctly
var retrieved Product
db.Preload("SubscriptionPlans").First(&retrieved, "id = ?", product.ID)
// retrieved.SubscriptionPlans now contains the associated plans
```

