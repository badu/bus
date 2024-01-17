package bus_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/badu/bus"
)

func TestSubTopicWhilePub(t *testing.T) {
	// scenario : we have a large number of subscribers.
	// we publish an event and while doing that,
	// we register another one on a different goroutine

	topic := bus.NewTopic[*Uint32AsyncEvent]()
	for i := 0; i < 4096; i++ {
		topic.Sub(func(v *Uint32AsyncEvent) {})
	}

	finishPubWait := make(chan struct{})
	finishSubWait := make(chan struct{})
	start := make(chan struct{})

	go func() {
		<-start
		topic.PubAsync(&Uint32AsyncEvent{u: 1})
		defer close(finishPubWait)
	}()

	newSubCalled := false

	go func() {
		<-start
		topic.Sub(func(v *Uint32AsyncEvent) {
			newSubCalled = true
		})
		close(finishSubWait)
	}()

	close(start) // start both goroutines

	<-finishPubWait // wait for pub to finish

	<-finishSubWait // wait for sub to finish

	if newSubCalled {
		t.Fatal("new subscriber should not be called")
	}
}

func TestReusePayloadPointerAsync(t *testing.T) {
	// if you reuse the payload, you can alter it's content, of course

	topic := bus.NewTopic[*Uint32AsyncEvent]()
	c := uint32(0)
	for i := 0; i < 4096; i++ {
		k := i
		topic.Sub(func(v *Uint32AsyncEvent) {
			if v.u == 2048 {
				atomic.AddUint32(&c, 1)
				return
			}

			if k == 2048 {
				v.u = 2048
			}
		})
	}

	finishPubWait := make(chan struct{})

	payload := Uint32AsyncEvent{u: 1}

	topic.Pub(&payload)

	close(finishPubWait)

	<-finishPubWait // wait for pub to finish

	t.Logf("altered payload %d for %d listeners", payload.u, c)
}

func TestAsyncBus(t *testing.T) {
	c := uint32(0)

	var wg sync.WaitGroup
	wg.Add(4096)
	bus.Sub(
		func(event Uint32AsyncEvent) {
			atomic.AddUint32(&c, 1)
			wg.Done()
		},
	)

	go func() {
		for i := 0; i < 1024; i++ {
			bus.PubAsync(Uint32AsyncEvent{})
		}
	}()
	go func() {
		for i := 0; i < 1024; i++ {
			bus.PubAsync(Uint32AsyncEvent{})
		}
	}()
	go func() {
		for i := 0; i < 1024; i++ {
			bus.PubAsync(Uint32AsyncEvent{})
		}
	}()
	go func() {
		for i := 0; i < 1024; i++ {
			bus.PubAsync(Uint32AsyncEvent{})
		}
	}()

	wg.Wait()

	if c != 4096 {
		t.Fatalf("error : counter should be 4096 but is %d", c)
	}

	t.Logf("%d", c)
}

func TestRange(t *testing.T) {

	type Event1 struct{}
	type Event2 struct{}
	type Event3 struct{}
	type Event4 struct{}
	type Event5 struct{}

	bus.Sub(func(e Event1) {})
	bus.Sub(func(e Event2) {})
	bus.Sub(func(e Event2) {})
	bus.Sub(func(e Event3) {})
	bus.Sub(func(e Event3) {})
	bus.Sub(func(e Event3) {})
	bus.Sub(func(e Event4) {})
	bus.Sub(func(e Event4) {})
	bus.Sub(func(e Event4) {})
	bus.Sub(func(e Event4) {})
	bus.Sub(func(e Event5) {})
	bus.Sub(func(e Event5) {})
	bus.Sub(func(e Event5) {})
	bus.Sub(func(e Event5) {})
	bus.Sub(func(e Event5) {})

	seen := map[string]struct{}{
		"bus_test.Event2": {},
		"bus_test.Event3": {},
		"bus_test.Event1": {},
		"bus_test.Event5": {},
		"bus_test.Event4": {},
	}

	bus.DefaultBus.Range(func(k, _ any) bool {
		if _, has := seen[k.(string)]; has {
			delete(seen, k.(string))
		}
		return true
	})

	if len(seen) > 0 {
		t.Fatalf("error : not all events were seen")
	}
}
