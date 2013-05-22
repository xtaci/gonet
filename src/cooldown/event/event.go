package event

import (
	"misc/timer"
	"sync"
)

type Event struct {
	OID	int32
	user_id int32
	timeout int64
}

var (
	_event_ch chan uint32
	_events map[uint32]*Event
	_events_lock sync.Mutex
)

const (
	EVENT_MAX = 200000
)

func init() {
	_event_ch = make(chan uint32, EVENT_MAX)
	go _expire()
}

func _expire() {
	for {
		event_id := <-_event_ch

		_events_lock.Lock()
		event := _events[event_id]
		delete(_events, event_id)
		_events_lock.Unlock()

		// process event
		Execute(event)
	}
}

//------------------------------------------------ Add a timeout for a object-id
func Add(oid int32, user_id int32, timeout int64) uint32 {
	_event_id := timer.Add(timeout, _event_ch)
	event := &Event{OID:oid, user_id:user_id, timeout:timeout}

	_events_lock.Lock()
	_events[_event_id] = event
	_events_lock.Unlock()

	return _event_id
}

//------------------------------------------------ cancel an timeout
func Cancel(event_id uint32) {
	timer.Cancel(event_id)

	_events_lock.Lock()
	delete(_events, event_id)
	_events_lock.Unlock()
}
