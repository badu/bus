package bus

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

var mapper sync.Map // holds key (event name - string) versus topic values

// we allow developers to override event names. They should be careful about name collisions
type iEventName interface {
	EventID() string //
}

// Listener is being returned when you subscribe to a topic, so you can unsubscribe or access the parent topic
type Listener[T any] struct {
	parent   *Topic[T]     // so we can call unsubscribe from parent
	callback func(event T) // the function that we're going to call
}

// Topic keeps the subscribers of one topic
type Topic[T any] struct {
	subs      []*Listener[T] // list of listeners
	mu        sync.RWMutex   // guards subs
	lisnsPool sync.Pool      // a pool of listeners
}

// NewTopic creates a new topic for a specie of events
func NewTopic[T any]() *Topic[T] {
	var result *Topic[T]
	result = &Topic[T]{
		lisnsPool: sync.Pool{
			New: func() any {
				return &Listener[T]{parent: result}
			},
		},
	}
	return result
}

// Sub adds a callback to be called when an event of that type is being published
func (t *Topic[T]) Sub(callback func(event T)) *Listener[T] {
	result, ok := t.lisnsPool.Get().(*Listener[T])
	if !ok {
		fmt.Println("[Topic][Sub] failed to allocate listener")
		return nil
	}
	result.callback = callback
	result.parent = t

	t.mu.Lock()
	t.subs = append(t.subs, result)
	t.mu.Unlock()

	return result
}

// cancel is private to the topic, but can be accessed via Listener
func (t *Topic[T]) cancel(which *Listener[T]) {
	t.mu.Lock()
	for i := range t.subs {
		if t.subs[i] != which {
			continue
		}

		t.subs[i] = t.subs[len(t.subs)-1]
		t.subs[len(t.subs)-1] = nil
		t.subs = t.subs[:len(t.subs)-1]
		break
	}
	t.mu.Unlock()

	which.callback = nil
	t.lisnsPool.Put(which)
}

// NumSubs in case you need to perform tests and check the number of subscribers of this particular topic
func (t *Topic[T]) NumSubs() int {
	t.mu.RLock()
	result := len(t.subs)
	t.mu.RUnlock()
	return result
}

// Cancel forgets the indicated callback
func (s *Listener[T]) Cancel() {
	s.parent.cancel(s)
}

// Topic gives access to the underlying topic
func (s *Listener[T]) Topic() *Topic[T] {
	return s.parent
}

// Pub allows you to publish an event in that topic
func (t *Topic[T]) Pub(event T) {
	t.mu.RLock()
	subsCopy := make([]*Listener[T], len(t.subs))
	copy(subsCopy, t.subs)
	t.mu.RUnlock()

	for _, sub := range subsCopy {
		if sub.callback == nil {
			continue
		}

		sub.callback(event)
	}
}

// Bus is being returned when you subscribe, so you can manually Cancel
type Bus[T any] struct {
	listener *Listener[T]
	stop     atomic.Uint32 // flag for unsubscribing after receiving one event
}

// Cancel allows callers to manually unsubscribe, in case they don't want to use SubCancel
func (o *Bus[T]) Cancel() {
	if o.stop.CompareAndSwap(0, 1) {
		go o.listener.Cancel()
	}
}

// SubCancel can be used if you need to unsubscribe immediately after receiving an event, by making your function return true
// recovererInfo is declared as variadic, in order not to break the previous version, but only the first element is taken into consideration.
// recovererInfo is a function that gets called if the go routine which does the callback panics
func SubCancel[T any](callback func(event T) bool, recoverInfo ...func(panicError error, panicStack ...[]byte)) *Bus[T] {
	key := getKey[T]()

	existingTopic, ok := mapper.Load(key)
	if !ok || existingTopic == nil {
		newTopic := NewTopic[T]()
		mapper.Store(key, newTopic)
		existingTopic = newTopic
	}

	topic, ok := existingTopic.(*Topic[T])
	if !ok {
		pc, _, _, pcOk := runtime.Caller(1)
		if pcOk {
			f := runtime.FuncForPC(pc)
			if f != nil {
				fmt.Println("[SubCancel] no such topic, called from ", f.Name())
			}
		} else {
			fmt.Println("[SubCancel] no such topic on unknown caller")
		}
		return nil
	}

	var result Bus[T]
	result.listener = topic.Sub(
		func(v T) {
			go func() {
				defer func() {
					if p := recover(); p != nil {
						recovered, isError := p.(error)
						if len(recoverInfo) > 0 {
							if isError {
								recoverInfo[0](recovered, debug.Stack())
							} else {
								recoverInfo[0](recovered)
							}
						}
					}
				}()

				if result.stop.Load() == 1 {
					return
				}

				shouldCancel := callback(v)
				if shouldCancel {
					result.Cancel()
				}
			}()
		},
	)

	return &result
}

// Sub subscribes a callback function to listen for a specie of events
// recovererInfo is declared as variadic, in order not to break the previous version, but only the first element is taken into consideration.
// recovererInfo is a function that gets called if the go routine which does the callback panics
func Sub[T any](callback func(event T), recovererInfo ...func(panicError error, panicStack ...[]byte)) *Bus[T] {
	key := getKey[T]()

	existingTopic, ok := mapper.Load(key)
	if !ok || existingTopic == nil {
		newTopic := NewTopic[T]()
		mapper.Store(key, newTopic)
		existingTopic = newTopic
	}

	topic, ok := existingTopic.(*Topic[T])
	if !ok {
		pc, _, _, pcOk := runtime.Caller(1)
		if pcOk {
			f := runtime.FuncForPC(pc)
			if f != nil {
				fmt.Println("[Sub] no such topic, called from ", f.Name())
			}
		} else {
			fmt.Println("[Sub] no such topic on unknown caller")
		}
		return nil
	}

	var result Bus[T]
	result.listener = topic.Sub(
		func(v T) {
			go func() {
				defer func() {
					if p := recover(); p != nil {
						recovered, isError := p.(error)
						if len(recovererInfo) > 0 {
							if isError {
								recovererInfo[0](recovered, debug.Stack())
							} else {
								recovererInfo[0](recovered)
							}
						}
					}
				}()

				if result.stop.Load() == 1 {
					return
				}

				callback(v)
			}()
		},
	)

	return &result
}

// Pub publishes an event which will be dispatched to all listeners
func Pub[T any](event T, failures ...func(err string)) {
	key := getKey[T]()
	existingTopic, ok := mapper.Load(key)
	if !ok || existingTopic == nil { // create a new topic, even if there are no listeners (otherwise we will have to panic)
		newTopic := NewTopic[T]()
		existingTopic, _ = mapper.LoadOrStore(key, newTopic)
	}

	switch topic := existingTopic.(type) {
	case *Topic[T]:
		topic.Pub(event)
	default:
		if len(failures) > 0 {
			failures[0](fmt.Sprintf("failt to publish on topic typed %T and event typed %T\n", topic, event))
		}
	}
}

// Range gives access to mapper Range
func Range(iter func(k, v any) bool) {
	mapper.Range(iter)
}

// Reset allows usage in tests that shares the bus
func Reset() {
	mapper = sync.Map{}
}

func getKey[T any]() string {
	var zero T
	if e, ok := any(zero).(iEventName); ok {
		return e.EventID()
	}

	return fmt.Sprintf("%T", zero)
}
