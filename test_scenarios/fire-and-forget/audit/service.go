package audit

import (
	"testing"

	"github.com/badu/bus"
	"github.com/badu/bus/test_scenarios/fire-and-forget/events"
)

type AuditServiceImpl struct {
	t  *testing.T
	bl *bus.Bus[events.SMSSentEvent]
}

func NewAuditService(t *testing.T) AuditServiceImpl {
	result := AuditServiceImpl{t: t}
	result.bl = bus.Listen(result.OnSMSSentEvent)
	return result
}

func (s *AuditServiceImpl) OnSMSSentEvent(event events.SMSSentEvent) bool {
	s.t.Logf("audit event : an sms was %s sent to %s with message %s", event.Status, event.Request.Number, event.Request.Message)
	return false
}

func (s *AuditServiceImpl) Stop() {
	s.bl.Unsub()
}
