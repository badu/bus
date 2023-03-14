package users

import (
	"context"
	"strings"

	"github.com/badu/bus"
	"github.com/badu/bus/test_scenarios/fire-and-forget/events"
)

type ServiceImpl struct {
	sb *strings.Builder
}

func NewService(sb *strings.Builder) ServiceImpl {
	result := ServiceImpl{sb: sb}
	return result
}

func (s *ServiceImpl) RegisterUser(ctx context.Context, name, phone string) {
	bus.Pub(events.UserRegisteredEvent{UserName: name, Phone: phone})
	bus.Pub(&events.DummyEvent{}) // nobody listens on this one
}
