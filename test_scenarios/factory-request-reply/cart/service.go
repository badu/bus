package cart

import (
	"context"
	"fmt"
	"strings"

	"github.com/badu/bus"
	"github.com/badu/bus/test_scenarios/factory-request-reply/events"
	"github.com/badu/bus/test_scenarios/factory-request-reply/inventory"
	"github.com/badu/bus/test_scenarios/factory-request-reply/prices"
)

type ServiceImpl struct {
	sb *strings.Builder
}

func NewService(sb *strings.Builder) ServiceImpl {
	result := ServiceImpl{sb: sb}
	return result
}

func (s *ServiceImpl) AddProductToCart(ctx context.Context, productID string) error {
	inventoryClientRequest := events.NewInventoryGRPCClientRequestEvent()
	bus.Pub(inventoryClientRequest)
	inventoryClientRequest.WaitReply()

	pricesClientRequest := events.NewPricesGRPCClientRequestEvent()
	bus.Pub(pricesClientRequest)
	pricesClientRequest.WaitReply()

	defer inventoryClientRequest.Conn.Close() // close GRPC connection when done
	stockResponse, err := inventoryClientRequest.Client.GetStockForProduct(ctx, &inventory.ProductIDRequest{ID: productID})
	if err != nil {
		return err
	}

	defer pricesClientRequest.Conn.Close() // close GRPC connection when done
	priceResponse, err := pricesClientRequest.Client.GetPricesForProduct(ctx, &prices.ProductIDRequest{ID: productID})
	if err != nil {
		return err
	}

	s.sb.WriteString(fmt.Sprintf("stock %0.2fpcs @ price %0.2f$\n", stockResponse.Stock, priceResponse.Price))

	return nil
}
