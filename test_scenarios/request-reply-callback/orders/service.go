package orders

import (
	"context"
	"fmt"
	"strings"

	"github.com/badu/bus"
	"github.com/badu/bus/test_scenarios/request-reply-callback/events"
)

type ServiceImpl struct {
	sb *strings.Builder
}

func NewService(sb *strings.Builder) ServiceImpl {
	result := ServiceImpl{sb: sb}
	return result
}

func (s *ServiceImpl) RegisterOrder(ctx context.Context, productIDs []int) (*Order, error) {
	event := events.NewRequestEvent[Order](Order{ProductIDs: productIDs})
	s.sb.WriteString(fmt.Sprintf("dispatching event typed %T\n", event))
	bus.Pub(event)
	<-event.Done            // wait for "reply"
	return event.Callback() // return the callback, which is containing the actual result
}

func (s *ServiceImpl) GetOrderStatus(ctx context.Context, orderID int) (*OrderStatus, error) {
	event := events.NewRequestEvent[OrderStatus](OrderStatus{OrderID: orderID})
	s.sb.WriteString(fmt.Sprintf("dispatching event typed %T\n", event))
	bus.Pub(event)
	<-event.Done            // wait for "reply"
	return event.Callback() // return the callback, which is containing the actual result
}
