package service

import (
	"github.com/google/uuid"
	apperrors "github.com/microservice-go/product-service/internal/errors"
	"github.com/microservice-go/product-service/internal/models"
	"github.com/microservice-go/product-service/internal/repository"
)

type SubscriptionService interface {
	CreateSubscriptionPlan(productID, planName string, duration int, price float64) (*models.SubscriptionPlan, error)
	GetSubscriptionPlan(id string) (*models.SubscriptionPlan, error)
	UpdateSubscriptionPlan(id, productID, planName string, duration int, price float64) (*models.SubscriptionPlan, error)
	DeleteSubscriptionPlan(id string) error
	ListSubscriptionPlans(productID string) ([]models.SubscriptionPlan, error)
}

type subscriptionService struct {
	repo        repository.SubscriptionRepository
	productRepo repository.ProductRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository, productRepo repository.ProductRepository) SubscriptionService {
	return &subscriptionService{
		repo:        repo,
		productRepo: productRepo,
	}
}

// CreateSubscriptionPlan creates a new subscription plan with validation
func (s *subscriptionService) CreateSubscriptionPlan(productID, planName string, duration int, price float64) (*models.SubscriptionPlan, error) {
	if err := validateSubscriptionInput(planName, duration, price); err != nil {
		return nil, err
	}

	prodID, err := parseProductID(productID)
	if err != nil {
		return nil, err
	}

	if _, err := s.productRepo.GetByID(prodID); err != nil {
		return nil, apperrors.NewNotFoundError("Product", productID)
	}

	plan := &models.SubscriptionPlan{
		ProductID: prodID,
		PlanName:  planName,
		Duration:  duration,
		Price:     price,
	}

	if err := s.repo.Create(plan); err != nil {
		return nil, apperrors.NewDatabaseError("create subscription plan", err)
	}

	return plan, nil
}

func (s *subscriptionService) GetSubscriptionPlan(id string) (*models.SubscriptionPlan, error) {
	planID, err := parsePlanID(id)
	if err != nil {
		return nil, err
	}

	plan, err := s.repo.GetByID(planID)
	if err != nil {
		return nil, apperrors.NewNotFoundError("SubscriptionPlan", id)
	}

	return plan, nil
}

func (s *subscriptionService) UpdateSubscriptionPlan(id, productID, planName string, duration int, price float64) (*models.SubscriptionPlan, error) {
	planID, err := parsePlanID(id)
	if err != nil {
		return nil, err
	}

	if _, err := s.repo.GetByID(planID); err != nil {
		return nil, apperrors.NewNotFoundError("SubscriptionPlan", id)
	}

	if err := validateSubscriptionInput(planName, duration, price); err != nil {
		return nil, err
	}

	prodID, err := parseProductID(productID)
	if err != nil {
		return nil, err
	}

	if _, err := s.productRepo.GetByID(prodID); err != nil {
		return nil, apperrors.NewNotFoundError("Product", productID)
	}

	plan := &models.SubscriptionPlan{
		ID:        planID,
		ProductID: prodID,
		PlanName:  planName,
		Duration:  duration,
		Price:     price,
	}

	if err := s.repo.Update(plan); err != nil {
		return nil, apperrors.NewDatabaseError("update subscription plan", err)
	}

	return s.repo.GetByID(planID)
}

func (s *subscriptionService) DeleteSubscriptionPlan(id string) error {
	planID, err := parsePlanID(id)
	if err != nil {
		return err
	}

	if _, err := s.repo.GetByID(planID); err != nil {
		return apperrors.NewNotFoundError("SubscriptionPlan", id)
	}

	if err := s.repo.Delete(planID); err != nil {
		return apperrors.NewDatabaseError("delete subscription plan", err)
	}

	return nil
}

func (s *subscriptionService) ListSubscriptionPlans(productID string) ([]models.SubscriptionPlan, error) {
	prodID, err := parseProductID(productID)
	if err != nil {
		return nil, err
	}

	plans, err := s.repo.ListByProductID(prodID)
	if err != nil {
		return nil, apperrors.NewDatabaseError("list subscription plans", err)
	}

	return plans, nil
}

func parsePlanID(id string) (uuid.UUID, error) {
	if id == "" {
		return uuid.Nil, apperrors.NewValidationError("id", "subscription plan ID is required")
	}

	planID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, apperrors.NewValidationError("id", "invalid subscription plan ID format")
	}

	return planID, nil
}

func validateSubscriptionInput(planName string, duration int, price float64) error {
	if planName == "" {
		return apperrors.NewValidationError("planName", "plan name is required")
	}
	if len(planName) > 255 {
		return apperrors.NewValidationError("planName", "plan name must be less than 255 characters")
	}
	if duration <= 0 {
		return apperrors.NewValidationError("duration", "duration must be positive")
	}
	if duration > 3650 { 
		return apperrors.NewValidationError("duration", "duration cannot exceed 3650 days")
	}
	if price < 0 {
		return apperrors.NewValidationError("price", "price cannot be negative")
	}
	return nil
}
