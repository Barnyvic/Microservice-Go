package main

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


type Product struct {
	ID          string 
	Name        string  `gorm:"not null"`
	Description string
	Price       float64 `gorm:"not null"`
	ProductType string  `gorm:"not null"`
	
	SubscriptionPlans []SubscriptionPlan
}

type SubscriptionPlan struct {
	ID        string 
	ProductID string
	PlanName  string  `gorm:"not null"`
	Duration  int     `gorm:"not null"`
	Price     float64 `gorm:"not null"`
	
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(&Product{}, &SubscriptionPlan{})

	product := Product{
		Name:        "Test Product",
		Description: "A test product",
		Price:       99.99,
		ProductType: "digital",
	}

	db.Create(&product)
	
	fmt.Printf("Created product with ID: %s\n", product.ID) 
	var fetchedProduct Product
	db.First(&fetchedProduct, "name = ?", "Test Product")
	
	fmt.Printf("Product has %d subscription plans\n", len(fetchedProduct.SubscriptionPlans))

	plan := SubscriptionPlan{
		ProductID: "invalid-uuid",
		PlanName:  "Monthly Plan",
		Duration:  30,
		Price:     29.99,
	}
	
	db.Create(&plan)
	
	db.Delete(&product)
}

/*
EXPECTED ERRORS/ISSUES:

1. Product ID will be empty string instead of a valid UUID
2. No automatic timestamp tracking (CreatedAt, UpdatedAt)
3. No soft delete support
4. Foreign key relationship not properly defined
5. No cascade delete - orphaned records
6. Associations not loaded (SubscriptionPlans will be empty)
7. No UUID validation
8. No BeforeCreate hook for UUID generation
9. String-based IDs instead of proper UUID type
10. No error handling for database operations

SYMPTOMS:
- Products created with empty IDs
- Cannot query by ID reliably
- Associations don't load automatically
- Deleting products leaves orphaned subscription plans
- No audit trail (missing timestamps)
- Cannot use soft deletes
*/

