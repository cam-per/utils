package event

import (
	"github.com/cam-per/kuroko/pkg/generic"
	"sync"
	"time"
)

type Event interface {
	Close()
}

type eventInstance[T any] struct {
	parent  *Emitter[T]
	index   int
	handler func(event Event, data T)
}

func (e *eventInstance[T]) Close() {
	e.parent.close(e.index)
}

type Emitter[T any] struct {
	mu     sync.RWMutex
	events []*eventInstance[T]
	ttl    generic.Set[string]
}

func NewEmitter[T any]() *Emitter[T] {
	return &Emitter[T]{
		ttl: generic.NewSet[string](),
	}
}

func (emitter *Emitter[T]) Register(handler func(event Event, data T)) Event {
	emitter.mu.Lock()
	defer emitter.mu.Unlock()

	event := &eventInstance[T]{
		parent:  emitter,
		handler: handler,
		index:   len(emitter.events),
	}
	emitter.events = append(emitter.events, event)
	return event
}

func (emitter *Emitter[T]) close(index int) {
	emitter.mu.Lock()
	defer emitter.mu.Unlock()

	if len(emitter.events) == 1 {
		emitter.events = nil
		return
	}

	events := make([]*eventInstance[T], 0, len(emitter.events)-1)
	for i, v := range emitter.events {
		if i == index {
			emitter.events[i] = nil
			continue
		}
		events = append(events, v)
	}
	emitter.events = events
}

func (emitter *Emitter[T]) Emit(data T) {
	emitter.mu.RLock()
	defer emitter.mu.RUnlock()

	var wg sync.WaitGroup
	wg.Add(len(emitter.events))
	for _, event := range emitter.events {
		go func(e *eventInstance[T]) {
			defer wg.Done()
			e.handler(e, data)
		}(event)
	}
}

func (emitter *Emitter[T]) EmitAsync(data T) {
	emitter.mu.RLock()
	defer emitter.mu.RUnlock()

	for _, e := range emitter.events {
		go e.handler(e, data)
	}
}

func (emitter *Emitter[T]) EmitTTL(data T, hash string, ttl time.Duration) {
	emitter.mu.RLock()
	defer emitter.mu.RUnlock()

	if emitter.ttl.Has(hash) {
		return
	}
	emitter.ttl.Add(hash)
	defer time.AfterFunc(ttl, func() { emitter.ttl.Delete(hash) })

	var wg sync.WaitGroup
	wg.Add(len(emitter.events))
	for _, event := range emitter.events {
		go func(e *eventInstance[T]) {
			defer wg.Done()
			e.handler(e, data)
		}(event)
	}
}

func (emitter *Emitter[T]) EmitAsyncTTL(data T, hash string, ttl time.Duration) {
	emitter.mu.RLock()
	defer emitter.mu.RUnlock()

	if emitter.ttl.Has(hash) {
		return
	}
	emitter.ttl.Add(hash)
	defer time.AfterFunc(ttl, func() { emitter.ttl.Delete(hash) })

	for _, e := range emitter.events {
		go e.handler(e, data)
	}
}
