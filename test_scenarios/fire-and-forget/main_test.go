package fire_and_forget

import (
	"context"
	"testing"
	"time"

	"github.com/badu/bus/test_scenarios/fire-and-forget/audit"
	"github.com/badu/bus/test_scenarios/fire-and-forget/notifications"
	"github.com/badu/bus/test_scenarios/fire-and-forget/users"
)

func TestUserRegistration(t *testing.T) {
	t.Log("Fire and forget test")

	userSvc := users.NewService()
	notifications.NewSmsService(t)
	notifications.NewEmailService(t)
	auditSvc := audit.NewAuditService(t)
	userSvc.RegisterUser(context.Background(), "Badu", "+40742222222")
	<-time.After(500 * time.Millisecond)

	// audit will be missing
	auditSvc.Stop()
	userSvc.RegisterUser(context.Background(), "Adina", "+40743333333")
	<-time.After(500 * time.Millisecond)

	t.Log("fire and forget testing concluded")
}
