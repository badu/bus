package audit

import (
	"fmt"
	"strings"

	"github.com/badu/bus"
	"github.com/badu/bus/test_scenarios/fire-and-forget/events"
)

type ServiceImpl struct {
	sb *strings.Builder
}

func NewAuditService(sb *strings.Builder) ServiceImpl {
	result := ServiceImpl{sb: sb}
	bus.Sub(result.OnUserRegisteredEvent)
	bus.SubCancel(result.OnSMSRequestEvent)
	bus.SubCancel(result.OnSMSSentEvent)
	return result
}

// OnUserRegisteredEvent is classic event handler
func (s *ServiceImpl) OnUserRegisteredEvent(event events.UserRegisteredEvent) {
	// we can save audit data here
}

// OnSMSRequestEvent is a pub-unsub type, we have to return 'false' to continue listening for this kind of events
func (s *ServiceImpl) OnSMSRequestEvent(event events.SMSRequestEvent) bool {
	return false
}

// OnSMSSentEvent is a pub-unsub type where we give up on listening after receiving first message
func (s *ServiceImpl) OnSMSSentEvent(event events.SMSSentEvent) bool {
	s.sb.WriteString(fmt.Sprintf("audit event : an sms was %s sent to %s with message %s\n", event.Status, event.Request.Number, event.Request.Message))
	return true // after first event, audit will give up listening for events
}
