package timer

import (
	"sync"
	"sync/atomic"
	"time"
)

type Event struct {
	Id      int32      // user specified Id
	Timeout int64      // timeout
	CH      chan int32 // event trigger channel
}

const (
	TIMER_LEVEL = uint(16) // num of time intervals, 10 means max 2^10 seconds
)

var (
	_eventlist [TIMER_LEVEL]map[uint32]*Event // 2^n based time interval

	_eventqueue      map[uint32]*Event // add queue
	_eventqueue_lock sync.Mutex

	_cancelqueue      []uint32 // cancel queue
	_cancelqueue_lock sync.Mutex

	_timer_id uint32	// 内部事件编号
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
	last := time.Now().Unix()

	for {
		// add pending events
		_eventqueue_lock.Lock()
		for k, v := range _eventqueue {
			diff := v.Timeout - time.Now().Unix()
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
		now := time.Now().Unix()
		nsecs := now - last // num of seconds passed since last trigger
		last = now

		if nsecs < 0 { // maybe someone changed system clock
			continue
		}

		for c := int64(0); c < nsecs; c++ {
			timer_count++

			for i := TIMER_LEVEL - 1; i > 0; i-- {
				mask := (uint32(1) << i) - 1
				if timer_count&mask == 0 {
					_trigger(i)
				}
			}

			_trigger(0)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func _trigger(level uint) {
	now := time.Now().Unix()
	list := _eventlist[level]

	for k, v := range list {
		if v.Timeout-now < 1<<level {
			// move to one closer timer-list or trigger
			if level == 0 {
				func() {
					defer func() {
						recover() // ignore closed channel
					}()

					v.CH <- v.Id
				}()
			} else {
				_eventlist[level-1][k] = v
			}

			delete(list, k)
		}
	}
}

//------------------------------------------------ 
// 添加一个定时，timeout为到期的Unix时间
// id 是调用者定义的编号, 事件发生时，会把id发送到ch
func Add(id int32, timeout int64, ch chan int32) {
	event := &Event{Id: id, CH: ch, Timeout: timeout}

	timer_id := atomic.AddUint32(&_timer_id, 1)
	_eventqueue_lock.Lock()
	_eventqueue[timer_id] = event
	_eventqueue_lock.Unlock()
}
