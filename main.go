package bus

import (
	"sync"
	"sync/atomic"
)

var mapper sync.Map // holds key (event id, typed string) versus topic values

// internal interface that all the events must implement
type iEvent interface {
	EventID() string //
	Async() bool     // if returns true, this event will be triggered by spinning a goroutine
}

// Listener is being returned when you subscribe to a topic, so you can unsubscribe or access the parent topic
type Listener[T iEvent] struct {
	parent   *Topic[T]     // so we can call unsubscribe from parent
	callback func(event T) // the function that we're going to call
}

// Topic keeps the subscribers of one topic
type Topic[T iEvent] struct {
	subs      []*Listener[T] // list of listeners
	rwMu      sync.RWMutex   // guards subs
	lisnsPool sync.Pool      // a pool of listeners
}

// NewTopic creates a new topic for a specie of events
func NewTopic[T iEvent]() *Topic[T] {
	result := &Topic[T]{}
	result.lisnsPool.New = func() any {
		return &Listener[T]{
			parent: result,
		}
	}
	return result
}

// Sub adds a callback to be called when an event of that type is being published
func (b *Topic[T]) Sub(callback func(v T)) *Listener[T] {
	result := b.lisnsPool.Get().(*Listener[T])
	result.callback = callback
	result.parent = b

	b.rwMu.Lock()
	b.subs = append(b.subs, result)
	b.rwMu.Unlock()

	return result
}

// unsub is private to the topic, but can be accessed via Listener
func (b *Topic[T]) unsub(who *Listener[T]) {
	b.rwMu.Lock()
	for i := range b.subs {
		if b.subs[i] != who {
			continue
		}

		b.subs[i] = b.subs[len(b.subs)-1]
		b.subs[len(b.subs)-1] = nil
		b.subs = b.subs[:len(b.subs)-1]
		break
	}
	b.rwMu.Unlock()

	who.callback = nil
	b.lisnsPool.Put(who)
}

// NumSubs in case you need to perform tests and check the number of subscribers of this particular topic
func (b *Topic[T]) NumSubs() int {
	b.rwMu.RLock()
	result := len(b.subs)
	b.rwMu.RUnlock()
	return result
}

// Unsub forgets the indicated callback
func (s *Listener[T]) Unsub() {
	s.parent.unsub(s)
}

// Topic gives access to the underlying topic
func (s *Listener[T]) Topic() *Topic[T] {
	return s.parent
}

// Pub allows you to publish an event in that topic
func (b *Topic[T]) Pub(event T) {
	b.rwMu.RLock()
	for topic := range b.subs {
		if event.Async() {
			go b.subs[topic].callback(event)
			continue
		}

		b.subs[topic].callback(event)
	}
	b.rwMu.RUnlock()
}

// Bus is being returned when you subscribe, so you can manually Unsub
type Bus[T iEvent] struct {
	listener *Listener[T]
	stop     atomic.Uint32 // flag for unsubscribing after receiving one event
}

// Unsub allows caller to manually unsubscribe, in case they don't want to use SubUnsub
func (o *Bus[T]) Unsub() {
	if o.stop.CompareAndSwap(0, 1) {
		go o.listener.Unsub()
	}
}

// SubUnsub can be used if you need to unsubscribe immediately after receiving an event, by making your function return true
func SubUnsub[T iEvent](callback func(event T) bool) *Bus[T] {
	var event T
	topic, ok := mapper.Load(event.EventID())
	if !ok || topic == nil {
		topic, _ = mapper.LoadOrStore(event.EventID(), NewTopic[T]())
	}

	var result Bus[T]

	result.listener = topic.(*Topic[T]).Sub(func(v T) {
		if result.stop.Load() == 1 {
			return
		}

		unsub := callback(v)
		if unsub {
			result.Unsub()
		}

	})

	return &result
}

// Sub subscribes a callback function to listen for a specie of events
func Sub[T iEvent](callback func(event T)) *Bus[T] {
	var event T
	topic, ok := mapper.Load(event.EventID())
	if !ok || topic == nil {
		topic, _ = mapper.LoadOrStore(event.EventID(), NewTopic[T]())
	}

	var result Bus[T]

	result.listener = topic.(*Topic[T]).Sub(func(v T) {
		if result.stop.Load() == 1 {
			return
		}
		callback(v)
	})

	return &result
}

// Pub publishes an event which will be dispatched to all listeners
func Pub[T iEvent](event T) {
	topic, ok := mapper.Load(event.EventID())
	if !ok || topic == nil { // create new topic, even if there are no listeners (otherwise we will have to panic)
		topic, _ = mapper.LoadOrStore(event.EventID(), NewTopic[T]())
	}
	topic.(*Topic[T]).Pub(event)
}
