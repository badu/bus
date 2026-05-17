package events

type UserRegisteredEvent struct {
	UserName string
	Phone    string
}

type SMSRequestEvent struct {
	Number  string
	Message string
}

type SMSSentEvent struct {
	Request SMSRequestEvent
	Status  string
}

type DummyEvent struct {
	Altered bool
}
