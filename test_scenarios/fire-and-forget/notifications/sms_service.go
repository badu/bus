package notifications

import (
	"testing"

	"github.com/badu/bus"
	"github.com/badu/bus/test_scenarios/fire-and-forget/events"
)

type SmsServiceImpl struct {
	t *testing.T
}

func NewSmsService(t *testing.T) SmsServiceImpl {
	result := SmsServiceImpl{t: t}
	bus.Listen(result.OnSMSSendRequest)
	return result
}

func (s *SmsServiceImpl) OnSMSSendRequest(event events.SMSRequestEvent) bool {
	s.t.Logf("sms sent requested for number %q with message %q", event.Number, event.Message)
	bus.Publish(events.SMSSentEvent{
		Request: event,
		Status:  "successfully sent",
	})
	return false
}
