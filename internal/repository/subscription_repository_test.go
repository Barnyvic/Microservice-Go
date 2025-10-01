//go:build cgo
// +build cgo

package repository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/microservice-go/product-service/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupSubscriptionTestDB(t *testing.T) *gorm.DB {
	// Use a file-based SQLite database instead of in-memory to avoid CGO issues
	db, err := gorm.Open(sqlite.Open("test_subscription.db"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.Product{}, &models.SubscriptionPlan{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Clean up test data before each test
	db.Exec("DELETE FROM subscription_plans")
	db.Exec("DELETE FROM products")

	return db
}

func TestSubscriptionRepository_Create(t *testing.T) {
	db := setupSubscriptionTestDB(t)
	repo := NewSubscriptionRepository(db)

	// Create a product first
	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		ProductType: "digital",
	}
	err := db.Create(product).Error
	assert.NoError(t, err)

	plan := &models.SubscriptionPlan{
		ProductID: product.ID,
		PlanName:  "Monthly Plan",
		Duration:  30,
		Price:     29.99,
	}

	err = repo.Create(plan)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, plan.ID)
}

func TestSubscriptionRepository_GetByID(t *testing.T) {
	db := setupSubscriptionTestDB(t)
	repo := NewSubscriptionRepository(db)

	// Create a product first
	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		ProductType: "digital",
	}
	err := db.Create(product).Error
	assert.NoError(t, err)

	plan := &models.SubscriptionPlan{
		ProductID: product.ID,
		PlanName:  "Monthly Plan",
		Duration:  30,
		Price:     29.99,
	}
	err = repo.Create(plan)
	assert.NoError(t, err)

	retrieved, err := repo.GetByID(plan.ID)

	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, plan.ID, retrieved.ID)
	assert.Equal(t, plan.PlanName, retrieved.PlanName)
	assert.Equal(t, plan.Duration, retrieved.Duration)
	assert.Equal(t, plan.Price, retrieved.Price)
	assert.Equal(t, plan.ProductID, retrieved.ProductID)
}

func TestSubscriptionRepository_GetByID_NotFound(t *testing.T) {
	db := setupSubscriptionTestDB(t)
	repo := NewSubscriptionRepository(db)

	nonExistentID := uuid.New()
	plan, err := repo.GetByID(nonExistentID)

	assert.Error(t, err)
	assert.Nil(t, plan)
	assert.Equal(t, "subscription plan not found", err.Error())
}

func TestSubscriptionRepository_Update(t *testing.T) {
	db := setupSubscriptionTestDB(t)
	repo := NewSubscriptionRepository(db)

	// Create a product first
	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		ProductType: "digital",
	}
	err := db.Create(product).Error
	assert.NoError(t, err)

	plan := &models.SubscriptionPlan{
		ProductID: product.ID,
		PlanName:  "Original Plan",
		Duration:  30,
		Price:     29.99,
	}
	err = repo.Create(plan)
	assert.NoError(t, err)

	plan.PlanName = "Updated Plan"
	plan.Duration = 60
	plan.Price = 49.99
	err = repo.Update(plan)

	assert.NoError(t, err)

	updated, err := repo.GetByID(plan.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Plan", updated.PlanName)
	assert.Equal(t, 60, updated.Duration)
	assert.Equal(t, 49.99, updated.Price)
}

func TestSubscriptionRepository_Update_NotFound(t *testing.T) {
	db := setupSubscriptionTestDB(t)
	repo := NewSubscriptionRepository(db)

	nonExistentPlan := &models.SubscriptionPlan{
		ID:        uuid.New(),
		ProductID: uuid.New(),
		PlanName:  "Non-existent Plan",
		Duration:  30,
		Price:     29.99,
	}

	err := repo.Update(nonExistentPlan)

	assert.Error(t, err)
	assert.Equal(t, "subscription plan not found", err.Error())
}

func TestSubscriptionRepository_Delete(t *testing.T) {
	db := setupSubscriptionTestDB(t)
	repo := NewSubscriptionRepository(db)

	// Create a product first
	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		ProductType: "digital",
	}
	err := db.Create(product).Error
	assert.NoError(t, err)

	plan := &models.SubscriptionPlan{
		ProductID: product.ID,
		PlanName:  "Monthly Plan",
		Duration:  30,
		Price:     29.99,
	}
	err = repo.Create(plan)
	assert.NoError(t, err)

	err = repo.Delete(plan.ID)
	assert.NoError(t, err)

	deleted, err := repo.GetByID(plan.ID)
	assert.Error(t, err)
	assert.Nil(t, deleted)
}

func TestSubscriptionRepository_Delete_NotFound(t *testing.T) {
	db := setupSubscriptionTestDB(t)
	repo := NewSubscriptionRepository(db)

	nonExistentID := uuid.New()
	err := repo.Delete(nonExistentID)

	assert.Error(t, err)
	assert.Equal(t, "subscription plan not found", err.Error())
}

func TestSubscriptionRepository_ListByProductID(t *testing.T) {
	db := setupSubscriptionTestDB(t)
	repo := NewSubscriptionRepository(db)

	// Create products
	product1 := &models.Product{
		Name:        "Product 1",
		Description: "Test Description",
		Price:       99.99,
		ProductType: "digital",
	}
	product2 := &models.Product{
		Name:        "Product 2",
		Description: "Test Description",
		Price:       199.99,
		ProductType: "physical",
	}
	err := db.Create(product1).Error
	assert.NoError(t, err)
	err = db.Create(product2).Error
	assert.NoError(t, err)

	// Create subscription plans
	plans := []*models.SubscriptionPlan{
		{ProductID: product1.ID, PlanName: "Monthly Plan", Duration: 30, Price: 29.99},
		{ProductID: product1.ID, PlanName: "Annual Plan", Duration: 365, Price: 299.99},
		{ProductID: product2.ID, PlanName: "Quarterly Plan", Duration: 90, Price: 79.99},
	}

	for _, plan := range plans {
		err := repo.Create(plan)
		assert.NoError(t, err)
	}

	// Test listing plans for product1
	product1Plans, err := repo.ListByProductID(product1.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(product1Plans))

	// Test listing plans for product2
	product2Plans, err := repo.ListByProductID(product2.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(product2Plans))

	// Test listing plans for non-existent product
	nonExistentProductID := uuid.New()
	emptyPlans, err := repo.ListByProductID(nonExistentProductID)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(emptyPlans))
}

func TestSubscriptionRepository_CascadeDelete(t *testing.T) {
	db := setupSubscriptionTestDB(t)
	repo := NewSubscriptionRepository(db)

	// Create a product
	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		ProductType: "digital",
	}
	err := db.Create(product).Error
	assert.NoError(t, err)

	// Create subscription plans
	plans := []*models.SubscriptionPlan{
		{ProductID: product.ID, PlanName: "Monthly Plan", Duration: 30, Price: 29.99},
		{ProductID: product.ID, PlanName: "Annual Plan", Duration: 365, Price: 299.99},
	}

	for _, plan := range plans {
		err := repo.Create(plan)
		assert.NoError(t, err)
	}

	// Verify plans exist
	existingPlans, err := repo.ListByProductID(product.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(existingPlans))

	// Delete the product (should cascade delete subscription plans)
	err = db.Delete(product).Error
	assert.NoError(t, err)

	// Verify subscription plans are also deleted
	deletedPlans, err := repo.ListByProductID(product.ID)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(deletedPlans))
}

func TestSubscriptionRepository_SoftDelete(t *testing.T) {
	db := setupSubscriptionTestDB(t)
	repo := NewSubscriptionRepository(db)

	// Create a product first
	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		ProductType: "digital",
	}
	err := db.Create(product).Error
	assert.NoError(t, err)

	plan := &models.SubscriptionPlan{
		ProductID: product.ID,
		PlanName:  "Monthly Plan",
		Duration:  30,
		Price:     29.99,
	}
	err = repo.Create(plan)
	assert.NoError(t, err)

	// Verify plan exists
	retrieved, err := repo.GetByID(plan.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)

	// Soft delete the plan
	err = repo.Delete(plan.ID)
	assert.NoError(t, err)

	// Verify plan is soft deleted (not found via GetByID)
	deleted, err := repo.GetByID(plan.ID)
	assert.Error(t, err)
	assert.Nil(t, deleted)

	// Verify plan still exists in database but with DeletedAt set
	var count int64
	err = db.Unscoped().Model(&models.SubscriptionPlan{}).Where("id = ?", plan.ID).Count(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}
