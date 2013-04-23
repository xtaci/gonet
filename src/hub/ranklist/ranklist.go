package ranklist

import (
	"sync"
	"errors"
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

//--------------------------------------------------------- change score of user
func ChangeScore(id, oldscore, newscore int) (err error){
	_lock.Lock()
	defer _lock.Unlock()

	var idlist []int
	defer func() {
		for i := range idlist {
			AddUser(idlist[i], oldscore)
		}
	}()

	for {
		n,_ := _ranklist.Score(oldscore)

		if n==nil {
			err = errors.New("cannot change user with score")
			return
		}

		if n.Data().(int) == id {
			_ranklist.DeleteNode(n)
			_ranklist.Insert(newscore,id)
			return
		} else {
			// temporary delete 
			_ranklist.DeleteNode(n)
			idlist = append(idlist, n.Data().(int))
		}
	}

	return
}

//--------------------------------------------------------- find user rank
func Find(id, score int) (rank int, err error){
	_lock.Lock()
	defer _lock.Unlock()

	var idlist []int
	defer func() {
		for i := range idlist {
			AddUser(idlist[i], score)
		}
	}()

	for {
		n, r := _ranklist.Score(score)

		if n==nil {
			err = errors.New("find user with score")
			return
		}

		_ranklist.DeleteNode(n)
		if n.Data().(int) == id {
			rank = r
			return
		}
	}

	return
}


//--------------------------------------------------------- get users from ranklist in [from, to] 
func GetRankList(from, to int) []int {
	sublist := make([]int, to-from+1)

	_lock.RLock()
	defer _lock.RUnlock()

	for i:=from;i<=to;i++ {
		sublist[i-from] = _ranklist.Rank(i).Data().(int)
	}

	return sublist
}
