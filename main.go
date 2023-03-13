package bus

import (
	"sync"
	"sync/atomic"
)

var mapper sync.Map // holds key (event id, typed string) versus topic values

type iEvent interface {
	EventID() string //
	Async() bool     // if returns true, this event will be triggered by spinning a goroutine
}

type Listener[T iEvent] struct {
	parent   *Topic[T]     // so we can call unsubscribe from parent
	callback func(event T) // the function that we're going to call
}

type Topic[T iEvent] struct {
	subs      []*Listener[T] // list of listeners
	rwMu      sync.RWMutex   // guards subs
	lisnsPool sync.Pool      // a pool of listeners
}

func NewTopic[T iEvent]() *Topic[T] {
	result := &Topic[T]{}
	result.lisnsPool.New = func() any {
		return &Listener[T]{
			parent: result,
		}
	}
	return result
}

func (t *Topic[T]) Sub(callback func(v T)) *Listener[T] {
	result := t.lisnsPool.Get().(*Listener[T])
	result.callback = callback
	result.parent = t

	t.rwMu.Lock()
	t.subs = append(t.subs, result)
	t.rwMu.Unlock()

	return result
}

func (t *Topic[T]) Forget(who *Listener[T]) {
	t.rwMu.Lock()
	for i := range t.subs {
		if t.subs[i] != who {
			continue
		}

		t.subs[i] = t.subs[len(t.subs)-1]
		t.subs[len(t.subs)-1] = nil
		t.subs = t.subs[:len(t.subs)-1]
		break
	}
	t.rwMu.Unlock()

	who.callback = nil
	t.lisnsPool.Put(who)
}

func (t *Topic[T]) NumSubs() uint64 {
	t.rwMu.RLock()
	result := uint64(len(t.subs))
	t.rwMu.RUnlock()
	return result
}

func (s *Listener[T]) Unsub() {
	s.parent.Forget(s)
}

func (s *Listener[T]) Topic() *Topic[T] {
	return s.parent
}

func (t *Topic[T]) Pub(event T) {
	t.rwMu.RLock()
	for topic := range t.subs {
		if event.Async() {
			go t.subs[topic].callback(event)
			continue
		}

		t.subs[topic].callback(event)
	}
	t.rwMu.RUnlock()
}

type Bus[T iEvent] struct {
	name     string
	listener *Listener[T]
	stop     atomic.Uint32 // flag for unsubscribing after receiving one event
}

func (bs *Bus[T]) Unsub() {
	if bs.stop.CompareAndSwap(0, 1) {
		go bs.listener.Unsub()
	}
}

func (bs *Bus[T]) String() string {
	return "bus for topic `" + string(bs.name) + "`"
}

func Listen[T iEvent](callback func(event T) bool) *Bus[T] {
	var event T
	topic, ok := mapper.Load(event.EventID())
	if !ok || topic == nil {
		topic, _ = mapper.LoadOrStore(event.EventID(), NewTopic[T]())
	}

	var result Bus[T]
	result.name = event.EventID()

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

func Publish[T iEvent](event T) {
	topic, ok := mapper.Load(event.EventID())
	if !ok || topic == nil { // create new topic, even if there are no listeners (otherwise we will have to panic)
		topic, _ = mapper.LoadOrStore(event.EventID(), NewTopic[T]())
	}
	topic.(*Topic[T]).Pub(event)
}
