package bus_test

import (
	"sync/atomic"
	"testing"

	"github.com/badu/bus"
)

type Uint32SyncEvent struct {
	u uint32
}

func (u Uint32SyncEvent) EventID() string {
	return "Uint32SyncEvent"
}

type Uint32AsyncEvent struct {
	u uint32
}

func (u Uint32AsyncEvent) EventID() string {
	return "Uint32AsyncEvent"
}

func BenchmarkBroadcast_0008Sync(b *testing.B) {
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

func BenchmarkBroadcast_0008Async(b *testing.B) {
	topic := bus.NewTopic[Uint32AsyncEvent]()
	c := uint32(0)
	for i := 0; i < 8; i++ {
		topic.Sub(func(v Uint32AsyncEvent) {
			atomic.AddUint32(&c, v.u)
		})
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			topic.PubAsync(Uint32AsyncEvent{u: 1})
		}
	})
}

func BenchmarkBroadcast_0008PtrSync(b *testing.B) {
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

func BenchmarkBroadcast_0008PtrAsync(b *testing.B) {
	topic := bus.NewTopic[*Uint32AsyncEvent]()
	c := uint32(0)
	for i := 0; i < 8; i++ {
		topic.Sub(func(v *Uint32AsyncEvent) {
			atomic.AddUint32(&c, v.u)
		})
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			topic.PubAsync(&Uint32AsyncEvent{u: 1})
		}
	})
}

func BenchmarkBroadcast_0256Sync(b *testing.B) {
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

func BenchmarkBroadcast_0256Async(b *testing.B) {
	topic := bus.NewTopic[Uint32AsyncEvent]()
	var c uint32
	for i := 0; i < 256; i++ {
		topic.Sub(func(v Uint32AsyncEvent) {
			atomic.AddUint32(&c, v.u)
		})
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			topic.PubAsync(Uint32AsyncEvent{u: 1})
		}
	})
}

func BenchmarkBroadcast_0256PtrSync(b *testing.B) {
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

func BenchmarkBroadcast_0256PtrAsync(b *testing.B) {
	topic := bus.NewTopic[*Uint32AsyncEvent]()
	var c uint32
	for i := 0; i < 256; i++ {
		topic.Sub(func(v *Uint32AsyncEvent) {
			atomic.AddUint32(&c, v.u)
		})
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			topic.PubAsync(&Uint32AsyncEvent{u: 1})
		}
	})
}

func BenchmarkBroadcast_1kSync(b *testing.B) {
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

func BenchmarkBroadcast_1kAsync(b *testing.B) {
	topic := bus.NewTopic[Uint32AsyncEvent]()
	var c uint32
	for i := 0; i < 1024; i++ {
		topic.Sub(func(v Uint32AsyncEvent) {
			atomic.AddUint32(&c, v.u)
		})
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			topic.PubAsync(Uint32AsyncEvent{u: 1})
		}
	})
}

func BenchmarkBroadcast_1kPtrSync(b *testing.B) {
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

func BenchmarkBroadcast_1kPtrAsync(b *testing.B) {
	topic := bus.NewTopic[*Uint32AsyncEvent]()
	var c uint32
	for i := 0; i < 1024; i++ {
		topic.Sub(func(v *Uint32AsyncEvent) {
			atomic.AddUint32(&c, v.u)
		})
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			topic.PubAsync(&Uint32AsyncEvent{u: 1})
		}
	})
}
func BenchmarkBroadcast_2kSync(b *testing.B) {
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

func BenchmarkBroadcast_2kAsync(b *testing.B) {
	topic := bus.NewTopic[Uint32AsyncEvent]()
	var c uint32
	for i := 0; i < 2048; i++ {
		topic.Sub(func(v Uint32AsyncEvent) {
			atomic.AddUint32(&c, v.u)
		})
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			topic.PubAsync(Uint32AsyncEvent{u: 1})
		}
	})
}

func BenchmarkBroadcast_2kPtrSync(b *testing.B) {
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

func BenchmarkBroadcast_2kPtrAsync(b *testing.B) {
	topic := bus.NewTopic[*Uint32AsyncEvent]()
	var c uint32
	for i := 0; i < 2048; i++ {
		topic.Sub(func(v *Uint32AsyncEvent) {
			atomic.AddUint32(&c, v.u)
		})
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			topic.PubAsync(&Uint32AsyncEvent{u: 1})
		}
	})
}
