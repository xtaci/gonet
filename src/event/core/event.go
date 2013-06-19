package core

import (
	"sync"
)

import (
	"db"
	"misc/naming"
	"misc/timer"
)

type Event struct {
	tblname uint32
	user_id int32
	oid     uint32
	timeout int64
}

var (
	_event_ch    chan int32
	_events      map[int32]*Event // mapping from  event_id-> Event
	_events_lock sync.RWMutex

	_hashtbl map[uint32]string // mapping from hash(tblname)-> tblname
)

const (
	EVENT_CHAN_MAX = 200000
	EVENTID_GEN    = "EVENTID_GEN"
)

func init() {
	_event_ch = make(chan int32, EVENT_CHAN_MAX)
	_events = make(map[int32]*Event)
	_hashtbl = make(map[uint32]string)
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
func Add(tblname string, oid uint32, user_id int32, timeout int64) int32 {
	h_tblname := naming.FNV1a(tblname)
	_hashtbl[h_tblname] = tblname

	event_id := db.NextVal(EVENTID_GEN)
	timer.Add(event_id, timeout, _event_ch)

	event := &Event{tblname: h_tblname, user_id: user_id, oid: oid, timeout: timeout}
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
func Load(tblname string, oid uint32, user_id int32, timeout int64, event_id int32) {
	h_tblname := naming.FNV1a(tblname)
	_hashtbl[h_tblname] = tblname

	timer.Add(event_id, timeout, _event_ch)

	event := &Event{tblname: h_tblname, user_id: user_id, oid: oid, timeout: timeout}
	_events[event_id] = event

	return
}
