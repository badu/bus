package users

import (
	"context"

	"github.com/badu/bus"
	"github.com/badu/bus/test_scenarios/fire-and-forget/events"
)

type ServiceImpl struct {
}

func NewService() ServiceImpl {
	result := ServiceImpl{}

	return result
}

func (s *ServiceImpl) RegisterUser(ctx context.Context, name, phone string) {
	bus.Publish(events.UserRegisteredEvent{UserName: name, Phone: phone})
	bus.Publish(&events.DummyEvent{}) // nobody listens on this one
}
