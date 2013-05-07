package ranklist

import (
	"sync"
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

const (
	RAID_TIME = 300 // seconds
)

var _raidtime_max int64

//--------------------------------------------------------- player info 
type PlayerInfo struct {
	Id	int32
	Name string
	Score int32
	State int
	ProtectTime time.Time
	RaidStart time.Time
}

var (
	_ranklist dos.Tree
	_lock sync.RWMutex
	_players map[int32]*PlayerInfo
	_count int32
)

func init() {
	_players = make(map[int32]*PlayerInfo)
	_raidtime_max = RAID_TIME
}

//--------------------------------------------------------- expires
func ExpireRoutine() {

	for {

		_lock.Lock()
		for _,v := range _players {
			if v.State == PROTECTED && v.ProtectTime.Unix() <  time.Now().Unix() {
				// PROTECTED->FREE
				v.State = FREE
			} else if v.State == BEING_RAID && time.Now().Unix() - v.RaidStart.Unix() > _raidtime_max {
				// expire BEING_RAID
				// in case someone do evil.
				v.State = FREE
			}
		}
		_lock.Unlock()

		time.Sleep(time.Minute)
	}
}

// add a user to rank list, only useful when startup & register
func AddUser(ud *User) {
	_lock.Lock()
	defer _lock.Unlock()

	info := &PlayerInfo{Id:ud.Id, Name:ud.Name, Score:ud.Score, State:FREE }
	_players[ud.Id] = info
	_ranklist.Insert(int(ud.Score), info)
}


//========================================================= Rank List operations
// change score of user
// the DOS tree scheme for changing a value is delete & insert
func ChangeScore(id, oldscore, newscore int32) bool {
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
			return false
		}

		if n.Data().(*PlayerInfo).Id == id {
			_ranklist.DeleteNode(n)
			n.Data().(*PlayerInfo).Score = newscore
			_ranklist.Insert(int(newscore), n.Data().(*PlayerInfo))
			return true
		} else {
			// temporary delete 
			_ranklist.DeleteNode(n)
			tmplist = append(tmplist, n.Data())
		}
	}

	return true
}

// players count
func Count() int {
	_lock.RLock()
	defer _lock.RUnlock()
	return _ranklist.Count()
}

// get ranklist snapshot in [from, to] 
func GetRankList(from, to int) []PlayerInfo {
	sublist := make([]PlayerInfo, to-from+1)

	_lock.RLock()
	defer _lock.RUnlock()

	for i := from; i <= to; i++ {
		sublist[i-from] = *_ranklist.Rank(i).Data().(*PlayerInfo)
	}

	return sublist
}


//========================================================= The State Machine Of Player
// the user is not allowed to login 
// when being attacked
func Login(id int32) bool {
	_lock.RLock()
	defer _lock.RUnlock()

	player := _players[id]

	if player.State == FREE || player.State == PROTECTED {
		player.State = ONLINE
		return true
	}

	return false
}

// raid a player
// make sure no user is attacked by more than 1 player
// check return value before proceed.
func Raid(id int32) bool {
	_lock.Lock()
	defer _lock.Unlock()

	player := _players[id]

	if player.State == FREE {
		player.State = BEING_RAID
		player.RaidStart = time.Now()
		return true
	}

	return false
}

// the player should be protected when the raid is over
func Protect(id int32, until time.Time) bool {
	_lock.Lock()
	defer _lock.Unlock()

	player := _players[id]

	if player.State == BEING_RAID {
		player.ProtectTime = until
		player.State = PROTECTED
		return true
	}

	return false
}

// the player should NOT be protected when the raid is over
func Free(id int32) bool {
	_lock.Lock()
	defer _lock.Unlock()

	player := _players[id]

	if player.State == BEING_RAID {
		player.State = FREE
		return true
	}

	return false
}

// Reader
func State(id int32) int {
	_lock.Lock()
	defer _lock.Unlock()

	return _players[id].State
}

func ProtectTime(id int32) time.Time {
	_lock.Lock()
	defer _lock.Unlock()

	return _players[id].ProtectTime
}
