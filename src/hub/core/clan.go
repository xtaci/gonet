package core

import (
	//"fmt"
	"labix.org/v2/mgo/bson"
	"log"
	"sort"
	"strconv"
	"sync"
)

import (
	"cfg"
	"db"
	. "types"
)

const (
	COLLECTION = "CLANS"
	CLANGEN    = "CLANGEN"
)

type MemberSlice struct {
	M []int32
}

//------------------------------------------------ Add a member,make sure not twice added
func (mem *MemberSlice) _add(user_id int32) {
	for k := range mem.M {
		if mem.M[k] == user_id {
			return
		}
	}

	mem.M = append(mem.M, user_id)
}

func (mem *MemberSlice) _remove(user_id int32) {
	idx := -1
	for k := range mem.M {
		if mem.M[k] == user_id {
			idx = k
			break
		}
	}

	if idx > 0 {
		mem.M = append(mem.M[:idx], mem.M[idx+1:]...)
	}

}

func (mem *MemberSlice) Len() int {
	return len(mem.M)
}

//------------------------------------------------ sort in descending order
func (mem *MemberSlice) Less(i, j int) bool {
	return Score(mem.M[i]) > Score(mem.M[j])
}

func (mem *MemberSlice) Sort() {
	sort.Sort(mem)
}

//------------------------------------------------ XOR swap
func (mem *MemberSlice) Swap(i, j int) {
	mem.M[i] = mem.M[i] ^ mem.M[j]
	mem.M[j] = mem.M[i] ^ mem.M[j]
	mem.M[i] = mem.M[i] ^ mem.M[j]
}

//------------------------------------------------ Clan
type ClanInfo struct {
	ClanId   int32
	Leader   int32
	Name     string
	Desc     string
	Messages []*IPCObject
	MaxMsgId uint32

	// runtime
	_members MemberSlice
}

var (
	_clans      map[int32]*ClanInfo  // id -> claninfo
	_clan_names map[string]*ClanInfo // name-> claninfo
	_lock       sync.RWMutex
)

func init() {
	_clans = make(map[int32]*ClanInfo)
	_clan_names = make(map[string]*ClanInfo)
}

//------------------------------------------------ create clan
func Create(creator_id int32, clanname string) (clanid int32, succ bool) {
	_lock.Lock()
	defer _lock.Unlock()

	if _clan_names[clanname] == nil {
		clanid := db.NextVal(CLANGEN)
		clan := &ClanInfo{ClanId: clanid, Name: clanname, Leader: creator_id}
		clan._members._add(creator_id)

		// index
		_clans[clanid] = clan
		_clan_names[clanname] = clan

		// save
		_save(clan)

		return clanid, true
	}

	return 0, false
}

func Clan(clanid int32) *ClanInfo {
	_lock.Lock()
	defer _lock.Unlock()
	return _clans[clanid]
}

func (clan *ClanInfo) Members() []int32 {
	_lock.RLock()
	defer _lock.RUnlock()

	m := make([]int32, len(clan._members.M))
	copy(m, clan._members.M)

	return m
}

//------------------------------------------------ Join clan
func (clan *ClanInfo) Join(user_id int32) {
	_lock.Lock()
	defer _lock.Unlock()

	clan._members._add(user_id)
	_save(clan)
}

//------------------------------------------------ leave clan
func (clan *ClanInfo) Leave(user_id int32) {
	_lock.Lock()
	defer _lock.Unlock()

	defer func() { // if no member, delete clan
		if clan._members.Len() == 0 {
			delete(_clans, clan.ClanId)
			delete(_clan_names, clan.Name)
			c := db.Collection(COLLECTION)
			err := c.Remove(bson.M{"clanid": clan.ClanId})
			if err != nil {
				log.Println(err)
			}
		}
	}()

	clan._members._remove(user_id)
}

//------------------------------------------------ get clan ranklist
func (clan *ClanInfo) Ranklist() []int32 {
	_lock.RLock()
	defer _lock.RUnlock()

	clan._members.Sort()
	m := make([]int32, len(clan._members.M))
	copy(m, clan._members.M)

	return m
}

//------------------------------------------------  push message to clan
func (clan *ClanInfo) Push(obj *IPCObject) {
	_lock.Lock()
	defer _lock.Unlock()

	// clan message max
	config := cfg.Get()
	msg_max, err := strconv.Atoi(config["clan_msg_max"])
	if err != nil {
		log.Println("clan_msg_max:", err)
	}

	if len(clan.Messages) >= msg_max {
		clan.Messages = clan.Messages[1:]
	}

	clan.Messages = append(clan.Messages, obj)
	clan.MaxMsgId += 1
	_save(clan)
}

//------------------------------------------------ Save to db
func _save(clan *ClanInfo) {
	c := db.Collection(COLLECTION)
	info, err := c.Upsert(bson.M{"clanid": clan.ClanId}, clan)
	if err != nil {
		log.Println(info, err)
	}
}
