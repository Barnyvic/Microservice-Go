package service

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/microservice-go/product-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSubscriptionRepository struct {
	mock.Mock
}

func (m *MockSubscriptionRepository) Create(plan *models.SubscriptionPlan) error {
	args := m.Called(plan)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) GetByID(id uuid.UUID) (*models.SubscriptionPlan, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionRepository) Update(plan *models.SubscriptionPlan) error {
	args := m.Called(plan)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) ListByProductID(productID uuid.UUID) ([]models.SubscriptionPlan, error) {
	args := m.Called(productID)
	return args.Get(0).([]models.SubscriptionPlan), args.Error(1)
}

type MockProductRepositoryForSubscription struct {
	mock.Mock
}

func (m *MockProductRepositoryForSubscription) Create(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepositoryForSubscription) GetByID(id uuid.UUID) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepositoryForSubscription) Update(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepositoryForSubscription) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductRepositoryForSubscription) List(productType string, page, pageSize int) ([]models.Product, int64, error) {
	args := m.Called(productType, page, pageSize)
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

func TestCreateSubscriptionPlan_Success(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProductRepo := new(MockProductRepositoryForSubscription)
	service := NewSubscriptionService(mockRepo, mockProductRepo)

	productID := uuid.New()
	expectedProduct := &models.Product{
		ID:          productID,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		ProductType: "digital",
	}

	mockProductRepo.On("GetByID", productID).Return(expectedProduct, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.SubscriptionPlan")).Return(nil)

	plan, err := service.CreateSubscriptionPlan(productID.String(), "Monthly Plan", 30, 29.99)

	assert.NoError(t, err)
	assert.NotNil(t, plan)
	assert.Equal(t, "Monthly Plan", plan.PlanName)
	assert.Equal(t, 30, plan.Duration)
	assert.Equal(t, 29.99, plan.Price)
	assert.Equal(t, productID, plan.ProductID)
	mockRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestCreateSubscriptionPlan_EmptyPlanName(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProductRepo := new(MockProductRepositoryForSubscription)
	service := NewSubscriptionService(mockRepo, mockProductRepo)

	plan, err := service.CreateSubscriptionPlan(uuid.New().String(), "", 30, 29.99)

	assert.Error(t, err)
	assert.Nil(t, plan)
	assert.Contains(t, err.Error(), "plan name is required")
}

func TestCreateSubscriptionPlan_InvalidDuration(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProductRepo := new(MockProductRepositoryForSubscription)
	service := NewSubscriptionService(mockRepo, mockProductRepo)

	plan, err := service.CreateSubscriptionPlan(uuid.New().String(), "Monthly Plan", 0, 29.99)

	assert.Error(t, err)
	assert.Nil(t, plan)
	assert.Contains(t, err.Error(), "duration must be positive")
}

func TestCreateSubscriptionPlan_NegativePrice(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProductRepo := new(MockProductRepositoryForSubscription)
	service := NewSubscriptionService(mockRepo, mockProductRepo)

	plan, err := service.CreateSubscriptionPlan(uuid.New().String(), "Monthly Plan", 30, -10.0)

	assert.Error(t, err)
	assert.Nil(t, plan)
	assert.Contains(t, err.Error(), "price cannot be negative")
}

func TestCreateSubscriptionPlan_InvalidProductID(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProductRepo := new(MockProductRepositoryForSubscription)
	service := NewSubscriptionService(mockRepo, mockProductRepo)

	plan, err := service.CreateSubscriptionPlan("invalid-uuid", "Monthly Plan", 30, 29.99)

	assert.Error(t, err)
	assert.Nil(t, plan)
	assert.Contains(t, err.Error(), "invalid product ID format")
}

func TestCreateSubscriptionPlan_ProductNotFound(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProductRepo := new(MockProductRepositoryForSubscription)
	service := NewSubscriptionService(mockRepo, mockProductRepo)

	productID := uuid.New()
	mockProductRepo.On("GetByID", productID).Return(nil, errors.New("product not found"))

	plan, err := service.CreateSubscriptionPlan(productID.String(), "Monthly Plan", 30, 29.99)

	assert.Error(t, err)
	assert.Nil(t, plan)
	assert.Contains(t, err.Error(), "Product with ID")
	mockProductRepo.AssertExpectations(t)
}

func TestGetSubscriptionPlan_Success(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProductRepo := new(MockProductRepositoryForSubscription)
	service := NewSubscriptionService(mockRepo, mockProductRepo)

	planID := uuid.New()
	expectedPlan := &models.SubscriptionPlan{
		ID:        planID,
		ProductID: uuid.New(),
		PlanName:  "Monthly Plan",
		Duration:  30,
		Price:     29.99,
	}

	mockRepo.On("GetByID", planID).Return(expectedPlan, nil)

	plan, err := service.GetSubscriptionPlan(planID.String())

	assert.NoError(t, err)
	assert.NotNil(t, plan)
	assert.Equal(t, expectedPlan.ID, plan.ID)
	assert.Equal(t, expectedPlan.PlanName, plan.PlanName)
	mockRepo.AssertExpectations(t)
}

func TestGetSubscriptionPlan_InvalidID(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProductRepo := new(MockProductRepositoryForSubscription)
	service := NewSubscriptionService(mockRepo, mockProductRepo)

	plan, err := service.GetSubscriptionPlan("invalid-uuid")

	assert.Error(t, err)
	assert.Nil(t, plan)
	assert.Contains(t, err.Error(), "invalid subscription plan ID format")
}

func TestGetSubscriptionPlan_NotFound(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProductRepo := new(MockProductRepositoryForSubscription)
	service := NewSubscriptionService(mockRepo, mockProductRepo)

	planID := uuid.New()
	mockRepo.On("GetByID", planID).Return(nil, errors.New("subscription plan not found"))

	plan, err := service.GetSubscriptionPlan(planID.String())

	assert.Error(t, err)
	assert.Nil(t, plan)
	assert.Contains(t, err.Error(), "SubscriptionPlan with ID")
	mockRepo.AssertExpectations(t)
}

func TestUpdateSubscriptionPlan_Success(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProductRepo := new(MockProductRepositoryForSubscription)
	service := NewSubscriptionService(mockRepo, mockProductRepo)

	planID := uuid.New()
	productID := uuid.New()
	expectedProduct := &models.Product{
		ID:          productID,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		ProductType: "digital",
	}
	expectedPlan := &models.SubscriptionPlan{
		ID:        planID,
		ProductID: productID,
		PlanName:  "Updated Plan",
		Duration:  60,
		Price:     49.99,
	}

	mockRepo.On("GetByID", planID).Return(&models.SubscriptionPlan{ID: planID}, nil).Once()
	mockProductRepo.On("GetByID", productID).Return(expectedProduct, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.SubscriptionPlan")).Return(nil)
	mockRepo.On("GetByID", planID).Return(expectedPlan, nil).Once()

	plan, err := service.UpdateSubscriptionPlan(planID.String(), productID.String(), "Updated Plan", 60, 49.99)

	assert.NoError(t, err)
	assert.NotNil(t, plan)
	assert.Equal(t, "Updated Plan", plan.PlanName)
	assert.Equal(t, 60, plan.Duration)
	assert.Equal(t, 49.99, plan.Price)
	mockRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestUpdateSubscriptionPlan_PlanNotFound(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProductRepo := new(MockProductRepositoryForSubscription)
	service := NewSubscriptionService(mockRepo, mockProductRepo)

	planID := uuid.New()
	mockRepo.On("GetByID", planID).Return(nil, errors.New("subscription plan not found"))

	plan, err := service.UpdateSubscriptionPlan(planID.String(), uuid.New().String(), "Updated Plan", 60, 49.99)

	assert.Error(t, err)
	assert.Nil(t, plan)
	assert.Contains(t, err.Error(), "SubscriptionPlan with ID")
	mockRepo.AssertExpectations(t)
}

func TestDeleteSubscriptionPlan_Success(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProductRepo := new(MockProductRepositoryForSubscription)
	service := NewSubscriptionService(mockRepo, mockProductRepo)

	planID := uuid.New()
	expectedPlan := &models.SubscriptionPlan{
		ID:        planID,
		ProductID: uuid.New(),
		PlanName:  "Monthly Plan",
		Duration:  30,
		Price:     29.99,
	}

	mockRepo.On("GetByID", planID).Return(expectedPlan, nil)
	mockRepo.On("Delete", planID).Return(nil)

	err := service.DeleteSubscriptionPlan(planID.String())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteSubscriptionPlan_PlanNotFound(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProductRepo := new(MockProductRepositoryForSubscription)
	service := NewSubscriptionService(mockRepo, mockProductRepo)

	planID := uuid.New()
	mockRepo.On("GetByID", planID).Return(nil, errors.New("subscription plan not found"))

	err := service.DeleteSubscriptionPlan(planID.String())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SubscriptionPlan with ID")
	mockRepo.AssertExpectations(t)
}

func TestListSubscriptionPlans_Success(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProductRepo := new(MockProductRepositoryForSubscription)
	service := NewSubscriptionService(mockRepo, mockProductRepo)

	productID := uuid.New()
	expectedPlans := []models.SubscriptionPlan{
		{ID: uuid.New(), ProductID: productID, PlanName: "Monthly Plan", Duration: 30, Price: 29.99},
		{ID: uuid.New(), ProductID: productID, PlanName: "Annual Plan", Duration: 365, Price: 299.99},
	}

	mockRepo.On("ListByProductID", productID).Return(expectedPlans, nil)

	plans, err := service.ListSubscriptionPlans(productID.String())

	assert.NoError(t, err)
	assert.Equal(t, 2, len(plans))
	assert.Equal(t, "Monthly Plan", plans[0].PlanName)
	assert.Equal(t, "Annual Plan", plans[1].PlanName)
	mockRepo.AssertExpectations(t)
}

func TestListSubscriptionPlans_InvalidProductID(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProductRepo := new(MockProductRepositoryForSubscription)
	service := NewSubscriptionService(mockRepo, mockProductRepo)

	plans, err := service.ListSubscriptionPlans("invalid-uuid")

	assert.Error(t, err)
	assert.Nil(t, plans)
	assert.Contains(t, err.Error(), "invalid product ID format")
}

func TestValidateSubscriptionInput_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		planName    string
		duration    int
		price       float64
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid input",
			planName:    "Monthly Plan",
			duration:    30,
			price:       29.99,
			expectError: false,
		},
		{
			name:        "Empty plan name",
			planName:    "",
			duration:    30,
			price:       29.99,
			expectError: true,
			errorMsg:    "plan name is required",
		},
		{
			name:        "Plan name too long",
			planName:    string(make([]byte, 256)),
			duration:    30,
			price:       29.99,
			expectError: true,
			errorMsg:    "plan name must be less than 255 characters",
		},
		{
			name:        "Zero duration",
			planName:    "Monthly Plan",
			duration:    0,
			price:       29.99,
			expectError: true,
			errorMsg:    "duration must be positive",
		},
		{
			name:        "Negative duration",
			planName:    "Monthly Plan",
			duration:    -10,
			price:       29.99,
			expectError: true,
			errorMsg:    "duration must be positive",
		},
		{
			name:        "Duration too long",
			planName:    "Monthly Plan",
			duration:    3651, 
			price:       29.99,
			expectError: true,
			errorMsg:    "duration cannot exceed 3650 days",
		},
		{
			name:        "Negative price",
			planName:    "Monthly Plan",
			duration:    30,
			price:       -10.0,
			expectError: true,
			errorMsg:    "price cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSubscriptionInput(tt.planName, tt.duration, tt.price)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
