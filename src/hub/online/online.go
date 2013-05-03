package online

import "sync"
import . "types"

var _active map[int32]*Session
var _lock sync.RWMutex

func Register(sess *Session, id int32) {
	defer _lock.Unlock()
	_lock.Lock()

	_active[id] = sess
}

func Unregister(id int32) {
	defer _lock.Unlock()
	_lock.Lock()

	delete(_active, id)
}

func Query(id int32) *Session {
	defer _lock.RUnlock()
	_lock.RLock()

	return _active[id]
}

func init() {
	_active = make(map[int32]*Session)
}
