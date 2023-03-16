package orders

import (
	"time"

	"github.com/badu/bus"
	"github.com/badu/bus/test_scenarios/request-reply-with-cancellation/events"
)

type Order struct {
	ID         int
	ProductIDs []int
}

type RepositoryImpl struct {
	calls int
}

func NewRepository() RepositoryImpl {
	result := RepositoryImpl{}
	bus.Sub(result.OnCreateOrder)
	return result
}

func (r *RepositoryImpl) OnCreateOrder(event events.CreateOrderEvent) {
	defer func() {
		r.calls++
	}()

	for {
		select {
		case <-time.After(4 * time.Second):
			event.NewOrder = &events.NewOrder{ID: r.calls}
			event.State.Close()
			return
		case <-event.State.Ctx.Done():
			event.State.Close()
			return
		}
	}
}
