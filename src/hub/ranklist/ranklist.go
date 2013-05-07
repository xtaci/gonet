package ranklist

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

import (
	"misc/alg/dos"
	. "types"
)

const (
	FREE = iota
	ONLINE
	BEING_RAID
	PROTECTED
)

//--------------------------------------------------------- striped version of user
type PlayerInfo struct {
	Id	int32
	Name string
	Score int32
	State int32
	ProtectTime time.Time
}

var (
	_ranklist dos.Tree
	_lock sync.RWMutex
	_players map[int32]*PlayerInfo
	_count int32
)

func init() {
	_players = make(map[int32]*PlayerInfo)
}

//--------------------------------------------------------- add a user to rank list, only useful when startup & register
func AddUser(ud *User) {
	_lock.Lock()
	defer _lock.Unlock()

	state := FREE
	if ud.ProtectTime.Unix() > time.Now().Unix() {
		state = PROTECTED
	}

	info := &PlayerInfo{Id:ud.Id, Name:ud.Name, Score:ud.Score, State:int32(state), ProtectTime:ud.ProtectTime}
	_players[ud.Id] = info
	_ranklist.Insert(int(ud.Score), info)
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

		if n.Data().(*PlayerInfo).Id == id {
			_ranklist.DeleteNode(n)
			n.Data().(*PlayerInfo).Score = newscore
			_ranklist.Insert(int(newscore), n.Data().(*PlayerInfo))
			return
		} else {
			// temporary delete 
			_ranklist.DeleteNode(n)
			tmplist = append(tmplist, n.Data())
		}
	}

	return
}

//--------------------------------------------------------- players count
func Count() int {
	_lock.RLock()
	defer _lock.RUnlock()
	return _ranklist.Count()
}

//--------------------------------------------------------- get players from ranklist in [from, to] 
func GetRankList(from, to int) []*PlayerInfo {
	sublist := make([]*PlayerInfo, to-from+1)

	_lock.RLock()
	defer _lock.RUnlock()

	for i := from; i <= to; i++ {
		sublist[i-from] = _ranklist.Rank(i).Data().(*PlayerInfo)
	}

	return sublist
}

// change player state
// the atomicity of state is guaranted by ranklist
func ChangeState(id int32, oldstate, newstate int32) bool {
	_lock.Lock()
	defer _lock.Unlock()

	player := _players[id]
	return atomic.CompareAndSwapInt32(&player.State, oldstate, newstate)
}

func GetState(id int32) int32 {
	_lock.RLock()
	defer _lock.RUnlock()
	return _players[id].State
}
