package production

import (
	"sync/atomic"
	"sync"
	"time"
)

type Event struct {
	User int32
	Timeout int64
}

const (
	TIMER_LEVEL = uint(10)
)

var (
	_incr uint32
	_eventlist [TIMER_LEVEL]map[uint32]*Event

	_eventqueue map[uint32]*Event
	_eventqueue_lock sync.Mutex

	_cancelqueue []uint32
	_cancelqueue_lock sync.Mutex
)

func init() {
	for k := range _eventlist {
		_eventlist[k] = make(map[uint32]*Event)
	}

	_eventqueue = make(map[uint32]*Event)
}

func TimerRoutine() {
	timer_count := uint32(0)

	for {
		time.Sleep(time.Second)
		timer_count++

		now := time.Now().Unix()
		// add pending events
		_eventqueue_lock.Lock()
		for k,v := range _eventqueue {
			diff := v.Timeout - now
			if diff <= 0 {	// in case of very near event
				diff = 1
			}

			for i:= TIMER_LEVEL-1;i>=0;i-- {
				if diff >= 1 << i {
					_eventlist[i][k] = v
					break
				}
			}
		}
		_eventqueue = make(map[uint32]*Event)
		_eventqueue_lock.Unlock()

		// cancelqueue
		for _,v := range _cancelqueue {
			for i:= TIMER_LEVEL-1;i>=0;i-- {
				list := _eventlist[i]
				delete(list,v)
			}
		}

		// triggers
		for i:= TIMER_LEVEL-1;i>0;i-- {
			mask := (uint32(1) << i)-1
			if timer_count & mask == 0 {
				_trigger(i)
			}
		}

		_trigger(0)
	}
}

func _trigger(level uint) {
	now := time.Now().Unix()
	list := _eventlist[level]

	for k,v := range list {
		if v.Timeout - now < 1 << level {
			// move to one closer timer or just removal
			if level == 0 {
			} else {
				_eventlist[level-1][k] = v
				delete(list, k)
			}
		}
	}
}

//----------------------------------------------- add a timeout for user
func Add(user_id int32, timeout int64) uint32 {
	event_id:= atomic.AddUint32(&_incr, 1)
	event := &Event{User:user_id, Timeout:timeout}

	_eventqueue_lock.Lock()
	_eventqueue[event_id] = event
	_eventqueue_lock.Unlock()
	return event_id
}

func Cancel(event_id uint32) {
	_cancelqueue_lock.Lock()
	_cancelqueue = append(_cancelqueue, event_id)
	_cancelqueue_lock.Unlock()
}
