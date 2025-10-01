package handler

import (
	"context"

	"github.com/microservice-go/product-service/internal/service"
	pb "github.com/microservice-go/product-service/proto/subscription"
)

type SubscriptionHandler struct {
	pb.UnimplementedSubscriptionServiceServer
	service service.SubscriptionService
}

func NewSubscriptionHandler(service service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

func (h *SubscriptionHandler) CreateSubscriptionPlan(ctx context.Context, req *pb.CreateSubscriptionPlanRequest) (*pb.SubscriptionPlanResponse, error) {
	plan, err := h.service.CreateSubscriptionPlan(req.ProductId, req.PlanName, int(req.Duration), req.Price)
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.SubscriptionPlanResponse{
		Plan: toSubscriptionPlanProto(plan),
	}, nil
}

func (h *SubscriptionHandler) GetSubscriptionPlan(ctx context.Context, req *pb.GetSubscriptionPlanRequest) (*pb.SubscriptionPlanResponse, error) {
	plan, err := h.service.GetSubscriptionPlan(req.Id)
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.SubscriptionPlanResponse{
		Plan: toSubscriptionPlanProto(plan),
	}, nil
}

func (h *SubscriptionHandler) UpdateSubscriptionPlan(ctx context.Context, req *pb.UpdateSubscriptionPlanRequest) (*pb.SubscriptionPlanResponse, error) {
	plan, err := h.service.UpdateSubscriptionPlan(req.Id, req.ProductId, req.PlanName, int(req.Duration), req.Price)
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.SubscriptionPlanResponse{
		Plan: toSubscriptionPlanProto(plan),
	}, nil
}

func (h *SubscriptionHandler) DeleteSubscriptionPlan(ctx context.Context, req *pb.DeleteSubscriptionPlanRequest) (*pb.DeleteSubscriptionPlanResponse, error) {
	err := h.service.DeleteSubscriptionPlan(req.Id)
	if err != nil {
		return &pb.DeleteSubscriptionPlanResponse{
			Success: false,
			Message: err.Error(),
		}, mapServiceError(err)
	}

	return &pb.DeleteSubscriptionPlanResponse{
		Success: true,
		Message: "Subscription plan deleted successfully",
	}, nil
}

func (h *SubscriptionHandler) ListSubscriptionPlans(ctx context.Context, req *pb.ListSubscriptionPlansRequest) (*pb.ListSubscriptionPlansResponse, error) {
	plans, err := h.service.ListSubscriptionPlans(req.ProductId)
	if err != nil {
		return nil, mapServiceError(err)
	}

	pbPlans := make([]*pb.SubscriptionPlan, len(plans))
	for i := range plans {
		pbPlans[i] = toSubscriptionPlanProto(&plans[i])
	}

	return &pb.ListSubscriptionPlansResponse{
		Plans: pbPlans,
		Total: int32(len(plans)),
	}, nil
}
