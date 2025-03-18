package pubsub

import (
	"sync"
)

type PubSub[T any] struct {
	subs       map[chan T]bool
	bufferSize int
	closed     bool

	m sync.RWMutex
}

func NewPubSub[T any](bufferSize int) *PubSub[T] {
	return &PubSub[T]{
		subs:       make(map[chan T]bool),
		bufferSize: bufferSize,
	}
}

func (p *PubSub[T]) Subscribe() chan T {
	p.m.Lock()
	defer p.m.Unlock()

	ch := make(chan T, p.bufferSize)
	p.subs[ch] = true
	return ch
}

func (p *PubSub[T]) Publish(msg T) {
	p.m.RLock()
	defer p.m.RUnlock()

	if p.closed {
		return
	}
	for ch := range p.subs {
		ch <- msg
	}
}

func (p *PubSub[T]) PublishNoWait(msg T) {
	p.m.RLock()
	defer p.m.RUnlock()

	if p.closed {
		return
	}

	for ch := range p.subs {
		select {
		case ch <- msg:
		default:
		}
	}
}

func (p *PubSub[T]) Unsubscribe(ch chan T) {
	p.m.Lock()
	defer p.m.Unlock()

	if _, ok := p.subs[ch]; !ok {
		return
	}

	delete(p.subs, ch)
}

func (p *PubSub[T]) Close() {
	p.m.Lock()
	defer p.m.Unlock()

	if !p.closed {
		p.closed = true
		for ch := range p.subs {
			delete(p.subs, ch)
			close(ch)
		}
	}
}
