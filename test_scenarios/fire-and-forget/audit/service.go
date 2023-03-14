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
	bus.SubUnsub(result.OnSMSSentEvent)
	return result
}

func (s *ServiceImpl) OnSMSSentEvent(event events.SMSSentEvent) bool {
	s.sb.WriteString(fmt.Sprintf("audit event : an sms was %s sent to %s with message %s\n", event.Status, event.Request.Number, event.Request.Message))
	return true // after first event, audit will give up listening for events
}
