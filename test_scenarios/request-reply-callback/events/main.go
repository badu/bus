package events

type RequestEvent[T any] struct {
	Payload  T
	Callback func() (*T, error)
	Done     chan struct{}
}

func NewRequestEvent[T any](payload T) *RequestEvent[T] {
	return &RequestEvent[T]{
		Payload: payload,
		Done:    make(chan struct{}),
	}
}

func (i *RequestEvent[T]) Async() bool {
	return true // this one is async
}
