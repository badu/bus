package fire_and_forget

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/badu/bus/test_scenarios/fire-and-forget/audit"
	"github.com/badu/bus/test_scenarios/fire-and-forget/notifications"
	"github.com/badu/bus/test_scenarios/fire-and-forget/users"
)

func TestUserRegistration(t *testing.T) {
	var sb strings.Builder

	userSvc := users.NewService(&sb)
	notifications.NewSmsService(&sb)
	notifications.NewEmailService(&sb)
	audit.NewAuditService(&sb)

	userSvc.RegisterUser(context.Background(), "Badu", "+40742222222")

	<-time.After(500 * time.Millisecond)

	userSvc.RegisterUser(context.Background(), "Adina", "+40743333333")

	<-time.After(500 * time.Millisecond)

	const expecting = "user Badu has registered - sending welcome email message\n" +
		"sms sent requested for number +40742222222 with message Badu your user account was created. Check your email for instructions\n" +
		"audit event : an sms was successfully sent sent to +40742222222 with message Badu your user account was created. Check your email for instructions\n" +
		"user Adina has registered - sending welcome email message\n" +
		"sms sent requested for number +40743333333 with message Adina your user account was created. Check your email for instructions\n"

	got := sb.String()
	if got != expecting {
		t.Fatalf("expecting :\n%s but got : \n%s", expecting, got)
	}
}
