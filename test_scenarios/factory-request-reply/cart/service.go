package cart

import (
	"context"
	"log"

	"github.com/badu/bus"
	"github.com/badu/bus/test_scenarios/factory-request-reply/events"
	"github.com/badu/bus/test_scenarios/factory-request-reply/inventory"
	"github.com/badu/bus/test_scenarios/factory-request-reply/prices"
)

type ServiceImpl struct {
}

func NewService() ServiceImpl {
	result := ServiceImpl{}

	return result
}

func (s *ServiceImpl) AddProductToCart(ctx context.Context, productID string) error {
	e1 := events.NewInventoryGRPCClientRequestEvent()
	bus.Publish(e1)
	e1.WaitReply()

	e2 := events.NewPricesGRPCClientRequestEvent()
	bus.Publish(e2)
	e2.WaitReply()

	defer e1.Conn.Close() // close GRPC connection when done
	stockResponse, err := e1.Client.GetStockForProduct(ctx, &inventory.ProductIDRequest{ID: productID})
	if err != nil {
		return err
	}

	defer e2.Conn.Close() // close GRPC connection when done
	priceResponse, err := e2.Client.GetPricesForProduct(ctx, &prices.ProductIDRequest{ID: productID})
	if err != nil {
		return err
	}

	log.Println("stock", stockResponse.Stock, "price", priceResponse.Price)
	return nil
}
