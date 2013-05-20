package timer

import (
	"sync"
	"sync/atomic"
	"time"
)

type Event struct {
	Timeout int64       // timeout
	CH      chan uint32 // event trigger channel
}

const (
	TIMER_LEVEL = uint(16) // num of time intervals, 10 means max 2^10 seconds
)

var (
	_incr      uint32                         // event id generator
	_eventlist [TIMER_LEVEL]map[uint32]*Event // 2^n based time interval

	_eventqueue      map[uint32]*Event // add queue
	_eventqueue_lock sync.Mutex

	_cancelqueue      []uint32 // cancel queue
	_cancelqueue_lock sync.Mutex
)

func init() {
	for k := range _eventlist {
		_eventlist[k] = make(map[uint32]*Event)
	}

	_eventqueue = make(map[uint32]*Event)

	go _timer()
}

//------------------------------------------------ Timer Routine
func _timer() {
	timer_count := uint32(0)

	for {
		time.Sleep(time.Second)
		timer_count++

		// add pending events
		now := time.Now().Unix()
		_eventqueue_lock.Lock()
		for k, v := range _eventqueue {
			diff := v.Timeout - now
			if diff <= 0 { // in case of very near event
				diff = 1
			}

			for i := TIMER_LEVEL - 1; i >= 0; i-- {
				if diff >= 1<<i {
					_eventlist[i][k] = v
					break
				}
			}
		}
		_eventqueue = make(map[uint32]*Event)
		_eventqueue_lock.Unlock()

		// cancelqueue
		_cancelqueue_lock.Lock()
		for _, v := range _cancelqueue {
			for i := TIMER_LEVEL - 1; i >= 0; i-- {
				list := _eventlist[i]
				delete(list, v)
			}
		}
		_cancelqueue = nil
		_cancelqueue_lock.Unlock()

		// triggers
		for i := TIMER_LEVEL - 1; i > 0; i-- {
			mask := (uint32(1) << i) - 1
			if timer_count&mask == 0 {
				_trigger(i)
			}
		}

		_trigger(0)
	}
}

func _trigger(level uint) {
	now := time.Now().Unix()
	list := _eventlist[level]

	for k, v := range list {
		if v.Timeout-now < 1<<level {
			// move to one closer timer or just removal
			if level == 0 {
				v.CH <- k
			} else {
				_eventlist[level-1][k] = v
			}

			delete(list, k)
		}
	}
}

//------------------------------------------------ add a timeout event
func Add(timeout int64, ch chan uint32) uint32 {
	event_id := atomic.AddUint32(&_incr, 1)
	event := &Event{CH: ch, Timeout: timeout}

	_eventqueue_lock.Lock()
	_eventqueue[event_id] = event
	_eventqueue_lock.Unlock()
	return event_id
}

//------------------------------------------------ cancel an event
func Cancel(event_id uint32) {
	_cancelqueue_lock.Lock()
	_cancelqueue = append(_cancelqueue, event_id)
	_cancelqueue_lock.Unlock()
}
