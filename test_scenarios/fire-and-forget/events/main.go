package events

type UserRegisteredEvent struct {
	UserName string
	Phone    string
}

func (e UserRegisteredEvent) Async() bool {
	return true
}

type SMSRequestEvent struct {
	Number  string
	Message string
}

func (e SMSRequestEvent) Async() bool {
	return true
}

type SMSSentEvent struct {
	Request SMSRequestEvent
	Status  string
}

func (e SMSSentEvent) Async() bool {
	return true
}

type DummyEvent struct {
	AlteredAsync bool
}

func (e *DummyEvent) Async() bool {
	return e.AlteredAsync
}
