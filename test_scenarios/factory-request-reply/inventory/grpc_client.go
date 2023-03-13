package inventory

import (
	"context"
)

type InventoryServiceClient interface {
	GetStockForProduct(ctx context.Context, in *ProductIDRequest) (*ProductStockResponse, error) // , opts ...grpc.CallOption) (*ProductStockResponse, error)
}

type ProductIDRequest struct {
	ID string
}

type ProductStockResponse struct {
	Stock float64
}
