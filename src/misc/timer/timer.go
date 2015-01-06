package timer

import (
	"log"
	"time"
)

const (
	TIMER_LEVEL = 20 // 时间段最大分级，最大时间段为 2^TIMERLEVEL
	QUEUE_MAX   = 65536
)

type _timer_event struct {
	Id      int32      // 用户定义的ID
	Timeout int64      // 到期时间 Unix Time
	CH      chan int32 // 发送通道
}

var (
	_eventlist  [TIMER_LEVEL]map[uint32]_timer_event // 事件列表
	_eventqueue chan _timer_event                    // 事件添加队列
)

func init() {
	for k := range _eventlist {
		_eventlist[k] = make(map[uint32]_timer_event)
	}

	_eventqueue = make(chan _timer_event, QUEUE_MAX)

	go _timer()
}

//------------------------------------------------
// 定时器 goroutine
// 根据程序启动后经过的秒数计数
func _timer() {
	defer func() {
		if x := recover(); x != nil {
			log.Println("TIMER CRASHED", x)
		}
	}()

	last := time.Now().Unix()
	timer_id := uint32(0) // 内部事件编号

	sleep_timer := time.NewTimer(time.Second)
	for {
		select {
		case new_event := <-_eventqueue:
			timer_id++
			// 最小的时间间隔，处理为1s
			diff := new_event.Timeout - time.Now().Unix()
			if diff <= 0 {
				diff = 1
			}

			// 发到合适的框
			for i := TIMER_LEVEL - 1; i >= 0; i-- {
				if diff >= 1<<uint(i) {
					_eventlist[i][timer_id] = new_event
					break
				}
			}
		case <-sleep_timer.C:
			// 重置定时器
			sleep_timer.Reset(time.Second)

			// 检查事件触发
			// 累计距离上一次触发的秒数,并逐秒触发
			// 如果校正了系统时间，时间前移，nsec为负数的时候，last的值不应该变动，否则会出现秒数的重复计数
			now := time.Now().Unix()
			if now-last < 1 {
				continue
			}

			// 开始逐秒触发
			for {
				last++
				for i := TIMER_LEVEL - 1; i >= 0; i-- {
					mask := (1 << uint(i)) - 1
					if last&int64(mask) == 0 {
						_trigger(i)
					}
				}

				// 如果到达当前时间，停止循环
				if last == now {
					break
				}
			}
		}
	}
}

//------------------------------------------------ 单级触发
func _trigger(level int) {
	now := time.Now().Unix()
	list := _eventlist[level]

	for k, v := range list {
		if v.Timeout-now < 1<<uint(level) {
			if level == 0 { // 触发事件
				select { // 发送必须为非阻塞, 传入的chan不能关闭
				case v.CH <- v.Id:
				default:
				}
			} else { // 移动到前一个更短间距的时间队列
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
	_eventqueue <- _timer_event{Id: id, CH: ch, Timeout: timeout}
}
