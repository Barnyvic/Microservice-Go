package service

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/microservice-go/product-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository is a mock implementation of ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) GetByID(id uuid.UUID) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) Update(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductRepository) List(productType string, page, pageSize int) ([]models.Product, int64, error) {
	args := m.Called(productType, page, pageSize)
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

func TestCreateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	mockRepo.On("Create", mock.AnythingOfType("*models.Product")).Return(nil)

	product, err := service.CreateProduct("Test Product", "Test Description", 99.99, "digital")

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "Test Product", product.Name)
	assert.Equal(t, "Test Description", product.Description)
	assert.Equal(t, 99.99, product.Price)
	assert.Equal(t, "digital", product.ProductType)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_EmptyName(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	product, err := service.CreateProduct("", "Test Description", 99.99, "digital")

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Contains(t, err.Error(), "product name is required")
}

func TestCreateProduct_NegativePrice(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	product, err := service.CreateProduct("Test Product", "Test Description", -10.0, "digital")

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Contains(t, err.Error(), "price cannot be negative")
}

func TestGetProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	productID := uuid.New()
	expectedProduct := &models.Product{
		ID:          productID,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		ProductType: "digital",
	}

	mockRepo.On("GetByID", productID).Return(expectedProduct, nil)

	product, err := service.GetProduct(productID.String())

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, expectedProduct.ID, product.ID)
	assert.Equal(t, expectedProduct.Name, product.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetProduct_InvalidID(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	product, err := service.GetProduct("invalid-uuid")

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Contains(t, err.Error(), "invalid product ID format")
}

func TestGetProduct_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	productID := uuid.New()
	mockRepo.On("GetByID", productID).Return(nil, errors.New("product not found"))

	product, err := service.GetProduct(productID.String())

	assert.Error(t, err)
	assert.Nil(t, product)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	productID := uuid.New()
	expectedProduct := &models.Product{
		ID:          productID,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		ProductType: "digital",
	}

	// Mock the GetByID call that happens before delete
	mockRepo.On("GetByID", productID).Return(expectedProduct, nil)
	mockRepo.On("Delete", productID).Return(nil)

	err := service.DeleteProduct(productID.String())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestListProducts_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	expectedProducts := []models.Product{
		{ID: uuid.New(), Name: "Product 1", Price: 10.0, ProductType: "digital"},
		{ID: uuid.New(), Name: "Product 2", Price: 20.0, ProductType: "physical"},
	}

	mockRepo.On("List", "digital", 1, 10).Return(expectedProducts, int64(2), nil)

	products, total, err := service.ListProducts("digital", 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(products))
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

