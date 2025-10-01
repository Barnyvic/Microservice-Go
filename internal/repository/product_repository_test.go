package repository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/microservice-go/product-service/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.Product{}, &models.SubscriptionPlan{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestProductRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProductRepository(db)

	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		ProductType: "digital",
	}

	err := repo.Create(product)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, product.ID)
}

func TestProductRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProductRepository(db)

	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		ProductType: "digital",
	}
	err := repo.Create(product)
	assert.NoError(t, err)

	retrieved, err := repo.GetByID(product.ID)

	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, product.ID, retrieved.ID)
	assert.Equal(t, product.Name, retrieved.Name)
	assert.Equal(t, product.Price, retrieved.Price)
}

func TestProductRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProductRepository(db)

	nonExistentID := uuid.New()
	product, err := repo.GetByID(nonExistentID)

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Equal(t, "product not found", err.Error())
}

func TestProductRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProductRepository(db)

	product := &models.Product{
		Name:        "Original Name",
		Description: "Original Description",
		Price:       99.99,
		ProductType: "digital",
	}
	err := repo.Create(product)
	assert.NoError(t, err)

	product.Name = "Updated Name"
	product.Price = 149.99
	err = repo.Update(product)

	assert.NoError(t, err)

	updated, err := repo.GetByID(product.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, 149.99, updated.Price)
}

func TestProductRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProductRepository(db)

	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		ProductType: "digital",
	}
	err := repo.Create(product)
	assert.NoError(t, err)

	err = repo.Delete(product.ID)
	assert.NoError(t, err)

	deleted, err := repo.GetByID(product.ID)
	assert.Error(t, err)
	assert.Nil(t, deleted)
}

func TestProductRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProductRepository(db)

	products := []*models.Product{
		{Name: "Product 1", Price: 10.0, ProductType: "digital"},
		{Name: "Product 2", Price: 20.0, ProductType: "physical"},
		{Name: "Product 3", Price: 30.0, ProductType: "digital"},
	}

	for _, p := range products {
		err := repo.Create(p)
		assert.NoError(t, err)
	}

	allProducts, total, err := repo.List("", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(allProducts))
	assert.Equal(t, int64(3), total)

	digitalProducts, total, err := repo.List("digital", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(digitalProducts))
	assert.Equal(t, int64(2), total)
}

func TestProductRepository_List_Pagination(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProductRepository(db)

	for i := 0; i < 5; i++ {
		product := &models.Product{
			Name:        "Product",
			Price:       10.0,
			ProductType: "digital",
		}
		err := repo.Create(product)
		assert.NoError(t, err)
	}

	products, total, err := repo.List("", 1, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(products))
	assert.Equal(t, int64(5), total)

	products, total, err = repo.List("", 2, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(products))
	assert.Equal(t, int64(5), total)
}

