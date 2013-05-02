package names

import "sync"
import . "types"

var names map[int32]*Session
var _lock sync.RWMutex

func Register(sess *Session, id int32) {
	defer _lock.Unlock()
	_lock.Lock()

	names[id] = sess
}

func Unregister(id int32) {
	defer _lock.Unlock()
	_lock.Lock()

	delete(names, id)
}

func Query(id int32) *Session {
	defer _lock.RUnlock()
	_lock.RLock()

	return names[id]
}

func init() {
	names = make(map[int32]*Session)
}
