package users

import (
	"context"
	"strings"

	"github.com/badu/bus"
	"github.com/badu/bus/test_scenarios/fire-and-forget/events"
)

type ServiceImpl struct {
	sb *strings.Builder
	c  int
}

func NewService(sb *strings.Builder) ServiceImpl {
	result := ServiceImpl{sb: sb}
	return result
}

func (s *ServiceImpl) RegisterUser(ctx context.Context, name, phone string) {
	s.c++
	bus.Pub(events.UserRegisteredEvent{UserName: name, Phone: phone})
	bus.Pub(&events.DummyEvent{AlteredAsync: s.c%2 == 0}) // nobody listens on this one
}
