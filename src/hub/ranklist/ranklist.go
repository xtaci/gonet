package ranklist

import (
	"errors"
	"sync"
	"sync/atomic"
)

import (
	"misc/alg/dos"
)

var _ranklist dos.Tree
var _lock sync.RWMutex

var _count int32

func Increase() int32 {
	return atomic.AddInt32(&_count, 1)
}

func Decrease() int32 {
	return atomic.AddInt32(&_count, -1)
}

//--------------------------------------------------------- add a user to rank list
func AddUser(id int32, score int) {
	_lock.Lock()
	defer _lock.Unlock()
	_ranklist.Insert(score, id)
	Increase()
}

//--------------------------------------------------------- change score of user
func ChangeScore(id int32, oldscore, newscore int) (err error) {
	_lock.Lock()
	defer _lock.Unlock()

	var idlist []int32
	defer func() {
		for i := range idlist {
			AddUser(idlist[i], oldscore)
		}
	}()

	for {
		n, _ := _ranklist.Score(oldscore)

		if n == nil {
			err = errors.New("cannot change user with score")
			return
		}

		if n.Data().(int32) == id {
			_ranklist.DeleteNode(n)
			_ranklist.Insert(newscore, id)
			return
		} else {
			// temporary delete 
			_ranklist.DeleteNode(n)
			idlist = append(idlist, n.Data().(int32))
		}
	}

	return
}

//--------------------------------------------------------- find user rank
func Find(id int32, score int) (rank int, err error) {
	_lock.Lock()
	defer _lock.Unlock()

	var idlist []int32
	defer func() {
		for i := range idlist {
			AddUser(idlist[i], score)
		}
	}()

	for {
		n, r := _ranklist.Score(score)

		if n == nil {
			err = errors.New("find user with score")
			return
		}

		_ranklist.DeleteNode(n)
		if n.Data().(int32) == id {
			rank = r
			return
		}
	}

	return
}

func Count() int {
	_lock.RLock()
	defer _lock.RUnlock()
	return _ranklist.Count()
}

//--------------------------------------------------------- get users from ranklist in [from, to] 
func GetRankList(from, to int) []int32 {
	sublist := make([]int32, to-from+1)

	_lock.RLock()
	defer _lock.RUnlock()

	for i := from; i <= to; i++ {
		sublist[i-from] = _ranklist.Rank(i).Data().(int32)
	}

	return sublist
}
