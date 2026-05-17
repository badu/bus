package bus_test

import (
	"sync/atomic"
	"testing"

	"github.com/badu/bus"
)

type Uint32SyncEvent struct {
	u uint32
}

type Uint32Event struct {
	u uint32
}

func BenchmarkBroadcast_0008(b *testing.B) {
	topic := bus.NewTopic[Uint32SyncEvent]()
	c := uint32(0)
	for i := 0; i < 8; i++ {
		topic.Sub(func(v Uint32SyncEvent) {
			atomic.AddUint32(&c, v.u)
		})
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			topic.Pub(Uint32SyncEvent{u: 1})
		}
	})
}

func BenchmarkBroadcast_0008Ptr(b *testing.B) {
	topic := bus.NewTopic[*Uint32SyncEvent]()
	c := uint32(0)
	for i := 0; i < 8; i++ {
		topic.Sub(func(v *Uint32SyncEvent) {
			atomic.AddUint32(&c, v.u)
		})
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			topic.Pub(&Uint32SyncEvent{u: 1})
		}
	})
}

func BenchmarkBroadcast_0256(b *testing.B) {
	topic := bus.NewTopic[Uint32SyncEvent]()
	var c uint32
	for i := 0; i < 256; i++ {
		topic.Sub(func(v Uint32SyncEvent) {
			atomic.AddUint32(&c, v.u)
		})
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			topic.Pub(Uint32SyncEvent{u: 1})
		}
	})
}

func BenchmarkBroadcast_0256Ptr(b *testing.B) {
	topic := bus.NewTopic[*Uint32SyncEvent]()
	var c uint32
	for i := 0; i < 256; i++ {
		topic.Sub(func(v *Uint32SyncEvent) {
			atomic.AddUint32(&c, v.u)
		})
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			topic.Pub(&Uint32SyncEvent{u: 1})
		}
	})
}

func BenchmarkBroadcast_1k(b *testing.B) {
	topic := bus.NewTopic[Uint32SyncEvent]()
	var c uint32
	for i := 0; i < 1024; i++ {
		topic.Sub(func(v Uint32SyncEvent) {
			atomic.AddUint32(&c, v.u)
		})
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			topic.Pub(Uint32SyncEvent{u: 1})
		}
	})
}

func BenchmarkBroadcast_1kPtr(b *testing.B) {
	topic := bus.NewTopic[*Uint32SyncEvent]()
	var c uint32
	for i := 0; i < 1024; i++ {
		topic.Sub(func(v *Uint32SyncEvent) {
			atomic.AddUint32(&c, v.u)
		})
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			topic.Pub(&Uint32SyncEvent{u: 1})
		}
	})
}

func BenchmarkBroadcast_2k(b *testing.B) {
	topic := bus.NewTopic[Uint32SyncEvent]()
	var c uint32
	for i := 0; i < 2048; i++ {
		topic.Sub(func(v Uint32SyncEvent) {
			atomic.AddUint32(&c, v.u)
		})
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			topic.Pub(Uint32SyncEvent{u: 1})
		}
	})
}

func BenchmarkBroadcast_2kPtr(b *testing.B) {
	topic := bus.NewTopic[*Uint32SyncEvent]()
	var c uint32
	for i := 0; i < 2048; i++ {
		topic.Sub(func(v *Uint32SyncEvent) {
			atomic.AddUint32(&c, v.u)
		})
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			topic.Pub(&Uint32SyncEvent{u: 1})
		}
	})
}
