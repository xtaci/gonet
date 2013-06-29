package core

import (
	"sync"
	"time"
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
	EVENT_CHAN_MAX   = 200000
	EVENTID_GEN      = "EVENTID_GEN"
	CLEANUP_INTERVAL = 3600
)

func init() {
	_event_ch = make(chan int32, EVENT_CHAN_MAX)
	_events = make(map[int32]*Event)
	go _expire()
	go _cleanup()
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

//---------------------------------------------------------- 周期性的清空数据库中已发生的事件 
func _cleanup() {
	for {
		time.Sleep(time.Second * CLEANUP_INTERVAL)
		event_tbl.RemoveOld(time.Now().Unix() - CLEANUP_INTERVAL)
	}
}

//---------------------------------------------------------- Add a timeout event, return event id
func Add(Type int16, user_id int32, timeout int64, params []byte) int32 {
	event_id := db.NextVal(EVENTID_GEN)

	// first add to hashmap
	event := &Event{Type: Type, UserId: user_id, EventId: event_id, Timeout: timeout, Params: params}
	_events_lock.Lock()
	_events[event_id] = event
	_events_lock.Unlock()

	// then store to db
	event_tbl.Add(event)

	// and finally, put in timer
	timer.Add(event_id, timeout, _event_ch)

	return event_id
}

//---------------------------------------------------------- cancel an event with id
func Cancel(event_id int32) {
	_events_lock.Lock()
	delete(_events, event_id)
	_events_lock.Unlock()
}

//---------------------------------------------------------- load a event at startup
func Load(ev *Event) {
	_events[ev.EventId] = ev
	timer.Add(ev.EventId, ev.Timeout, _event_ch)
}
