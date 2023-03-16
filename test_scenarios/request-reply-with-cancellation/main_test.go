package request_reply

import (
	"context"
	"testing"
	"time"

	"github.com/badu/bus/test_scenarios/request-reply-with-cancellation/orders"
)

func TestRequestReplyWithCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	svc := orders.NewService()
	orders.NewRepository()

	response, err := svc.CreateOrder(ctx, []int{1, 2, 3})
	switch err {
	default:
		t.Fatalf("error : it supposed to timeout, but it responded %#v and the error is %#v", response, err)
	case context.DeadlineExceeded:
		// what we were expecting
	}

	cancel()
}
