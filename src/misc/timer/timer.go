package production

import (
	"sync/atomic"
	"time"
)

type Item struct {
	Id uint32
	User int32
	Expire int32
	Finish int64
}

const (
	TIMER_LEVEL = uint(10)
)

var (
	_incr uint32
	_eventlist [TIMER_LEVEL]map[uint32]*Item
)

func init() {
	for k := range _eventlist {
		_eventlist[k] = make(map[uint32]*Item)
	}
}

func TimerRoutine() {
	timer_count := uint32(0)

	for {
		time.Sleep(time.Second)
		timer_count++

		_trigger(0)

		for i:=uint(1);i<TIMER_LEVEL;i++ {
			mask := (uint32(1) << i)-1
			if timer_count & mask == 0 {
				_trigger(i)
			}
		}
	}
}

func _trigger(level uint) {
	now := time.Now().Unix()
	list := _eventlist[level]

	for k,v := range list {
		if v.Finish - now <= 1 << level {
			// move to one closer timer or just removal
			if level == 0 {
				// trigger event
			} else {
				delete(list, k)
			}
		}
	}
}

//----------------------------------------------- add a timeout for user, 'finish' is unix time
func Add(id int32, finish int64) uint32 {
	//expire:= int32(finish - time.Now().Unix())
	timer_id := atomic.AddUint32(&_incr, 1)

	return timer_id
}

func Cancel(id int32, product int32) {
}
