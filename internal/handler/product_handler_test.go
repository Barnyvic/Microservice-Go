package handler

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/microservice-go/product-service/internal/models"
	pb "github.com/microservice-go/product-service/proto/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductService is a mock implementation of ProductService
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(name, description string, price float64, productType string) (*models.Product, error) {
	args := m.Called(name, description, price, productType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) GetProduct(id string) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) UpdateProduct(id, name, description string, price float64, productType string) (*models.Product, error) {
	args := m.Called(id, name, description, price, productType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) DeleteProduct(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductService) ListProducts(productType string, page, pageSize int) ([]models.Product, int64, error) {
	args := m.Called(productType, page, pageSize)
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

func TestProductHandler_CreateProduct(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	productID := uuid.New()
	expectedProduct := &models.Product{
		ID:          productID,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		ProductType: "digital",
	}

	mockService.On("CreateProduct", "Test Product", "Test Description", 99.99, "digital").
		Return(expectedProduct, nil)

	req := &pb.CreateProductRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		ProductType: "digital",
	}

	resp, err := handler.CreateProduct(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, productID.String(), resp.Product.Id)
	assert.Equal(t, "Test Product", resp.Product.Name)
	assert.Equal(t, 99.99, resp.Product.Price)
	mockService.AssertExpectations(t)
}

func TestProductHandler_GetProduct(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	productID := uuid.New()
	expectedProduct := &models.Product{
		ID:          productID,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		ProductType: "digital",
	}

	mockService.On("GetProduct", productID.String()).Return(expectedProduct, nil)

	req := &pb.GetProductRequest{
		Id: productID.String(),
	}

	resp, err := handler.GetProduct(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, productID.String(), resp.Product.Id)
	assert.Equal(t, "Test Product", resp.Product.Name)
	mockService.AssertExpectations(t)
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	productID := uuid.New()
	mockService.On("DeleteProduct", productID.String()).Return(nil)

	req := &pb.DeleteProductRequest{
		Id: productID.String(),
	}

	resp, err := handler.DeleteProduct(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, "Product deleted successfully", resp.Message)
	mockService.AssertExpectations(t)
}

func TestProductHandler_ListProducts(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	products := []models.Product{
		{ID: uuid.New(), Name: "Product 1", Price: 10.0, ProductType: "digital"},
		{ID: uuid.New(), Name: "Product 2", Price: 20.0, ProductType: "digital"},
	}

	mockService.On("ListProducts", "digital", 1, 10).Return(products, int64(2), nil)

	req := &pb.ListProductsRequest{
		ProductType: "digital",
		Page:        1,
		PageSize:    10,
	}

	resp, err := handler.ListProducts(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Products))
	assert.Equal(t, int32(2), resp.Total)
	mockService.AssertExpectations(t)
}

