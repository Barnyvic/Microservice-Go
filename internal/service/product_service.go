package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/microservice-go/product-service/internal/constants"
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
	if err := validateProductInput(name, price, productType); err != nil {
		return nil, err
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
		return nil, errors.New(constants.ErrInvalidProductID)
	}

	return s.repo.GetByID(productID)
}

// UpdateProduct updates an existing product
func (s *productService) UpdateProduct(id, name, description string, price float64, productType string) (*models.Product, error) {
	productID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New(constants.ErrInvalidProductID)
	}

	if err := validateProductInput(name, price, productType); err != nil {
		return nil, err
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
		return errors.New(constants.ErrInvalidProductID)
	}

	return s.repo.Delete(productID)
}

// ListProducts retrieves a list of products with optional filtering
func (s *productService) ListProducts(productType string, page, pageSize int) ([]models.Product, int64, error) {
	page = normalizePage(page)
	pageSize = normalizePageSize(pageSize)

	return s.repo.List(productType, page, pageSize)
}

// validateProductInput validates product input fields
func validateProductInput(name string, price float64, productType string) error {
	if name == "" {
		return errors.New(constants.ErrProductNameRequired)
	}
	if price < 0 {
		return errors.New(constants.ErrPriceNegative)
	}
	if productType == "" {
		return errors.New(constants.ErrProductTypeRequired)
	}
	return nil
}

// normalizePage ensures page number is within valid range
func normalizePage(page int) int {
	if page < constants.MinPageSize {
		return constants.DefaultPage
	}
	return page
}

// normalizePageSize ensures page size is within valid range
func normalizePageSize(pageSize int) int {
	if pageSize < constants.MinPageSize {
		return constants.DefaultPageSize
	}
	if pageSize > constants.MaxPageSize {
		return constants.MaxPageSize
	}
	return pageSize
}
