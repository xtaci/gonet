package player

import "sync"

var names map[int]chan string
var _lock sync.RWMutex

func RegisterChannel(ch chan string, id int) {
	_lock.Lock()
	names[id] = ch
	_lock.Unlock()
}

func QueryChannel(id int) chan string {
	var ch chan string

	_lock.RLock()
	ch = names[id]
	_lock.RUnlock()

	return ch
}

func InitNames() {
	names = make(map[int]chan string)
}
