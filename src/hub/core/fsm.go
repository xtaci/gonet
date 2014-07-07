package core

import (
	"sync"
	"time"
)

import (
	"misc/timer"
	. "types"
)

const (
	AUTO_EXPIRE = 300 // seconds
	EVENT_MAX   = 50000
)

//--------------------------------------------------------- player info
type PlayerInfo struct {
	Id             int32
	State          byte
	Host           int32      // host
	ProtectTimeout int64      // real protect timeout
	RaidStart      int64      // raid start time
	sync.Mutex // Record lock
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
	_lock_players sync.Mutex            // lock players
	_players      map[int32]*PlayerInfo // all players
	_waits_ch     chan int32
)

func init() {
	_players = make(map[int32]*PlayerInfo)
	_waits_ch = make(chan int32, EVENT_MAX)

	go _expire()
}

//--------------------------------------------------------- expires
func _expire() {
	for {
		select {
		case user_id := <-_waits_ch:
			_lock_players.Lock()
			player := _players[user_id]
			_lock_players.Unlock()

			// modify state, using defensive strategy
			player.Lock()
			switch player.State {
			case ON_PROT:
				if player.ProtectTimeout <= time.Now().Unix() {
					player.State = ON_FREE
				}
			case OFF_PROT:
				if player.ProtectTimeout <= time.Now().Unix() {
					player.State = OFF_FREE
				}
			case OFF_RAID:
				if player.RaidStart+AUTO_EXPIRE <= time.Now().Unix() {
					player.State = OFF_FREE
				}
			}
			player.Unlock()
		}
	}
}

//------------------------------------------------ add a user to finite state machine manager
func _add_fsm(user *User) {
	player := &PlayerInfo{Id: user.Id}
	player.ProtectTimeout = user.ProtectTimeout

	if user.ProtectTimeout > time.Now().Unix() { // 有保护时间
		player.State = OFF_PROT

		_lock_players.Lock()
		_players[player.Id] = player
		_lock_players.Unlock()

		timer.Add(user.Id, user.ProtectTimeout, _waits_ch)
	} else {
		player.State = OFF_FREE

		_lock_players.Lock()
		_players[user.Id] = player
		_lock_players.Unlock()
	}
}

//------------------------------------------------ when game disconnect, perform a logout for all players on that gs
func LogoutServer(host int32) {
	_lock_players.Lock()
	snapshot := make(map[int32]*PlayerInfo)
	for k, v := range _players {
		snapshot[k] = v
	}
	_lock_players.Unlock()

	for _, v := range snapshot {
		v.Lock()
		if v.Host == host {
			// state correction
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
	_lock_players.Lock()
	player := _players[id]
	_lock_players.Unlock()

	if player != nil {
		player.Lock()
		defer player.Unlock()

		switch player.State {
		case OFF_FREE:
			player.State = ON_FREE
			player.Host = host
		case OFF_PROT:
			player.State = ON_PROT
			player.Host = host
		default:
			return false
		}
		return true
	}
	return false
}

func Logout(id int32) bool {
	_lock_players.Lock()
	player := _players[id]
	_lock_players.Unlock()

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
	_lock_players.Lock()
	player := _players[id]
	_lock_players.Unlock()

	if player != nil {
		player.Lock()
		defer player.Unlock()

		switch player.State {
		case OFF_FREE:
			player.State = OFF_RAID
			player.RaidStart = time.Now().Unix()
			timeout := time.Now().Unix() + AUTO_EXPIRE // automatic expire when in raid
			timer.Add(player.Id, timeout, _waits_ch)
		default:
			return false
		}
		return true
	}
	return false
}

func Free(id int32) bool {
	_lock_players.Lock()
	player := _players[id]
	_lock_players.Unlock()

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
	_lock_players.Lock()
	player := _players[id]
	_lock_players.Unlock()

	if player != nil {
		player.Lock()
		defer player.Unlock()

		switch player.State {
		case ON_FREE:
			player.State = ON_PROT
			player.ProtectTimeout = until
			timer.Add(id, until, _waits_ch)
		case OFF_RAID:
			player.State = OFF_PROT
			player.ProtectTimeout = until
			timer.Add(id, until, _waits_ch)
		case ON_PROT: // protect + protect
			player.ProtectTimeout = until
			timer.Add(id, until, _waits_ch)
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
	_lock_players.Lock()
	player := _players[id]
	_lock_players.Unlock()

	if player != nil {
		player.Lock()
		ret = player.State
		player.Unlock()
	}
	return
}

func Host(id int32) (ret int32) {
	_lock_players.Lock()
	player := _players[id]
	_lock_players.Unlock()

	if player != nil {
		player.Lock()
		ret = player.Host
		player.Unlock()
	}

	return
}
