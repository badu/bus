package request_reply_callback

import (
	"context"
	"strings"
	"testing"

	"github.com/badu/bus/test_scenarios/request-reply-callback/orders"
)

func TestRequestReplyCallback(t *testing.T) {
	var sb strings.Builder
	orders.NewRepository(&sb)
	svc := orders.NewService(&sb)

	ctx := context.Background()

	newOrder0, err := svc.RegisterOrder(ctx, []int{1, 2, 3})
	if err != nil {
		t.Fatalf("error creating order : %#v", err)
	}

	t.Logf("new order #0 : %#v", newOrder0)

	newOrder1, err := svc.RegisterOrder(ctx, []int{4, 5, 6})
	if err != nil {
		t.Fatalf("error creating order : %#v", err)
	}

	t.Logf("new order #1 : %#v", newOrder1)
	newOrder2, err := svc.RegisterOrder(ctx, []int{7, 8, 9})
	if err != nil {
		t.Fatalf("error creating order : %#v", err)
	}

	t.Logf("new order #2 : %#v", newOrder2)

	stat0, err := svc.GetOrderStatus(ctx, newOrder0.OrderID)
	if err != nil {
		t.Fatalf("error getting order status : %#v", err)
	}
	t.Logf("order #0 status : %s", stat0.Status)

	stat1, err := svc.GetOrderStatus(ctx, newOrder1.OrderID)
	if err != nil {
		t.Fatalf("error getting order status : %#v", err)
	}

	t.Logf("order #1 status : %s", stat1.Status)

	stat2, err := svc.GetOrderStatus(ctx, newOrder2.OrderID)
	if err != nil {
		t.Fatalf("error getting order status : %#v", err)
	}
	t.Logf("order #2 status : %s", stat2.Status)

	t.Logf("%s", sb.String())

}
