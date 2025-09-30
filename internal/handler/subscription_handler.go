package handler

import (
	"context"

	"github.com/microservice-go/product-service/internal/service"
	pb "github.com/microservice-go/product-service/proto/subscription"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SubscriptionHandler implements the gRPC SubscriptionService
type SubscriptionHandler struct {
	pb.UnimplementedSubscriptionServiceServer
	service service.SubscriptionService
}

// NewSubscriptionHandler creates a new SubscriptionHandler
func NewSubscriptionHandler(service service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// CreateSubscriptionPlan handles the CreateSubscriptionPlan RPC
func (h *SubscriptionHandler) CreateSubscriptionPlan(ctx context.Context, req *pb.CreateSubscriptionPlanRequest) (*pb.SubscriptionPlanResponse, error) {
	plan, err := h.service.CreateSubscriptionPlan(req.ProductId, req.PlanName, int(req.Duration), req.Price)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create subscription plan: %v", err)
	}

	return &pb.SubscriptionPlanResponse{
		Plan: &pb.SubscriptionPlan{
			Id:        plan.ID.String(),
			ProductId: plan.ProductID.String(),
			PlanName:  plan.PlanName,
			Duration:  int32(plan.Duration),
			Price:     plan.Price,
			CreatedAt: timestamppb.New(plan.CreatedAt),
			UpdatedAt: timestamppb.New(plan.UpdatedAt),
		},
	}, nil
}

// GetSubscriptionPlan handles the GetSubscriptionPlan RPC
func (h *SubscriptionHandler) GetSubscriptionPlan(ctx context.Context, req *pb.GetSubscriptionPlanRequest) (*pb.SubscriptionPlanResponse, error) {
	plan, err := h.service.GetSubscriptionPlan(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "subscription plan not found: %v", err)
	}

	return &pb.SubscriptionPlanResponse{
		Plan: &pb.SubscriptionPlan{
			Id:        plan.ID.String(),
			ProductId: plan.ProductID.String(),
			PlanName:  plan.PlanName,
			Duration:  int32(plan.Duration),
			Price:     plan.Price,
			CreatedAt: timestamppb.New(plan.CreatedAt),
			UpdatedAt: timestamppb.New(plan.UpdatedAt),
		},
	}, nil
}

// UpdateSubscriptionPlan handles the UpdateSubscriptionPlan RPC
func (h *SubscriptionHandler) UpdateSubscriptionPlan(ctx context.Context, req *pb.UpdateSubscriptionPlanRequest) (*pb.SubscriptionPlanResponse, error) {
	plan, err := h.service.UpdateSubscriptionPlan(req.Id, req.ProductId, req.PlanName, int(req.Duration), req.Price)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update subscription plan: %v", err)
	}

	return &pb.SubscriptionPlanResponse{
		Plan: &pb.SubscriptionPlan{
			Id:        plan.ID.String(),
			ProductId: plan.ProductID.String(),
			PlanName:  plan.PlanName,
			Duration:  int32(plan.Duration),
			Price:     plan.Price,
			CreatedAt: timestamppb.New(plan.CreatedAt),
			UpdatedAt: timestamppb.New(plan.UpdatedAt),
		},
	}, nil
}

// DeleteSubscriptionPlan handles the DeleteSubscriptionPlan RPC
func (h *SubscriptionHandler) DeleteSubscriptionPlan(ctx context.Context, req *pb.DeleteSubscriptionPlanRequest) (*pb.DeleteSubscriptionPlanResponse, error) {
	err := h.service.DeleteSubscriptionPlan(req.Id)
	if err != nil {
		return &pb.DeleteSubscriptionPlanResponse{
			Success: false,
			Message: err.Error(),
		}, status.Errorf(codes.Internal, "failed to delete subscription plan: %v", err)
	}

	return &pb.DeleteSubscriptionPlanResponse{
		Success: true,
		Message: "Subscription plan deleted successfully",
	}, nil
}

// ListSubscriptionPlans handles the ListSubscriptionPlans RPC
func (h *SubscriptionHandler) ListSubscriptionPlans(ctx context.Context, req *pb.ListSubscriptionPlansRequest) (*pb.ListSubscriptionPlansResponse, error) {
	plans, err := h.service.ListSubscriptionPlans(req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list subscription plans: %v", err)
	}

	pbPlans := make([]*pb.SubscriptionPlan, len(plans))
	for i, plan := range plans {
		pbPlans[i] = &pb.SubscriptionPlan{
			Id:        plan.ID.String(),
			ProductId: plan.ProductID.String(),
			PlanName:  plan.PlanName,
			Duration:  int32(plan.Duration),
			Price:     plan.Price,
			CreatedAt: timestamppb.New(plan.CreatedAt),
			UpdatedAt: timestamppb.New(plan.UpdatedAt),
		}
	}

	return &pb.ListSubscriptionPlansResponse{
		Plans: pbPlans,
		Total: int32(len(plans)),
	}, nil
}

