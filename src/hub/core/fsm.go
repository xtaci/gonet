package core

import (
	"sync"
	"sync/atomic"
	"time"
)

import (
	"misc/timer"
	. "types"
)

//------------------------------------------------ 状态机定义
const (
	UNKNOWN = byte(iota)
	OFF_FREE
	OFF_RAID
	ON_FREE
	ON_PROT
	OFF_PROT
)

const (
	RAID_TIME = 300 // seconds
	EVENT_MAX = 50000
)

//--------------------------------------------------------- player info
type PlayerInfo struct {
	Id             int32
	State          byte
	ProtectTimeout int64      // unix time
	RaidTimeout    int64      // unix time
	Host           int32      // host
	WaitEventId    int32      // current waiting event id, a user will only wait on ONE timeout event,  PROTECTTIMEOUT of RAID TIMEOUT
	_lock          sync.Mutex // Record lock
}

func (p *PlayerInfo) Lock() {
	p._lock.Lock()
}

func (p *PlayerInfo) Unlock() {
	p._lock.Unlock()
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
 * A: players
 * B: playerinfo.LCK
 **********************************************************/
var (
	_lock_players sync.RWMutex          // lock players
	_players      map[int32]*PlayerInfo // all players

	_waits      map[int32]*PlayerInfo
	_waits_lock sync.Mutex
	_waits_ch   chan int32

	_event_id_gen int32
)

func init() {
	_players = make(map[int32]*PlayerInfo)
	_waits = make(map[int32]*PlayerInfo)
	_waits_ch = make(chan int32, EVENT_MAX)

	go _expire()
}

//--------------------------------------------------------- expires
func _expire() {
	for {
		select {
		case event_id := <-_waits_ch:
			_waits_lock.Lock()
			player := _waits[event_id]
			delete(_waits, event_id)
			_waits_lock.Unlock()

			if player != nil {
				player.Lock()
				if player.WaitEventId == event_id { // check if it is the waiting event, or just ignore
					switch player.State {
					case ON_PROT:
						player.State = ON_FREE
					case OFF_PROT:
						player.State = OFF_FREE
					}
				}
				player.Unlock()
			}
		}
	}
}

//------------------------------------------------ add a user to finite state machine manager
func _add_fsm(user *User) {
	player := &PlayerInfo{Id: user.Id}

	if user.ProtectTimeout > time.Now().Unix() {
		player.ProtectTimeout = user.ProtectTimeout
		player.State = OFF_PROT
	} else {
		player.State = OFF_FREE
	}

	_lock_players.Lock()
	_players[user.Id] = player
	_lock_players.Unlock()
}

//------------------------------------------------ when game disconnect, perform a logout for all players on that gs
func LogoutServer(host int32) {
	_lock_players.RLock()
	snapshot := make(map[int32]*PlayerInfo)
	for k, v := range _players {
		snapshot[k] = v
	}
	_lock_players.RUnlock()

	for _, v := range snapshot {
		v.Lock()
		if v.Host == host {
			switch v.State {
			case ON_FREE:
				v.State = OFF_FREE
			case ON_PROT:
				v.State = OFF_PROT
			}
		}
		v.Unlock()
	}
}

//------------------------------------------------ The State Machine Of Player
func Login(id, host int32) bool {
	_lock_players.RLock()
	player := _players[id]
	_lock_players.RUnlock()

	if player != nil {
		player.Lock()
		defer player.Unlock()

		switch player.State {
		case OFF_FREE:
			player.State = ON_FREE
			player.Host = host
			return true
		case OFF_PROT:
			player.State = ON_PROT
			player.Host = host
			return true
		default:
			return false
		}
	}
	return false
}

func Logout(id int32) bool {
	_lock_players.RLock()
	player := _players[id]
	_lock_players.RUnlock()

	if player != nil {
		player.Lock()
		defer player.Unlock()

		switch player.State {
		case ON_FREE:
			player.State = OFF_FREE
		case ON_PROT:
			player.State = OFF_PROT
		default:
			return false
		}
		return true
	}
	return false
}

func Raid(id int32) bool {
	_lock_players.RLock()
	player := _players[id]
	_lock_players.RUnlock()

	if player != nil {
		player.Lock()
		defer player.Unlock()

		switch player.State {
		case OFF_FREE:
			timeout := time.Now().Unix() + RAID_TIME
			event_id := atomic.AddInt32(&_event_id_gen, 1)
			timer.Add(event_id, timeout, _waits_ch) // generate timer

			player.State = OFF_RAID
			player.RaidTimeout = timeout
			player.WaitEventId = event_id

			_waits_lock.Lock()
			_waits[event_id] = player
			_waits_lock.Unlock()
		default:
			return false
		}
		return true
	}
	return false
}

func Free(id int32) bool {
	_lock_players.RLock()
	player := _players[id]
	_lock_players.RUnlock()

	if player != nil {
		player.Lock()
		defer player.Unlock()

		switch player.State {
		case ON_PROT:
			player.State = ON_FREE
		case OFF_RAID:
			player.State = OFF_FREE
		default:
			return false
		}
		return true
	}
	return false
}

func Protect(id int32, until int64) bool {
	_lock_players.RLock()
	player := _players[id]
	_lock_players.RUnlock()

	if player != nil {
		player.Lock()
		defer player.Unlock()

		switch player.State {
		case ON_FREE:
			event_id := atomic.AddInt32(&_event_id_gen, 1)
			timer.Add(event_id, until, _waits_ch)

			player.State = ON_PROT
			player.ProtectTimeout = until
			player.WaitEventId = event_id

			_waits_lock.Lock()
			_waits[event_id] = player
			_waits_lock.Unlock()
		case OFF_RAID:
			event_id := atomic.AddInt32(&_event_id_gen, 1)
			timer.Add(event_id, until, _waits_ch)

			player.State = OFF_PROT
			player.ProtectTimeout = until
			player.WaitEventId = event_id

			_waits_lock.Lock()
			_waits[event_id] = player
			_waits_lock.Unlock()
		case ON_PROT:
			event_id := atomic.AddInt32(&_event_id_gen, 1)
			timer.Add(event_id, until, _waits_ch)

			player.State = ON_PROT
			player.ProtectTimeout = until
			player.WaitEventId = event_id

			_waits_lock.Lock()
			_waits[event_id] = player
			_waits_lock.Unlock()
		default:
			return false
		}
		return true
	}

	return false
}

// State Readers
// A->B
func State(id int32) (ret byte) {
	_lock_players.RLock()
	player := _players[id]
	_lock_players.RUnlock()

	if player != nil {
		player.Lock()
		ret = player.State
		player.Unlock()
	}
	return
}

func ProtectTimeout(id int32) (ret int64) {
	_lock_players.RLock()
	player := _players[id]
	_lock_players.RUnlock()

	if player != nil {
		player.Lock()
		ret = player.ProtectTimeout
		player.Unlock()
	}

	return
}

func Host(id int32) (ret int32) {
	_lock_players.RLock()
	player := _players[id]
	_lock_players.RUnlock()

	if player != nil {
		player.Lock()
		ret = player.Host
		player.Unlock()
	}

	return
}
