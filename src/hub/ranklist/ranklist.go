package ranklist

import (
	"errors"
	"sync"
	"sync/atomic"
)

import (
	"misc/alg/dos"
	. "types"
)

var _ranklist dos.Tree
var _lock sync.RWMutex

var _users map[int32]*User
var _userlock sync.RWMutex

var _count int32

func Increase() int32 {
	return atomic.AddInt32(&_count, 1)
}

func Decrease() int32 {
	return atomic.AddInt32(&_count, -1)
}

func init() {
	_users = make(map[int32]*User)
}

//--------------------------------------------------------- add a user to rank list, only useful when startup & register
func AddUser(ud *User, score int) {
	_lock.Lock()
	_userlock.Lock()

	defer func() {
		_lock.Unlock()
		_userlock.Unlock()
	}()

	_ranklist.Insert(score, ud)
	_users[ud.Id] = ud

	// atomic ops
	Increase()
}

//--------------------------------------------------------- change score of user
func ChangeScore(id, oldscore, newscore int32) (err error) {
	_lock.Lock()
	defer _lock.Unlock()

	var tmplist []interface{}
	defer func() {
		for i := range tmplist {
			_ranklist.Insert(int(oldscore), tmplist[i])
		}
	}()

	for {
		n, _ := _ranklist.Score(int(oldscore))

		if n == nil {
			err = errors.New("cannot change user with score")
			return
		}

		if n.Data().(*User).Id == id {
			_ranklist.DeleteNode(n)
			n.Data().(*User).Score = newscore
			_ranklist.Insert(int(newscore), n.Data().(*User))
			return
		} else {
			// temporary delete 
			_ranklist.DeleteNode(n)
			tmplist = append(tmplist, n.Data())
		}
	}

	return
}

//--------------------------------------------------------- find user rank, return a copy
func GetUserCopy(id int32) (ud User) {
	_userlock.RLock()
	defer _userlock.RUnlock()
	return *_users[id]
}

func Count() int {
	_lock.RLock()
	defer _lock.RUnlock()
	return _ranklist.Count()
}

//--------------------------------------------------------- get users from ranklist in [from, to] 
func GetRankList(from, to int) []*User {
	sublist := make([]*User, to-from+1)

	_lock.RLock()
	defer _lock.RUnlock()

	for i := from; i <= to; i++ {
		sublist[i-from] = _ranklist.Rank(i).Data().(*User)
	}

	return sublist
}

//--------------------------------------------------------- change user state
func ChangeState(id int32, oldstate, newstate int32) bool {
	_userlock.Lock()
	defer _userlock.Unlock()

	user := _users[id]
	return atomic.CompareAndSwapInt32(&user.State, oldstate, newstate)
}

func GetState(id int32) int32 {
	_userlock.RLock()
	defer _userlock.RUnlock()

	user := _users[id]
	return atomic.LoadInt32(&user.State)
}
