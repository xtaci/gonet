package names

import "sync"

var names map[int]chan string
var _lock sync.RWMutex

func Register(ch chan string, id int) {
	_lock.Lock()
	names[id] = ch
	_lock.Unlock()
}

func Unregister(id int) {
	_lock.Lock()
	delete(names, id)
	_lock.Unlock()
}

func Query(id int) chan string {
	var ch chan string

	_lock.RLock()
	ch = names[id]
	_lock.RUnlock()

	return ch
}

func init() {
	names = make(map[int]chan string)
}
