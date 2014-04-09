package timer

import (
	"sync"
	"sync/atomic"
	"time"
)

type _timer_event struct {
	Id      int32      // 用户定义的ID
	Timeout int64      // 到期时间 Unix Time
	CH      chan int32 // 发送通道
}

const (
	TIMER_LEVEL = uint(16) // 时间段最大分级，最大时间段为 2^TIMERLEVEL
)

var (
	_eventlist [TIMER_LEVEL]map[uint32]*_timer_event // 事件列表

	_eventqueue      map[uint32]*_timer_event // 事件添加队列
	_eventqueue_lock sync.Mutex

	_timer_id uint32 // 内部事件编号
)

func init() {
	for k := range _eventlist {
		_eventlist[k] = make(map[uint32]*_timer_event)
	}

	_eventqueue = make(map[uint32]*_timer_event)

	go _timer()
}

//------------------------------------------------
// 定时器 goroutine
// 根据程序启动后经过的秒数计数
func _timer() {
	timer_count := uint32(0)
	last := time.Now().Unix()

	for {
		time.Sleep(100 * time.Millisecond)

		// 处理排队
		// 最小的时间间隔，处理为1s
		_eventqueue_lock.Lock()
		for k, v := range _eventqueue {
			// 处理微小间隔
			diff := v.Timeout - time.Now().Unix()
			if diff <= 0 {
				diff = 1
			}

			// 发到合适的框
			for i := TIMER_LEVEL - 1; i >= 0; i-- {
				if diff >= 1<<i {
					_eventlist[i][k] = v
					break
				}
			}

			// 从队列中删除这个事件
			delete(_eventqueue, k)
		}
		_eventqueue_lock.Unlock()

		// 检查事件触发
		// 累计距离上一次触发的秒数,并逐秒触发
		// 如果校正了系统时间，时间前移，nsec为负数的时候，last的值不应该变动，否则会出现秒数的重复计数
		now := time.Now().Unix()
		nsecs := now - last

		if nsecs <= 0 {
			continue
		} else {
			last = now
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
	}
}

//------------------------------------------------ 单级触发
func _trigger(level uint) {
	now := time.Now().Unix()
	list := _eventlist[level]

	for k, v := range list {
		if v.Timeout-now < 1<<level {
			// 移动到前一个更短间距的LIST
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
	event := &_timer_event{Id: id, CH: ch, Timeout: timeout}

	timer_id := atomic.AddUint32(&_timer_id, 1)
	_eventqueue_lock.Lock()
	_eventqueue[timer_id] = event
	_eventqueue_lock.Unlock()
}
