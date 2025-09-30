package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Product model - BROKEN VERSION
type Product struct {
	ID          string // Wrong type - should be uuid.UUID
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// SubscriptionPlan model - BROKEN VERSION
type SubscriptionPlan struct {
	ID        string // Wrong type - should be uuid.UUID
	ProductID string // Wrong type - should be uuid.UUID
	PlanName  string
	Duration  int
	Price     float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	// Connect to database
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Product{}, &SubscriptionPlan{})

	// Create a product
	product := Product{
		ID:          uuid.New().String(), // Manual UUID generation
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
	}

	// This will fail because GORM expects proper UUID type for foreign key relationships
	result := db.Create(&product)
	if result.Error != nil {
		fmt.Printf("Error creating product: %v\n", result.Error)
	}

	// Create a subscription plan
	plan := SubscriptionPlan{
		ID:        uuid.New().String(),
		ProductID: product.ID, // This relationship won't work properly
		PlanName:  "Monthly Plan",
		Duration:  30,
		Price:     29.99,
	}

	result = db.Create(&plan)
	if result.Error != nil {
		fmt.Printf("Error creating subscription plan: %v\n", result.Error)
	}

	// Try to query with preload - this will fail
	var retrievedProduct Product
	result = db.Preload("SubscriptionPlans").First(&retrievedProduct, "id = ?", product.ID)
	if result.Error != nil {
		fmt.Printf("Error retrieving product: %v\n", result.Error)
	}

	fmt.Println("Operations completed")
}

