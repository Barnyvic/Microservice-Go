package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/microservice-go/product-service/internal/models"
	"gorm.io/gorm"
)

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	Create(product *models.Product) error
	GetByID(id uuid.UUID) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id uuid.UUID) error
	List(productType string, page, pageSize int) ([]models.Product, int64, error)
}

// productRepository implements ProductRepository interface
type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new instance of ProductRepository
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// Create inserts a new product into the database
func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

// GetByID retrieves a product by its ID
func (r *productRepository) GetByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("SubscriptionPlans").First(&product, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

// Update modifies an existing product
func (r *productRepository) Update(product *models.Product) error {
	result := r.db.Model(&models.Product{}).Where("id = ?", product.ID).Updates(product)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}

// Delete removes a product from the database (soft delete)
func (r *productRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.Product{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}

// List retrieves products with optional filtering and pagination
func (r *productRepository) List(productType string, page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{})

	// Apply filter if productType is provided
	if productType != "" {
		query = query.Where("product_type = ?", productType)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	// Fetch products with preloaded subscription plans
	if err := query.Preload("SubscriptionPlans").Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

