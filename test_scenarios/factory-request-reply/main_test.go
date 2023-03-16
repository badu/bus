package factory_request_reply

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/badu/bus"
	"github.com/badu/bus/test_scenarios/factory-request-reply/cart"
	"github.com/badu/bus/test_scenarios/factory-request-reply/events"
	"github.com/badu/bus/test_scenarios/factory-request-reply/inventory"
	"github.com/badu/bus/test_scenarios/factory-request-reply/prices"
)

var sb strings.Builder

type pricesClientStub struct{}

func (s *pricesClientStub) GetPricesForProduct(ctx context.Context, in *prices.ProductIDRequest) (*prices.ProductPriceResponse, error) {
	return &prices.ProductPriceResponse{Price: 10.30}, nil
}

type fakeCloser struct {
}

func (f *fakeCloser) Close() error {
	return nil
}

func OnPricesGRPCClientStubRequest(e *events.PricesGRPCClientRequestEvent) {
	sb.WriteString("OnPricesGRPCClientStubRequest\n")
	e.Client = &pricesClientStub{}
	e.Conn = &fakeCloser{}
	<-time.After(300 * time.Millisecond)
	e.Reply()
}

type inventoryClientStub struct{}

func (s *inventoryClientStub) GetStockForProduct(ctx context.Context, in *inventory.ProductIDRequest) (*inventory.ProductStockResponse, error) {
	return &inventory.ProductStockResponse{Stock: 200}, nil
}

func OnInventoryGRPCClientStubRequest(e *events.InventoryGRPCClientRequestEvent) {
	sb.WriteString("OnInventoryGRPCClientStubRequest\n")
	e.Client = &inventoryClientStub{}
	e.Conn = &fakeCloser{}
	<-time.After(300 * time.Millisecond)
	e.Reply()
}

func TestGRPCClientStub(t *testing.T) {

	cartSvc := cart.NewService(&sb)

	bus.Sub(OnInventoryGRPCClientStubRequest)
	bus.Sub(OnPricesGRPCClientStubRequest)

	err := cartSvc.AddProductToCart(context.Background(), "1")
	if err != nil {
		t.Fatalf("error adding product to cart : %#v", err)
	}

	err = cartSvc.AddProductToCart(context.Background(), "2")
	if err != nil {
		t.Fatalf("error adding product to cart : %#v", err)
	}

	const expecting = "OnInventoryGRPCClientStubRequest\n" +
		"OnPricesGRPCClientStubRequest\n" +
		"stock 200.00pcs @ price 10.30$\n" +
		"OnInventoryGRPCClientStubRequest\n" +
		"OnPricesGRPCClientStubRequest\n" +
		"stock 200.00pcs @ price 10.30$\n"

	got := sb.String()
	if got != expecting {
		t.Fatalf("expecting :\n%s but got : \n%s", expecting, got)
	}
}
