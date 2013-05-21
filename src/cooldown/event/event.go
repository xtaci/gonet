package event

import (
	"misc/timer"
	"sync"
)

type Event struct {
	user_id int32
	obj_id int32
	obj_type int32
	obj_nextlevel int32
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

		// TODO : process event
	}
}

func Add(user_id int32, obj_id int32, obj_type int32, obj_nextlevel int32, timeout int64) uint32 {
	_event_id := timer.Add(timeout, _event_ch)
	event := &Event{user_id:user_id, obj_id:obj_id, obj_type:obj_type, obj_nextlevel:obj_nextlevel, timeout:timeout}

	_events_lock.Lock()
	_events[_event_id] = event
	_events_lock.Unlock()

	return _event_id
}

func Cancel(event_id uint32) {
	timer.Cancel(event_id)

	_events_lock.Lock()
	delete(_events, event_id)
	_events_lock.Unlock()
}
