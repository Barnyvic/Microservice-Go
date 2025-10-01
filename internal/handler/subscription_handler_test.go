package handler

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/microservice-go/product-service/internal/models"
	pb "github.com/microservice-go/product-service/proto/subscription"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSubscriptionService struct {
	mock.Mock
}

func (m *MockSubscriptionService) CreateSubscriptionPlan(productID, planName string, duration int, price float64) (*models.SubscriptionPlan, error) {
	args := m.Called(productID, planName, duration, price)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionService) GetSubscriptionPlan(id string) (*models.SubscriptionPlan, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionService) UpdateSubscriptionPlan(id, productID, planName string, duration int, price float64) (*models.SubscriptionPlan, error) {
	args := m.Called(id, productID, planName, duration, price)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionService) DeleteSubscriptionPlan(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSubscriptionService) ListSubscriptionPlans(productID string) ([]models.SubscriptionPlan, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.SubscriptionPlan), args.Error(1)
}

func TestSubscriptionHandler_CreateSubscriptionPlan(t *testing.T) {
	mockService := new(MockSubscriptionService)
	handler := NewSubscriptionHandler(mockService)

	planID := uuid.New()
	productID := uuid.New()
	expectedPlan := &models.SubscriptionPlan{
		ID:        planID,
		ProductID: productID,
		PlanName:  "Monthly Plan",
		Duration:  30,
		Price:     29.99,
	}

	mockService.On("CreateSubscriptionPlan", productID.String(), "Monthly Plan", 30, 29.99).
		Return(expectedPlan, nil)

	req := &pb.CreateSubscriptionPlanRequest{
		ProductId: productID.String(),
		PlanName:  "Monthly Plan",
		Duration:  30,
		Price:     29.99,
	}

	resp, err := handler.CreateSubscriptionPlan(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, planID.String(), resp.Plan.Id)
	assert.Equal(t, productID.String(), resp.Plan.ProductId)
	assert.Equal(t, "Monthly Plan", resp.Plan.PlanName)
	assert.Equal(t, int32(30), resp.Plan.Duration)
	assert.Equal(t, 29.99, resp.Plan.Price)
	mockService.AssertExpectations(t)
}

func TestSubscriptionHandler_GetSubscriptionPlan(t *testing.T) {
	mockService := new(MockSubscriptionService)
	handler := NewSubscriptionHandler(mockService)

	planID := uuid.New()
	productID := uuid.New()
	expectedPlan := &models.SubscriptionPlan{
		ID:        planID,
		ProductID: productID,
		PlanName:  "Monthly Plan",
		Duration:  30,
		Price:     29.99,
	}

	mockService.On("GetSubscriptionPlan", planID.String()).Return(expectedPlan, nil)

	req := &pb.GetSubscriptionPlanRequest{
		Id: planID.String(),
	}

	resp, err := handler.GetSubscriptionPlan(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, planID.String(), resp.Plan.Id)
	assert.Equal(t, productID.String(), resp.Plan.ProductId)
	assert.Equal(t, "Monthly Plan", resp.Plan.PlanName)
	assert.Equal(t, int32(30), resp.Plan.Duration)
	assert.Equal(t, 29.99, resp.Plan.Price)
	mockService.AssertExpectations(t)
}

func TestSubscriptionHandler_UpdateSubscriptionPlan(t *testing.T) {
	mockService := new(MockSubscriptionService)
	handler := NewSubscriptionHandler(mockService)

	planID := uuid.New()
	productID := uuid.New()
	expectedPlan := &models.SubscriptionPlan{
		ID:        planID,
		ProductID: productID,
		PlanName:  "Updated Plan",
		Duration:  60,
		Price:     49.99,
	}

	mockService.On("UpdateSubscriptionPlan", planID.String(), productID.String(), "Updated Plan", 60, 49.99).
		Return(expectedPlan, nil)

	req := &pb.UpdateSubscriptionPlanRequest{
		Id:        planID.String(),
		ProductId: productID.String(),
		PlanName:  "Updated Plan",
		Duration:  60,
		Price:     49.99,
	}

	resp, err := handler.UpdateSubscriptionPlan(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, planID.String(), resp.Plan.Id)
	assert.Equal(t, productID.String(), resp.Plan.ProductId)
	assert.Equal(t, "Updated Plan", resp.Plan.PlanName)
	assert.Equal(t, int32(60), resp.Plan.Duration)
	assert.Equal(t, 49.99, resp.Plan.Price)
	mockService.AssertExpectations(t)
}

func TestSubscriptionHandler_DeleteSubscriptionPlan(t *testing.T) {
	mockService := new(MockSubscriptionService)
	handler := NewSubscriptionHandler(mockService)

	planID := uuid.New()
	mockService.On("DeleteSubscriptionPlan", planID.String()).Return(nil)

	req := &pb.DeleteSubscriptionPlanRequest{
		Id: planID.String(),
	}

	resp, err := handler.DeleteSubscriptionPlan(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, "Subscription plan deleted successfully", resp.Message)
	mockService.AssertExpectations(t)
}

func TestSubscriptionHandler_ListSubscriptionPlans(t *testing.T) {
	mockService := new(MockSubscriptionService)
	handler := NewSubscriptionHandler(mockService)

	productID := uuid.New()
	plans := []models.SubscriptionPlan{
		{ID: uuid.New(), ProductID: productID, PlanName: "Monthly Plan", Duration: 30, Price: 29.99},
		{ID: uuid.New(), ProductID: productID, PlanName: "Annual Plan", Duration: 365, Price: 299.99},
	}

	mockService.On("ListSubscriptionPlans", productID.String()).Return(plans, nil)

	req := &pb.ListSubscriptionPlansRequest{
		ProductId: productID.String(),
	}

	resp, err := handler.ListSubscriptionPlans(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Plans))
	assert.Equal(t, int32(2), resp.Total)
	assert.Equal(t, "Monthly Plan", resp.Plans[0].PlanName)
	assert.Equal(t, "Annual Plan", resp.Plans[1].PlanName)
	mockService.AssertExpectations(t)
}

func TestSubscriptionHandler_CreateSubscriptionPlan_ServiceError(t *testing.T) {
	mockService := new(MockSubscriptionService)
	handler := NewSubscriptionHandler(mockService)

	productID := uuid.New()
	mockService.On("CreateSubscriptionPlan", productID.String(), "Monthly Plan", 30, 29.99).
		Return(nil, assert.AnError)

	req := &pb.CreateSubscriptionPlanRequest{
		ProductId: productID.String(),
		PlanName:  "Monthly Plan",
		Duration:  30,
		Price:     29.99,
	}

	resp, err := handler.CreateSubscriptionPlan(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockService.AssertExpectations(t)
}

func TestSubscriptionHandler_GetSubscriptionPlan_ServiceError(t *testing.T) {
	mockService := new(MockSubscriptionService)
	handler := NewSubscriptionHandler(mockService)

	planID := uuid.New()
	mockService.On("GetSubscriptionPlan", planID.String()).Return(nil, assert.AnError)

	req := &pb.GetSubscriptionPlanRequest{
		Id: planID.String(),
	}

	resp, err := handler.GetSubscriptionPlan(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockService.AssertExpectations(t)
}

func TestSubscriptionHandler_UpdateSubscriptionPlan_ServiceError(t *testing.T) {
	mockService := new(MockSubscriptionService)
	handler := NewSubscriptionHandler(mockService)

	planID := uuid.New()
	productID := uuid.New()
	mockService.On("UpdateSubscriptionPlan", planID.String(), productID.String(), "Updated Plan", 60, 49.99).
		Return(nil, assert.AnError)

	req := &pb.UpdateSubscriptionPlanRequest{
		Id:        planID.String(),
		ProductId: productID.String(),
		PlanName:  "Updated Plan",
		Duration:  60,
		Price:     49.99,
	}

	resp, err := handler.UpdateSubscriptionPlan(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockService.AssertExpectations(t)
}

func TestSubscriptionHandler_DeleteSubscriptionPlan_ServiceError(t *testing.T) {
	mockService := new(MockSubscriptionService)
	handler := NewSubscriptionHandler(mockService)

	planID := uuid.New()
	mockService.On("DeleteSubscriptionPlan", planID.String()).Return(assert.AnError)

	req := &pb.DeleteSubscriptionPlanRequest{
		Id: planID.String(),
	}

	resp, err := handler.DeleteSubscriptionPlan(context.Background(), req)

	assert.Error(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Success)
	assert.Contains(t, resp.Message, "error")
	mockService.AssertExpectations(t)
}

func TestSubscriptionHandler_ListSubscriptionPlans_ServiceError(t *testing.T) {
	mockService := new(MockSubscriptionService)
	handler := NewSubscriptionHandler(mockService)

	productID := uuid.New()
	mockService.On("ListSubscriptionPlans", productID.String()).Return(nil, assert.AnError)

	req := &pb.ListSubscriptionPlansRequest{
		ProductId: productID.String(),
	}

	resp, err := handler.ListSubscriptionPlans(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockService.AssertExpectations(t)
}
