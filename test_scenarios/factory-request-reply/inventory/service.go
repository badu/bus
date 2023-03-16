package inventory

import (
	"context"
)

type ServiceImpl struct {
}

func NewService() ServiceImpl {
	result := ServiceImpl{}

	return result
}

func (s *ServiceImpl) GetStockForProduct(ctx context.Context, productID string) {

}
