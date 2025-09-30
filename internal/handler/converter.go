package handler

import (
	"strings"

	"github.com/microservice-go/product-service/internal/constants"
	"github.com/microservice-go/product-service/internal/models"
	productpb "github.com/microservice-go/product-service/proto/product"
	subscriptionpb "github.com/microservice-go/product-service/proto/subscription"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// toProductProto converts a Product model to a Product protobuf message
func toProductProto(product *models.Product) *productpb.Product {
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

	errMsg := err.Error()

	// Map validation errors to InvalidArgument
	switch errMsg {
	case constants.ErrProductNameRequired,
		constants.ErrPriceNegative,
		constants.ErrProductTypeRequired,
		constants.ErrPlanNameRequired,
		constants.ErrDurationPositive,
		constants.ErrInvalidProductID,
		constants.ErrInvalidPlanID:
		return status.Error(codes.InvalidArgument, errMsg)

	// Map not found errors to NotFound
	case constants.ErrProductNotFound,
		constants.ErrPlanNotFound:
		return status.Error(codes.NotFound, errMsg)

	// Default to Internal for unknown errors
	default:
		if strings.Contains(errMsg, "not found") {
			return status.Error(codes.NotFound, errMsg)
		}
		return status.Error(codes.Internal, errMsg)
	}
}

