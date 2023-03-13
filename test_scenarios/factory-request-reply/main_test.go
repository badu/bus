package factory_request_reply

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/badu/bus"
	"github.com/badu/bus/test_scenarios/factory-request-reply/cart"
	"github.com/badu/bus/test_scenarios/factory-request-reply/events"
	"github.com/badu/bus/test_scenarios/factory-request-reply/inventory"
	"github.com/badu/bus/test_scenarios/factory-request-reply/prices"
)

type pricesClientStub struct {
	prices.PricesServiceClient
}

func (s *pricesClientStub) GetPricesForProduct(ctx context.Context, in *prices.ProductIDRequest) (*prices.ProductPriceResponse, error) {
	return &prices.ProductPriceResponse{Price: 10.30}, nil
}

type fakeCloser struct {
}

func (f *fakeCloser) Close() error {
	return nil
}

func OnPricesGRPCClientStubRequest(e *events.PricesGRPCClientRequestEvent) bool {
	log.Println("OnPricesGRPCClientStubRequest", e)
	e.Client = &pricesClientStub{}
	e.Conn = &fakeCloser{}
	<-time.After(300 * time.Millisecond)
	e.Reply()
	return false
}

type inventoryClientStub struct {
	inventory.InventoryServiceClient
}

func (s *inventoryClientStub) GetStockForProduct(ctx context.Context, in *inventory.ProductIDRequest) (*inventory.ProductStockResponse, error) {
	return &inventory.ProductStockResponse{Stock: 200}, nil
}

func OnInventoryGRPCClientStubRequest(e *events.InventoryGRPCClientRequestEvent) bool {
	log.Println("OnInventoryGRPCClientStubRequest", e)
	e.Client = &inventoryClientStub{}
	e.Conn = &fakeCloser{}
	<-time.After(300 * time.Millisecond)
	e.Reply()
	return false
}

func TestGRPCClientStub(t *testing.T) {
	t.Log("GRPC client stub test")
	cartSvc := cart.NewService()

	bus.Listen(OnInventoryGRPCClientStubRequest)
	bus.Listen(OnPricesGRPCClientStubRequest)

	cartSvc.AddProductToCart(context.Background(), "1")
	cartSvc.AddProductToCart(context.Background(), "2")

	t.Log("GRPC client stubbing testing concluded")
}
