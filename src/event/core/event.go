package core

import (
	"misc/naming"
	"misc/timer"
	"sync"
)

type Event struct {
	tblname uint32
	user_id int32
	oid     uint32
	timeout int64
}

var (
	_event_ch     chan uint32
	_events       map[uint32]*Event // mapping from  event_id-> Event
	_events_lock  sync.RWMutex
	_event_id_gen uint32

	_hashtbl map[uint32]string // mapping from hash(tblname)-> tblname
)

const (
	EVENT_CHAN_MAX = 200000
)

func init() {
	_event_ch = make(chan uint32, EVENT_CHAN_MAX)
	_events = make(map[uint32]*Event)
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

//------------------------------------------------ Add a timeout event
func Add(tblname string, oid uint32, user_id int32, timeout int64) uint32 {
	h_tblname := naming.FNV1a(tblname)
	_hashtbl[h_tblname] = tblname

	_event_id_gen++
	event_id := _event_id_gen
	timer.Add(event_id, timeout, _event_ch)

	event := &Event{tblname: h_tblname, user_id: user_id, oid: oid, timeout: timeout}
	_events_lock.Lock()
	_events[event_id] = event
	_events_lock.Unlock()

	return event_id
}

//------------------------------------------------ Load a timeout event
func Load(tblname string, oid uint32, user_id int32, timeout int64, event_id uint32) {
	h_tblname := naming.FNV1a(tblname)
	_hashtbl[h_tblname] = tblname

	timer.Add(event_id, timeout, _event_ch)

	event := &Event{tblname: h_tblname, user_id: user_id, oid: oid, timeout: timeout}
	_events[event_id] = event

	return
}

//------------------------------------------------ cancel an oid's timeout
func Cancel(event_id uint32) {
	_events_lock.Lock()
	delete(_events, event_id)
	_events_lock.Unlock()
}
