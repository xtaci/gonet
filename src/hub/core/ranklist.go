package core

import (
	"sync"
)

import (
	"misc/alg/dos"
	. "types"
)

var (
	_ranklist      dos.Tree // dynamic order statistics
	_id_score      map[int32]int32
	_lock_ranklist sync.RWMutex
)

func init() {
	_id_score = make(map[int32]int32)
}

//------------------------------------------------ add a user to rank list
func _add_rank(user *User) {
	_lock_ranklist.Lock()
	defer _lock_ranklist.Unlock()
	_ranklist.Insert(int(user.Score), user.Id)
	_id_score[user.Id] = user.Score
}

//------------------------------------------------ update score of a player
func UpdateScore(id, oldscore, newscore int32) bool {
	_lock_ranklist.Lock()
	defer _lock_ranklist.Unlock()

	tmplist := make([]interface{}, 0, 64)

	defer func() {
		for i := range tmplist {
			_ranklist.Insert(int(oldscore), tmplist[i])
		}
	}()

	for {
		n, _ := _ranklist.ByScore(int(oldscore))

		if n == nil {
			return false
		}

		if n.Data().(int32) == id {
			_ranklist.DeleteNode(n)
			_ranklist.Insert(int(newscore), id)
			_id_score[id] = newscore
			return true
		} else {
			// temporary delete
			_ranklist.DeleteNode(n)
			tmplist = append(tmplist, n.Data())
		}
	}

	return true
}

//------------------------------------------------ get players from ranklist within [A,B]
func GetList(A, B int) (id []int32, score []int32) {
	if A <= 0 {
		A = 1
	}

	_lock_ranklist.Lock()
	defer _lock_ranklist.Unlock()

	if A > _ranklist.Count() {
		return
	}

	if B > _ranklist.Count() || B < A {
		B = _ranklist.Count()
	}

	id = make([]int32, B-A+1)
	score = make([]int32, B-A+1)

	for i := A; i <= B; i++ {
		n := _ranklist.Rank(i)
		id[i-A] = n.Data().(int32)
		score[i-A] = int32(n.Score())
	}

	return
}

//------------------------------------------------ get score
func Score(id int32) (ret int32) {
	_lock_ranklist.RLock()
	defer _lock_ranklist.RUnlock()
	return _id_score[id]
}
