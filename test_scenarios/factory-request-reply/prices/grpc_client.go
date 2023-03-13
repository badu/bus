package prices

import (
	"context"
)

type PricesServiceClient interface {
	GetPricesForProduct(ctx context.Context, in *ProductIDRequest) (*ProductPriceResponse, error) // , opts ...grpc.CallOption) (*ProductPriceResponse, error)
}

type ProductIDRequest struct {
	ID string
}

type ProductPriceResponse struct {
	Price float64
}
