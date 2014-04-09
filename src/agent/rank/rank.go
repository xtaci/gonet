package rank

import (
	"sync"
)

import (
	"db/user_tbl"
	. "helper"
	"misc/alg/dos"
)

var (
	_ranklist      dos.Tree // dynamic order statistics
	_id_cup        map[int32]int32
	_lock_ranklist sync.RWMutex
)

func init() {
	INFO("Loading Dynamic Order Statstistics.")
	_id_cup = make(map[int32]int32)

	go _init_ranks()
}

//---------------------------------------------------------- 异步加载排行榜
func _init_ranks() {
	// load users
	uds := user_tbl.GetAll()
	for i := range uds {
		Update(uds[i].Id, uds[i].Score)
	}
	INFO("Dynamic Order Statstistics Load Complete.")
}

//---------------------------------------------------------- update cup of a player
func Update(id, newcup int32) bool {
	_lock_ranklist.Lock()
	defer _lock_ranklist.Unlock()

	oldcup, ok := _id_cup[id]
	if !ok { // new user
		_ranklist.Insert(newcup, id)
		_id_cup[id] = newcup
		return true
	} else { // old user
		_, n := _ranklist.Locate(oldcup, id)
		if n == nil {
			ERR("没有在DOS中查到玩家", id)
			return false
		}

		_ranklist.DeleteNode(n)
		_ranklist.Insert(newcup, id)
		_id_cup[id] = newcup
		return true
	}
}

//---------------------------------------------------------- get players from ranklist within [A,B]
func GetList(A, B int) (id []int32, cup []int32) {
	if A < 1 || A > B {
		return
	}

	_lock_ranklist.RLock()
	defer _lock_ranklist.RUnlock()

	if A > _ranklist.Count() {
		return
	}

	if B > _ranklist.Count() {
		B = _ranklist.Count()
	}

	id, cup = make([]int32, B-A+1), make([]int32, B-A+1)
	for i := A; i <= B; i++ {
		n := _ranklist.Rank(i)
		id[i-A] = n.Id()
		cup[i-A] = n.Score()
	}

	return
}

//---------------------------------------------------------- get user of rank n
func RankN(n int32) int32 {
	_lock_ranklist.RLock()
	defer _lock_ranklist.RUnlock()
	node := _ranklist.Rank(int(n))
	if node != nil {
		return node.Id()
	}
	return -1
}

//---------------------------------------------------------- get rank for a user
func Rank(id int32) int32 {
	_lock_ranklist.Lock()
	defer _lock_ranklist.Unlock()
	rank, _ := _ranklist.Locate(_id_cup[id], id)
	return int32(rank)
}
