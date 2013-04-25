package ranklist

import (
	"sync"
)

import (
	"misc/alg/dos"
)

var _ranklist dos.Tree
var _lock sync.RWMutex

//--------------------------------------------------------- add a user to rank list
func AddUser(id, score int) {
	_lock.Lock()
	defer _lock.Unlock()
	_ranklist.Insert(score, id)
}

//--------------------------------------------------------- delete user on rank list
func DeleteUser(rank int) {
	_lock.Lock()
	defer _lock.Unlock()

	n := _ranklist.Get(rank)
	_ranklist.DeleteNode(n)
}

//--------------------------------------------------------- get users from ranklist in [from, to] 
func GetRank(from, to, int) []int {
	sublist := make([]int, to-from+1)

	_lock.RLock()
	defer _lock.RUnlock()

	for i:=from;i<=to;i++ {
		sublist[i-from] = _ranklist.Get(i)
	}
}
