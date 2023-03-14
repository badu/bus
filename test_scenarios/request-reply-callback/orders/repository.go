package orders

import (
	"strings"
	"time"

	"github.com/badu/bus"
	"github.com/badu/bus/test_scenarios/request-reply-callback/events"
)

type Order struct {
	OrderID    int
	ProductIDs []int
}

type OrderStatus struct {
	OrderID int
	Status  string
}

type RepositoryImpl struct {
	sb    *strings.Builder
	calls int
}

func NewRepository(
	sb *strings.Builder,
	cBus *bus.Topic[*events.RequestEvent[Order]],
	sBus *bus.Topic[*events.RequestEvent[OrderStatus]],
) RepositoryImpl {
	result := RepositoryImpl{sb: sb}
	cBus.Sub(result.onCreateOrder)
	sBus.Sub(result.onGetOrderStatus)
	return result
}

func (r *RepositoryImpl) onCreateOrder(event *events.RequestEvent[Order]) {
	defer func() { r.calls++ }()

	<-time.After(500 * time.Millisecond) // simulate heavy database call

	event.Callback = func() (*Order, error) {
		return &Order{OrderID: r.calls, ProductIDs: event.Payload.ProductIDs}, nil
	}

	close(event.Done)
}

func (r *RepositoryImpl) onGetOrderStatus(event *events.RequestEvent[OrderStatus]) {
	<-time.After(300 * time.Millisecond) // simulate heavy database call

	event.Callback = func() (*OrderStatus, error) {
		status := "in_progress"
		if event.Payload.OrderID == 3 {
			status = "cancelled"
		}
		return &OrderStatus{OrderID: event.Payload.OrderID, Status: status}, nil
	}

	close(event.Done)
}
