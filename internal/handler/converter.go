package handler

import (
	apperrors "github.com/microservice-go/product-service/internal/errors"
	"github.com/microservice-go/product-service/internal/models"
	productpb "github.com/microservice-go/product-service/proto/product"
	subscriptionpb "github.com/microservice-go/product-service/proto/subscription"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// toProductProto converts a Product model to a Product protobuf message
func toProductProto(product *models.Product) *productpb.Product {
	if product == nil {
		return nil
	}

	return &productpb.Product{
		Id:          product.ID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		ProductType: product.ProductType,
		CreatedAt:   timestamppb.New(product.CreatedAt),
		UpdatedAt:   timestamppb.New(product.UpdatedAt),
	}
}

// toSubscriptionPlanProto converts a SubscriptionPlan model to a SubscriptionPlan protobuf message
func toSubscriptionPlanProto(plan *models.SubscriptionPlan) *subscriptionpb.SubscriptionPlan {
	if plan == nil {
		return nil
	}

	return &subscriptionpb.SubscriptionPlan{
		Id:        plan.ID.String(),
		ProductId: plan.ProductID.String(),
		PlanName:  plan.PlanName,
		Duration:  int32(plan.Duration),
		Price:     plan.Price,
		CreatedAt: timestamppb.New(plan.CreatedAt),
		UpdatedAt: timestamppb.New(plan.UpdatedAt),
	}
}

// mapServiceError maps service layer errors to appropriate gRPC status codes
func mapServiceError(err error) error {
	if err == nil {
		return nil
	}

	// Check for custom error types
	if apperrors.IsValidationError(err) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if apperrors.IsNotFoundError(err) {
		return status.Error(codes.NotFound, err.Error())
	}

	if apperrors.IsDatabaseError(err) {
		return status.Error(codes.Internal, err.Error())
	}

	// Default to Internal for unknown errors
	return status.Error(codes.Internal, err.Error())
}

