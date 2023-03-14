package orders

import (
	"context"
	"strings"

	"github.com/badu/bus"
	"github.com/badu/bus/test_scenarios/request-reply-callback/events"
)

type ServiceImpl struct {
	sb   *strings.Builder
	CBus *bus.Topic[*events.RequestEvent[Order]]
	SBus *bus.Topic[*events.RequestEvent[OrderStatus]]
}

func NewService(sb *strings.Builder) ServiceImpl {
	result := ServiceImpl{sb: sb}
	result.CBus = bus.NewTopic[*events.RequestEvent[Order]]()
	result.SBus = bus.NewTopic[*events.RequestEvent[OrderStatus]]()
	NewRepository(sb, result.CBus, result.SBus)
	return result
}

func (s *ServiceImpl) RegisterOrder(ctx context.Context, productIDs []int) (*Order, error) {
	event := events.RequestEvent[Order]{Payload: Order{ProductIDs: productIDs}, Done: make(chan struct{})}
	s.CBus.Pub(&event)
	<-event.Done
	return event.Callback()
}

func (s *ServiceImpl) GetOrderStatus(ctx context.Context, orderID int) (*OrderStatus, error) {
	event := events.RequestEvent[OrderStatus]{Payload: OrderStatus{OrderID: orderID}, Done: make(chan struct{})}
	s.SBus.Pub(&event)
	<-event.Done
	return event.Callback()
}
