package orders

import (
	"context"

	"github.com/badu/bus"
	"github.com/badu/bus/test_scenarios/request-reply-with-cancellation/events"
)

type ServiceImpl struct {
}

func NewService() ServiceImpl {
	result := ServiceImpl{}
	return result
}

func (s *ServiceImpl) CreateOrder(ctx context.Context, productIDs []int) (*Order, error) {
	event := events.CreateOrderEvent{State: events.NewEventState(ctx), ProductIDs: productIDs, NewOrder: &events.NewOrder{}}
	bus.Pub(event)
	<-event.State.Done

	if event.NewOrder != nil && event.State.Error == nil {
		return &Order{ID: event.NewOrder.ID, ProductIDs: productIDs}, nil
	}

	return nil, event.State.Error
}
