package events

const (
	UserRegisteredEventType string = "UserRegisteredEvent"
	SMSRequestEventType     string = "SMSRequestEvent"
	SMSSentEventType        string = "SmsSentEvent"
	DummyEventType          string = "DummyEvent"
)

type UserRegisteredEvent struct {
	UserName string
	Phone    string
}

func (e UserRegisteredEvent) EventID() string {
	return UserRegisteredEventType
}

func (e UserRegisteredEvent) Async() bool {
	return true
}

type SMSRequestEvent struct {
	Number  string
	Message string
}

func (e SMSRequestEvent) EventID() string {
	return SMSRequestEventType
}

func (e SMSRequestEvent) Async() bool {
	return true
}

type SMSSentEvent struct {
	Request SMSRequestEvent
	Status  string
}

func (e SMSSentEvent) EventID() string {
	return SMSSentEventType
}

func (e SMSSentEvent) Async() bool {
	return true
}

type DummyEvent struct {
}

func (e *DummyEvent) EventID() string {
	return DummyEventType
}

func (e *DummyEvent) Async() bool {
	return true
}
