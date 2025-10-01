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

func mapServiceError(err error) error {
	if err == nil {
		return nil
	}


	if apperrors.IsValidationError(err) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if apperrors.IsNotFoundError(err) {
		return status.Error(codes.NotFound, err.Error())
	}

	if apperrors.IsDatabaseError(err) {
		return status.Error(codes.Internal, err.Error())
	}

	return status.Error(codes.Internal, err.Error())
}

