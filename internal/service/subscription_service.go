package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/microservice-go/product-service/internal/constants"
	"github.com/microservice-go/product-service/internal/models"
	"github.com/microservice-go/product-service/internal/repository"
)

// SubscriptionService defines the interface for subscription plan business logic
type SubscriptionService interface {
	CreateSubscriptionPlan(productID, planName string, duration int, price float64) (*models.SubscriptionPlan, error)
	GetSubscriptionPlan(id string) (*models.SubscriptionPlan, error)
	UpdateSubscriptionPlan(id, productID, planName string, duration int, price float64) (*models.SubscriptionPlan, error)
	DeleteSubscriptionPlan(id string) error
	ListSubscriptionPlans(productID string) ([]models.SubscriptionPlan, error)
}

// subscriptionService implements SubscriptionService interface
type subscriptionService struct {
	repo        repository.SubscriptionRepository
	productRepo repository.ProductRepository
}

// NewSubscriptionService creates a new instance of SubscriptionService
func NewSubscriptionService(repo repository.SubscriptionRepository, productRepo repository.ProductRepository) SubscriptionService {
	return &subscriptionService{
		repo:        repo,
		productRepo: productRepo,
	}
}

// CreateSubscriptionPlan creates a new subscription plan
func (s *subscriptionService) CreateSubscriptionPlan(productID, planName string, duration int, price float64) (*models.SubscriptionPlan, error) {
	if err := validateSubscriptionInput(planName, duration, price); err != nil {
		return nil, err
	}

	prodID, err := uuid.Parse(productID)
	if err != nil {
		return nil, errors.New(constants.ErrInvalidProductID)
	}

	// Verify product exists
	_, err = s.productRepo.GetByID(prodID)
	if err != nil {
		return nil, errors.New(constants.ErrProductNotFound)
	}

	plan := &models.SubscriptionPlan{
		ProductID: prodID,
		PlanName:  planName,
		Duration:  duration,
		Price:     price,
	}

	if err := s.repo.Create(plan); err != nil {
		return nil, err
	}

	return plan, nil
}

// GetSubscriptionPlan retrieves a subscription plan by ID
func (s *subscriptionService) GetSubscriptionPlan(id string) (*models.SubscriptionPlan, error) {
	planID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New(constants.ErrInvalidPlanID)
	}

	return s.repo.GetByID(planID)
}

// UpdateSubscriptionPlan updates an existing subscription plan
func (s *subscriptionService) UpdateSubscriptionPlan(id, productID, planName string, duration int, price float64) (*models.SubscriptionPlan, error) {
	planID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New(constants.ErrInvalidPlanID)
	}

	if err := validateSubscriptionInput(planName, duration, price); err != nil {
		return nil, err
	}

	prodID, err := uuid.Parse(productID)
	if err != nil {
		return nil, errors.New(constants.ErrInvalidProductID)
	}

	// Verify product exists
	_, err = s.productRepo.GetByID(prodID)
	if err != nil {
		return nil, errors.New(constants.ErrProductNotFound)
	}

	plan := &models.SubscriptionPlan{
		ID:        planID,
		ProductID: prodID,
		PlanName:  planName,
		Duration:  duration,
		Price:     price,
	}

	if err := s.repo.Update(plan); err != nil {
		return nil, err
	}

	return s.repo.GetByID(planID)
}

// DeleteSubscriptionPlan deletes a subscription plan by ID
func (s *subscriptionService) DeleteSubscriptionPlan(id string) error {
	planID, err := uuid.Parse(id)
	if err != nil {
		return errors.New(constants.ErrInvalidPlanID)
	}

	return s.repo.Delete(planID)
}

// ListSubscriptionPlans retrieves all subscription plans for a product
func (s *subscriptionService) ListSubscriptionPlans(productID string) ([]models.SubscriptionPlan, error) {
	prodID, err := uuid.Parse(productID)
	if err != nil {
		return nil, errors.New(constants.ErrInvalidProductID)
	}

	return s.repo.ListByProductID(prodID)
}

// validateSubscriptionInput validates subscription plan input fields
func validateSubscriptionInput(planName string, duration int, price float64) error {
	if planName == "" {
		return errors.New(constants.ErrPlanNameRequired)
	}
	if duration <= 0 {
		return errors.New(constants.ErrDurationPositive)
	}
	if price < 0 {
		return errors.New(constants.ErrPriceNegative)
	}
	return nil
}
