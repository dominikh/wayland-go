package wayland

import (
	"reflect"
	"sync"
)

type EventQueue struct {
	// signals the availability of events
	ch chan struct{}

	mu sync.Mutex
	events []event
}

func NewEventQueue() *EventQueue {
	return &EventQueue{
		ch:make(chan struct{}, 1),
	}
}

func (q *EventQueue) push(ev event) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.events=append(q.events, ev)
	select{
	case q.ch <- struct{}{}:
	default:
	}
}

func (q *EventQueue) Dispatch() {
	<-q.ch
	q.mu.Lock()
	defer q.mu.Unlock()
	for _, ev := range q.events {
		p := ev.obj.GetProxy()
		cb := p.eventHandlers[ev.ev]
		if cb != nil {
			args := []reflect.Value{reflect.ValueOf(ev.obj)}
			args = append(args, ev.args...)
			reflect.ValueOf(cb).Call(args)
		}
	}
	q.events=q.events[:0]
}
