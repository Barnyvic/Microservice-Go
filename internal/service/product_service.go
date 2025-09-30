package service

import (
	"github.com/microservice-go/product-service/internal/constants"
	apperrors "github.com/microservice-go/product-service/internal/errors"
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

// CreateProduct creates a new product with validation
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
		return nil, apperrors.NewDatabaseError("create product", err)
	}

	return product, nil
}

// GetProduct retrieves a product by ID
func (s *productService) GetProduct(id string) (*models.Product, error) {
	productID, err := parseProductID(id)
	if err != nil {
		return nil, err
	}

	product, err := s.repo.GetByID(productID)
	if err != nil {
		return nil, apperrors.NewNotFoundError("Product", id)
	}

	return product, nil
}

// UpdateProduct updates an existing product
func (s *productService) UpdateProduct(id, name, description string, price float64, productType string) (*models.Product, error) {
	productID, err := parseProductID(id)
	if err != nil {
		return nil, err
	}

	// Verify product exists
	if _, err := s.repo.GetByID(productID); err != nil {
		return nil, apperrors.NewNotFoundError("Product", id)
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
		return nil, apperrors.NewDatabaseError("update product", err)
	}

	return s.repo.GetByID(productID)
}

// DeleteProduct deletes a product by ID
func (s *productService) DeleteProduct(id string) error {
	productID, err := parseProductID(id)
	if err != nil {
		return err
	}

	// Verify product exists before deletion
	if _, err := s.repo.GetByID(productID); err != nil {
		return apperrors.NewNotFoundError("Product", id)
	}

	if err := s.repo.Delete(productID); err != nil {
		return apperrors.NewDatabaseError("delete product", err)
	}

	return nil
}

// ListProducts retrieves a list of products with optional filtering and pagination
func (s *productService) ListProducts(productType string, page, pageSize int) ([]models.Product, int64, error) {
	page = normalizePage(page)
	pageSize = normalizePageSize(pageSize)

	products, total, err := s.repo.List(productType, page, pageSize)
	if err != nil {
		return nil, 0, apperrors.NewDatabaseError("list products", err)
	}

	return products, total, nil
}

// validateProductInput validates product input fields
func validateProductInput(name string, price float64, productType string) error {
	if name == "" {
		return apperrors.NewValidationError("name", "product name is required")
	}
	if len(name) > 255 {
		return apperrors.NewValidationError("name", "product name must be less than 255 characters")
	}
	if price < 0 {
		return apperrors.NewValidationError("price", "price cannot be negative")
	}
	if productType == "" {
		return apperrors.NewValidationError("productType", "product type is required")
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
