package notifications

import (
	"fmt"
	"strings"

	"github.com/badu/bus"
	"github.com/badu/bus/test_scenarios/fire-and-forget/events"
)

type SmsServiceImpl struct {
	sb *strings.Builder
}

func NewSmsService(sb *strings.Builder) SmsServiceImpl {
	result := SmsServiceImpl{sb: sb}
	bus.Sub(result.OnSMSSendRequest)
	return result
}

func (s *SmsServiceImpl) OnSMSSendRequest(event events.SMSRequestEvent) {
	s.sb.WriteString(fmt.Sprintf("sms sent requested for number %s with message %s\n", event.Number, event.Message))
	bus.Pub(events.SMSSentEvent{
		Request: event,
		Status:  "successfully sent",
	})
}
