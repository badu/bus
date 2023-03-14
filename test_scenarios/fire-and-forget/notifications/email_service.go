package notifications

import (
	"fmt"
	"strings"

	"github.com/badu/bus"
	"github.com/badu/bus/test_scenarios/fire-and-forget/events"
)

type EmailServiceImpl struct {
	sb *strings.Builder
}

func NewEmailService(sb *strings.Builder) EmailServiceImpl {
	result := EmailServiceImpl{sb: sb}
	bus.Sub(result.OnUserRegisteredEvent)
	return result
}

func (s *EmailServiceImpl) OnUserRegisteredEvent(e events.UserRegisteredEvent) {
	s.sb.WriteString(fmt.Sprintf("user %s has registered - sending welcome email message\n", e.UserName))
	bus.Pub(events.SMSRequestEvent{
		Number:  e.Phone,
		Message: fmt.Sprintf("%s your user account was created. Check your email for instructions", e.UserName),
	})
}
