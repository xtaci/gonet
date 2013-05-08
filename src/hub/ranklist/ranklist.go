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
	HISHIFT = 16

	// 
	OFFLINE = 1 << HISHIFT
	ONLINE  = 1 << 1 << HISHIFT

	// battle status
	RAID = 1 << 2
	PROTECTED  = 1 << 3
	FREE       = 1 << 4

	LOMASK = 0xFFFF
	HIMASK = 0xFFFF0000
)

// OFFLINE|RAID, OFFLINE|PROTECTED , OFFLINE |FREE
// ONLINE|PROTECTED , ONLINE|FREE

const (
	RAID_TIME = 300 // seconds
)

var _raidtime_max int64

//--------------------------------------------------------- player info 
type PlayerInfo struct {
	Id          int32
	Name        string
	Score       int32
	State       int32
	ProtectTime int64 // unix time
	RaidStart   int64 // unix time
}

var (
	_ranklist dos.Tree // dynamic order statistics
	_lock     sync.RWMutex
	_players  map[int32]*PlayerInfo // free players
	_raids    map[int32]*PlayerInfo // being raided
	_protects map[int32]*PlayerInfo // protecting
)

func init() {
	_players = make(map[int32]*PlayerInfo)
	_raids = make(map[int32]*PlayerInfo)
	_protects = make(map[int32]*PlayerInfo)
	_raidtime_max = RAID_TIME
}

//--------------------------------------------------------- expires
func ExpireRoutine() {
	for {
		_lock.Lock()
		now := time.Now().Unix()
		for k, v := range _protects {
			if v.ProtectTime < now {
				// PROTECTED->FREE
				v.State = v.State&(^PROTECTED) | FREE
				delete(_protects, k)
			}

		}
		_lock.Unlock()
		_lock.Lock()
		for k, v := range _raids {
			if now-v.RaidStart > _raidtime_max {
				// RAID->FREE
				v.State = v.State&(^RAID) | FREE
				delete(_raids, k)
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

	info := &PlayerInfo{Id: ud.Id, Name: ud.Name, Score: ud.Score, State: OFFLINE | FREE}
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
// OFFLINE->ONLINE
func Login(id int32) bool {
	_lock.Lock()
	defer _lock.Unlock()

	player := _players[id]

	if player != nil {
		state := player.State

		if state&OFFLINE != 0 && state&RAID == 0 {
			player.State = int32(ONLINE | (state & LOMASK))
			return true
		}
	}

	return false
}

// Logout a User
func Logout(id int32) bool {
	_lock.Lock()
	defer _lock.Unlock()

	player := _players[id]

	if player != nil {
		state := player.State

		if state&ONLINE != 0 {
			player.State = int32(OFFLINE|(state & LOMASK))
			return true
		}
	}

	return false
}

// raid a player
// make sure no user is attacked by more than 1 player
// check return value before proceed.
// FREE->RAID
func Raid(id int32) bool {
	_lock.Lock()
	defer _lock.Unlock()

	player := _players[id]

	if player != nil {
		state := player.State

		if state&OFFLINE != 0 && state&FREE != 0 {
			player.State = int32(OFFLINE | RAID)
			player.RaidStart = time.Now().Unix()
			_raids[id] = player // add to being raided queue
			return true
		}
	}

	return false
}

// the player should be protected when the raid is over
// RAID->PROTECTED
func Protect(id int32, until time.Time) bool {
	_lock.Lock()
	defer _lock.Unlock()

	player := _raids[id]

	if player != nil {
		state := player.State

		if state&RAID != 0 {
			player.State = int32(OFFLINE | PROTECTED)
			player.ProtectTime = until.Unix()
			delete(_raids, id)     // remove from raids
			_protects[id] = player // add to protects
			return true
		}
	}

	return false
}

// the player should NOT be protected when the raid is over
// RAID->FREE
func Free(id int32) bool {
	_lock.Lock()
	defer _lock.Unlock()

	player := _raids[id]

	if player != nil {
		if player.State&RAID != 0 {
			player.State = int32(OFFLINE|FREE)
			delete(_raids, id) // remove from raids
			return true
		}
	}

	return false
}

// a online player spontanenous give up protection
// PROTECT->FREE
func Unprotect(id int32) bool {
	_lock.Lock()
	defer _lock.Unlock()

	player := _protects[id]

	if player != nil {
		if player.State&RAID != 0 {
			player.State = int32(ONLINE|FREE)
			delete(_protects, id) // remove from raids
			return true
		}
	}

	return false
}

// Readers
func State(id int32) int32 {
	_lock.RLock()
	defer _lock.RUnlock()

	return _players[id].State
}

func ProtectTime(id int32) int64 {
	_lock.RLock()
	defer _lock.RUnlock()

	return _players[id].ProtectTime
}
