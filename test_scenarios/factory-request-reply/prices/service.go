package prices

import (
	"context"
)

type ServiceImpl struct {
}

func NewService() ServiceImpl {
	result := ServiceImpl{}

	return result
}

func (s *ServiceImpl) GetPricesForProduct(ctx context.Context, productID string) {

}
