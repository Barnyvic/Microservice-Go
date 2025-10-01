package handler

import (
	"context"

	"github.com/microservice-go/product-service/internal/service"
	pb "github.com/microservice-go/product-service/proto/product"
)

type ProductHandler struct {
	pb.UnimplementedProductServiceServer
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	product, err := h.service.CreateProduct(req.Name, req.Description, req.Price, req.ProductType)
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.ProductResponse{
		Product: toProductProto(product),
	}, nil
}

func (h *ProductHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	product, err := h.service.GetProduct(req.Id)
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.ProductResponse{
		Product: toProductProto(product),
	}, nil
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	product, err := h.service.UpdateProduct(req.Id, req.Name, req.Description, req.Price, req.ProductType)
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.ProductResponse{
		Product: toProductProto(product),
	}, nil
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	err := h.service.DeleteProduct(req.Id)
	if err != nil {
		return &pb.DeleteProductResponse{
			Success: false,
			Message: err.Error(),
		}, mapServiceError(err)
	}

	return &pb.DeleteProductResponse{
		Success: true,
		Message: "Product deleted successfully",
	}, nil
}

func (h *ProductHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, total, err := h.service.ListProducts(req.ProductType, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, mapServiceError(err)
	}

	pbProducts := make([]*pb.Product, len(products))
	for i := range products {
		pbProducts[i] = toProductProto(&products[i])
	}

	return &pb.ListProductsResponse{
		Products: pbProducts,
		Total:    int32(total),
	}, nil
}
