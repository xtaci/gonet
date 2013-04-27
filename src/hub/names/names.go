package names

import "sync"

var names map[int32]chan interface{}
var _lock sync.RWMutex

func Register(ch chan interface{}, id int32) {
	defer _lock.Unlock()
	_lock.Lock()

	names[id] = ch
}

func Unregister(id int32) {
	defer _lock.Unlock()
	_lock.Lock()

	delete(names, id)
}

func Query(id int32) chan interface{} {
	defer _lock.RUnlock()
	_lock.RLock()

	return names[id]
}

func init() {
	names = make(map[int32]chan interface{})
}
