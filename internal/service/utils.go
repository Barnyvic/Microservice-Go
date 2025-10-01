package service

import (
	"github.com/google/uuid"
	apperrors "github.com/microservice-go/product-service/internal/errors"
)

func parseProductID(id string) (uuid.UUID, error) {
	if id == "" {
		return uuid.Nil, apperrors.NewValidationError("productId", "product ID is required")
	}

	productID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, apperrors.NewValidationError("productId", "invalid product ID format")
	}

	return productID, nil
}

