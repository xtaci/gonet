package player

import "sync"

var names map[string]chan string
var _lock sync.RWMutex

func RegisterChannel(ch chan string, name string) {
	_lock.Lock()
	names[name] = ch
	_lock.Unlock()
}

func QueryChannel(name string) chan string {
	var ch chan string

	_lock.RLock()
	ch = names[name]
	_lock.RUnlock()

	return ch
}

func InitNames() {
	names = make(map[string]chan string)
}
