package ranklist

import (
	"sync"
	"time"
)

import (
	. "types"
)

// OFFLINE|RAID, OFFLINE|PROTECTED , OFFLINE |FREE
// ONLINE|PROTECTED , ONLINE|FREE
const (
	HISHIFT = 16

	// 
	OFFLINE = 1 << HISHIFT
	ONLINE  = 1 << 1 << HISHIFT

	// battle status
	RAID      = 1 << 2
	PROTECTED = 1 << 3
	FREE      = 1 << 4

	LOMASK = 0xFFFF
	HIMASK = 0xFFFF0000
)

const (
	RAID_TIME = 300 // seconds
)

var _raidtime_max int64

//--------------------------------------------------------- player info 
type PlayerInfo struct {
	Id          int32
	State       int32
	ProtectTime int64 // unix time
	RaidStart   int64 // unix time
	Clan        int32 // clan info
	Host		int32 // host
	Name        string
	LCK         sync.Mutex // Record lock
}

/**********************************************************
 * consider following deadlock situations. 
 *
 * A(B) means lock A,lock B, unlock B, unlock A
 * A->B means lockA unlockA,then lockB, unlockB
 *
 * p:A(B), q:B(A), possible circular wait, deadlock!!!
 * p:A(B), q:A(B), ok 
 * p.A(B), q:B or A, ok
 * p:A->B, q: B->A, ok 
 *
 * make sure acquiring the lock IN SEQUENCE. i.e.
 * A: one of [players, raids, protects]
 * B: playerinfo.LCK
 **********************************************************/
var (
	_lock_players sync.RWMutex          // lock players
	_players      map[int32]*PlayerInfo // all players

	_lock_raids sync.Mutex            // lock raids
	_raids      map[int32]*PlayerInfo // being raided

	_lock_protects sync.Mutex            // lock protects
	_protects      map[int32]*PlayerInfo // protecting
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

		_lock_protects.Lock()
		now := time.Now().Unix()
		for k, v := range _protects {
			v.LCK.Lock()
			if v.ProtectTime < now {
				// PROTECTED->FREE
				v.State = v.State&(^PROTECTED) | FREE
				delete(_protects, k)
			}
			v.LCK.Unlock()

		}
		_lock_protects.Unlock()

		_lock_raids.Lock()
		for k, v := range _raids {
			v.LCK.Lock()
			if now-v.RaidStart > _raidtime_max {
				// RAID->FREE
				v.State = v.State&(^RAID) | FREE
				delete(_raids, k)
			}
			v.LCK.Unlock()

		}
		_lock_raids.Unlock()

		time.Sleep(time.Minute)
	}
}

//------------------------------------------------ add a user to finite state machine manager
func _add_fsm(ud *User) {
	info := &PlayerInfo{Id: ud.Id, Name: ud.Name, State: OFFLINE | FREE}

	_lock_players.Lock()
	_players[ud.Id] = info
	_lock_players.Unlock()
}

// The State Machine Of Player
//----------------------------------------------- OFFLINE->ONLINE
// A->B
func Login(id, host int32) bool {
	_lock_players.RLock()
	player := _players[id]
	_lock_players.RUnlock()

	if player != nil {
		player.LCK.Lock()
		defer player.LCK.Unlock()
		state := player.State

		if state&OFFLINE != 0 && state&RAID == 0 {
			player.State = int32(ONLINE | (state & LOMASK))
			player.Host = host
			return true
		}
	}

	return false
}

//----------------------------------------------- ONLINE->OFFLINE
// A->B
func Logout(id int32) bool {
	_lock_players.RLock()
	player := _players[id]
	_lock_players.RUnlock()

	if player != nil {
		player.LCK.Lock()
		defer player.LCK.Unlock()

		state := player.State

		if state&ONLINE != 0 {
			player.State = int32(OFFLINE | (state & LOMASK))
			return true
		}
	}

	return false
}

//----------------------------------------------- FREE->RAID
// A->B->A
func Raid(id int32) bool {
	_lock_players.RLock()
	player := _players[id]
	_lock_players.RUnlock()

	if player != nil {
		player.LCK.Lock()

		state := player.State

		if state&OFFLINE != 0 && state&FREE != 0 {
			player.State = int32(OFFLINE | RAID)
			player.RaidStart = time.Now().Unix()
			player.LCK.Unlock()

			_lock_raids.Lock()
			_raids[id] = player // add to raid queue
			_lock_raids.Unlock()
			return true
		}

		player.LCK.Unlock()
	}

	return false
}

//----------------------------------------------- RAID->PROTECTED
// A->B->A->A
func Protect(id int32, until time.Time) bool {
	_lock_raids.Lock()
	player := _raids[id]
	_lock_raids.Unlock()

	if player != nil {
		player.LCK.Lock()

		state := player.State

		if state&RAID != 0 {
			player.State = int32(OFFLINE | PROTECTED)
			player.ProtectTime = until.Unix()
			player.LCK.Unlock()

			_lock_raids.Lock()
			delete(_raids, id) // remove from raids
			_lock_raids.Unlock()

			_lock_protects.Lock()
			_protects[id] = player // add to protects
			_lock_protects.Unlock()

			return true
		}
		player.LCK.Unlock()
	}

	return false
}

//----------------------------------------------- RAID->FREE
// A(B)
func Free(id int32) bool {
	_lock_raids.Lock()
	defer _lock_raids.Unlock()

	player := _raids[id]

	if player != nil {
		player.LCK.Lock()
		defer player.LCK.Unlock()

		if player.State&RAID != 0 {
			player.State = int32(OFFLINE | FREE)
			delete(_raids, id) // remove from raids
			return true
		}
	}

	return false
}

//----------------------------------------------- PROTECT->FREE
// A(B)
func Unprotect(id int32) bool {
	_lock_protects.Lock()
	defer _lock_protects.Unlock()

	player := _protects[id]

	if player != nil {
		player.LCK.Lock()
		defer player.LCK.Unlock()

		if player.State&RAID != 0 {
			player.State = int32(ONLINE | FREE)
			delete(_protects, id) // remove from protects
			return true
		}
	}

	return false
}

// State Readers
// A->B 
func State(id int32) (ret int32) {
	_lock_players.RLock()
	player := _players[id]
	_lock_players.RUnlock()

	if player != nil {
		player.LCK.Lock()
		ret = player.State
		player.LCK.Unlock()
	}
	return
}

func ProtectTime(id int32) (ret int64) {
	_lock_players.RLock()
	player := _players[id]
	_lock_players.RUnlock()

	if player != nil {
		player.LCK.Lock()
		ret = player.ProtectTime
		player.LCK.Unlock()
	}

	return
}

func Name(id int32) (ret string) {
	_lock_players.RLock()
	player := _players[id]
	_lock_players.RUnlock()

	if player != nil {
		player.LCK.Lock()
		ret = player.Name
		player.LCK.Unlock()
	}
	return
}

func Host(id int32) (ret int32) {
	_lock_players.RLock()
	player := _players[id]
	_lock_players.RUnlock()

	if player != nil {
		player.LCK.Lock()
		ret = player.Host
		player.LCK.Unlock()
	}

	return
}

func Clan(id int32) (ret int32) {
	_lock_players.RLock()
	player := _players[id]
	_lock_players.RUnlock()

	if player != nil {
		player.LCK.Lock()
		ret = player.Clan
		player.LCK.Unlock()
	}

	return
}
