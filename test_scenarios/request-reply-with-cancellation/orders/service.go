package orders

import (
	"context"

	"github.com/badu/bus"
	"github.com/badu/bus/test_scenarios/request-reply-with-cancellation/events"
)

type ServiceImpl struct {
	bus *bus.Topic[events.CreateOrderEvent]
}

func NewService(bus *bus.Topic[events.CreateOrderEvent]) ServiceImpl {
	result := ServiceImpl{bus: bus}
	return result
}

func (s *ServiceImpl) CreateOrder(ctx context.Context, productIDs []int) (*Order, error) {
	event := events.CreateOrderEvent{State: events.NewEventState(ctx), ProductIDs: productIDs}
	s.bus.Pub(event)
	<-event.State.Done
	return &Order{ID: event.OrderID, ProductIDs: productIDs}, event.State.Error
}
