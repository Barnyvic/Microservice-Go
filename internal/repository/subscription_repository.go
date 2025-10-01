package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/microservice-go/product-service/internal/models"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(plan *models.SubscriptionPlan) error
	GetByID(id uuid.UUID) (*models.SubscriptionPlan, error)
	Update(plan *models.SubscriptionPlan) error
	Delete(id uuid.UUID) error
	ListByProductID(productID uuid.UUID) ([]models.SubscriptionPlan, error)
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(plan *models.SubscriptionPlan) error {
	return r.db.Create(plan).Error
}

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

func (r *subscriptionRepository) ListByProductID(productID uuid.UUID) ([]models.SubscriptionPlan, error) {
	var plans []models.SubscriptionPlan
	err := r.db.Where("product_id = ?", productID).Find(&plans).Error
	if err != nil {
		return nil, err
	}
	return plans, nil
}
