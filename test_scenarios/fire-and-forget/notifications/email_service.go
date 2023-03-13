package notifications

import (
	"fmt"
	"testing"

	"github.com/badu/bus"
	"github.com/badu/bus/test_scenarios/fire-and-forget/events"
)

type EmailServiceImpl struct {
	t *testing.T
}

func NewEmailService(t *testing.T) EmailServiceImpl {
	result := EmailServiceImpl{t: t}
	bus.Listen(result.OnUserRegisteredEvent)
	return result
}

func (s *EmailServiceImpl) OnUserRegisteredEvent(e events.UserRegisteredEvent) bool {
	s.t.Logf("user %s has registered - sending welcome email message.", e.UserName)
	bus.Publish(events.SMSRequestEvent{
		Number:  e.Phone,
		Message: fmt.Sprintf("%s your user account was created. Check your email for instructions", e.UserName),
	})
	return false
}
