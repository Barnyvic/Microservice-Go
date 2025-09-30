package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/microservice-go/product-service/internal/models"
	"gorm.io/gorm"
)

// SubscriptionRepository defines the interface for subscription plan data access
type SubscriptionRepository interface {
	Create(plan *models.SubscriptionPlan) error
	GetByID(id uuid.UUID) (*models.SubscriptionPlan, error)
	Update(plan *models.SubscriptionPlan) error
	Delete(id uuid.UUID) error
	ListByProductID(productID uuid.UUID) ([]models.SubscriptionPlan, error)
}

// subscriptionRepository implements SubscriptionRepository interface
type subscriptionRepository struct {
	db *gorm.DB
}

// NewSubscriptionRepository creates a new instance of SubscriptionRepository
func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

// Create inserts a new subscription plan into the database
func (r *subscriptionRepository) Create(plan *models.SubscriptionPlan) error {
	return r.db.Create(plan).Error
}

// GetByID retrieves a subscription plan by its ID
func (r *subscriptionRepository) GetByID(id uuid.UUID) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan
	err := r.db.Preload("Product").First(&plan, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subscription plan not found")
		}
		return nil, err
	}
	return &plan, nil
}

// Update modifies an existing subscription plan
func (r *subscriptionRepository) Update(plan *models.SubscriptionPlan) error {
	result := r.db.Model(&models.SubscriptionPlan{}).Where("id = ?", plan.ID).Updates(plan)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("subscription plan not found")
	}
	return nil
}

// Delete removes a subscription plan from the database (soft delete)
func (r *subscriptionRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.SubscriptionPlan{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("subscription plan not found")
	}
	return nil
}

// ListByProductID retrieves all subscription plans for a specific product
func (r *subscriptionRepository) ListByProductID(productID uuid.UUID) ([]models.SubscriptionPlan, error) {
	var plans []models.SubscriptionPlan
	err := r.db.Where("product_id = ?", productID).Find(&plans).Error
	if err != nil {
		return nil, err
	}
	return plans, nil
}

