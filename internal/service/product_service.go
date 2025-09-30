package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/microservice-go/product-service/internal/models"
	"github.com/microservice-go/product-service/internal/repository"
)

// ProductService defines the interface for product business logic
type ProductService interface {
	CreateProduct(name, description string, price float64, productType string) (*models.Product, error)
	GetProduct(id string) (*models.Product, error)
	UpdateProduct(id, name, description string, price float64, productType string) (*models.Product, error)
	DeleteProduct(id string) error
	ListProducts(productType string, page, pageSize int) ([]models.Product, int64, error)
}

// productService implements ProductService interface
type productService struct {
	repo repository.ProductRepository
}

// NewProductService creates a new instance of ProductService
func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

// CreateProduct creates a new product
func (s *productService) CreateProduct(name, description string, price float64, productType string) (*models.Product, error) {
	if name == "" {
		return nil, errors.New("product name is required")
	}
	if price < 0 {
		return nil, errors.New("price cannot be negative")
	}
	if productType == "" {
		return nil, errors.New("product type is required")
	}

	product := &models.Product{
		Name:        name,
		Description: description,
		Price:       price,
		ProductType: productType,
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

// GetProduct retrieves a product by ID
func (s *productService) GetProduct(id string) (*models.Product, error) {
	productID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid product ID format")
	}

	return s.repo.GetByID(productID)
}

// UpdateProduct updates an existing product
func (s *productService) UpdateProduct(id, name, description string, price float64, productType string) (*models.Product, error) {
	productID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid product ID format")
	}

	if name == "" {
		return nil, errors.New("product name is required")
	}
	if price < 0 {
		return nil, errors.New("price cannot be negative")
	}
	if productType == "" {
		return nil, errors.New("product type is required")
	}

	product := &models.Product{
		ID:          productID,
		Name:        name,
		Description: description,
		Price:       price,
		ProductType: productType,
	}

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}

	return s.repo.GetByID(productID)
}

// DeleteProduct deletes a product by ID
func (s *productService) DeleteProduct(id string) error {
	productID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid product ID format")
	}

	return s.repo.Delete(productID)
}

// ListProducts retrieves a list of products with optional filtering
func (s *productService) ListProducts(productType string, page, pageSize int) ([]models.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return s.repo.List(productType, page, pageSize)
}

