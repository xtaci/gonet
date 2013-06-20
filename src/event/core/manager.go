package core

import (
	"sync"
)

import (
	"db"
	"db/event_tbl"
	"misc/timer"
	. "types"
)

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
func Add(Type int16, user_id int32, timeout int64, params []byte) int32 {
	event_id := db.NextVal(EVENTID_GEN)
	timer.Add(event_id, timeout, _event_ch)

	event := &Event{Type: Type, UserId: user_id, EventId: event_id, Timeout: timeout, Params: params}
	_events_lock.Lock()
	_events[event_id] = event
	_events_lock.Unlock()

	// store to db
	event_tbl.Add(event)

	return event_id
}

//---------------------------------------------------------- cancel an oid's timeout
func Cancel(event_id int32) {
	_events_lock.Lock()
	delete(_events, event_id)
	_events_lock.Unlock()
}

//---------------------------------------------------------- Load a timeout event at startup
func Load(event_id int32, Type int16, user_id int32, timeout int64, params []byte) {
	timer.Add(event_id, timeout, _event_ch)
	event := &Event{EventId: event_id, Type: Type, UserId: user_id, Timeout: timeout, Params: params}
	_events[event_id] = event

	return
}
