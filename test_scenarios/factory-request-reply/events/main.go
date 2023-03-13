package events

import (
	"sync"

	"github.com/badu/bus/test_scenarios/factory-request-reply/inventory"
	"github.com/badu/bus/test_scenarios/factory-request-reply/prices"
)

const (
	InventoryGRPCClientRequestEventType = "InventoryGRPCClientRequestEvent"
	PricesGRPCClientRequestEventType    = "PricesGRPCClientRequestEvent"
)

type InventoryGRPCClientRequestEvent struct {
	wg     sync.WaitGroup
	Conn   Closer // should be *grpc.ClientConn, but we're avoiding the import
	Client inventory.InventoryServiceClient
}

func NewInventoryGRPCClientRequestEvent() *InventoryGRPCClientRequestEvent {
	result := InventoryGRPCClientRequestEvent{}
	result.wg.Add(1)
	return &result
}

func (i *InventoryGRPCClientRequestEvent) EventID() string {
	return InventoryGRPCClientRequestEventType
}

func (i *InventoryGRPCClientRequestEvent) Async() bool {
	return true // this one is async
}

func (i *InventoryGRPCClientRequestEvent) WaitReply() {
	i.wg.Wait()
}

func (i *InventoryGRPCClientRequestEvent) Reply() {
	i.wg.Done()
}

type PricesGRPCClientRequestEvent struct {
	wg     sync.WaitGroup
	Conn   Closer // should be *grpc.ClientConn, but we're avoiding the import
	Client prices.PricesServiceClient
}

func NewPricesGRPCClientRequestEvent() *PricesGRPCClientRequestEvent {
	result := PricesGRPCClientRequestEvent{}
	result.wg.Add(1)
	return &result
}

func (p *PricesGRPCClientRequestEvent) EventID() string {
	return PricesGRPCClientRequestEventType
}

func (p *PricesGRPCClientRequestEvent) Async() bool {
	return false // this one is sync
}

func (p *PricesGRPCClientRequestEvent) WaitReply() {
	p.wg.Wait()
}

func (p *PricesGRPCClientRequestEvent) Reply() {
	p.wg.Done()
}

type Closer interface {
	Close() error
}
