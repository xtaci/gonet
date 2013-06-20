package core

import (
	"sync"
)

import (
	"db"
	"misc/timer"
)

type Event struct {
	_type   int16
	user_id int32
	timeout int64
}

var (
	_event_ch    chan int32
	_events      map[int32]*Event // mapping from  event_id-> Event
	_events_lock sync.RWMutex
)

const (
	EVENT_CHAN_MAX = 200000
	EVENTID_GEN    = "EVENTID_GEN"
)

func init() {
	_event_ch = make(chan int32, EVENT_CHAN_MAX)
	_events = make(map[int32]*Event)
	go _expire()
}

func _expire() {
	for {
		event_id := <-_event_ch
		_events_lock.RLock()
		event := _events[event_id]
		_events_lock.RUnlock()

		// process event, sequentially
		if Execute(event, event_id) {
			_events_lock.Lock()
			delete(_events, event_id)
			_events_lock.Unlock()
		}
	}
}

//---------------------------------------------------------- Add a timeout event
func Add(Type int16, user_id int32, timeout int64) int32 {
	event_id := db.NextVal(EVENTID_GEN)
	timer.Add(event_id, timeout, _event_ch)

	event := &Event{_type: Type, user_id: user_id, timeout: timeout}
	_events_lock.Lock()
	_events[event_id] = event
	_events_lock.Unlock()

	return event_id
}

//---------------------------------------------------------- cancel an oid's timeout
func Cancel(event_id int32) {
	_events_lock.Lock()
	delete(_events, event_id)
	_events_lock.Unlock()
}

//---------------------------------------------------------- Load a timeout event at startup
func Load(Type int16, user_id int32, timeout int64, event_id int32) {
	timer.Add(event_id, timeout, _event_ch)
	event := &Event{_type: Type, user_id: user_id, timeout: timeout}
	_events[event_id] = event

	return
}
