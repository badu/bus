package events

import (
	"context"
)

type EventState struct {
	Ctx   context.Context
	Done  chan struct{} `json:"-"`
	Error error
}

func NewEventState(ctx context.Context) *EventState {
	return &EventState{
		Ctx:  ctx,
		Done: make(chan struct{}),
	}
}

func (s *EventState) Close() {
	s.Error = s.Ctx.Err()
	close(s.Done)
}

type NewOrder struct {
	ID int
}

type CreateOrderEvent struct {
	NewOrder   *NewOrder
	ProductIDs []int
	State      *EventState
}

func (c CreateOrderEvent) EventID() string {
	return "CreateOrderEventType"
}

func (c CreateOrderEvent) Async() bool {
	return false
}
