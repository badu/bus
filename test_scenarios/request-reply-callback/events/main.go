package events

const RequestEventType = "RequestEvent"

type RequestEvent[T any] struct {
	Payload  T
	Callback func() (*T, error)
	Done     chan struct{}
}

func (i *RequestEvent[T]) EventID() string {
	return RequestEventType
}

func (i *RequestEvent[T]) Async() bool {
	return true // this one is async
}
