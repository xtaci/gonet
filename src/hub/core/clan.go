package core

import (
	//"fmt"
	"labix.org/v2/mgo/bson"
	"log"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
)

import (
	"cfg"
	"db"
	. "types"
)

const (
	COLLECTION = "CLANS"
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
	ClanId   uint32
	Leader   int32
	Name     string
	Desc     string
	Messages []*IPCObject
	MaxMsgId uint32

	// runtime
	_members MemberSlice
}

var (
	_clans      map[uint32]*ClanInfo // id -> claninfo
	_clan_names map[string]*ClanInfo // name-> claninfo
	_lock       sync.RWMutex

	_clanid_gen uint32
)

func init() {
	_clans = make(map[uint32]*ClanInfo)
	_clan_names = make(map[string]*ClanInfo)
}

//------------------------------------------------ create clan
func Create(creator_id int32, clanname string) (clanid uint32, succ bool) {
	_lock.Lock()
	defer _lock.Unlock()

	if _clan_names[clanname] == nil {
		clanid := atomic.AddUint32(&_clanid_gen, 1)
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

//------------------------------------------------ Join clan
func Join(user_id int32, clanid uint32) bool {
	_lock.Lock()
	defer _lock.Unlock()

	//fmt.Println(_clans, clanid)
	clan := _clans[clanid]
	if clan != nil {
		clan._members._add(user_id)
		_save(clan)
		return true
	}
	return false
}

//------------------------------------------------ leave clan
func Leave(user_id int32, clanid uint32) bool {
	_lock.Lock()
	defer _lock.Unlock()

	clan := _clans[clanid]

	if clan != nil {
		defer func() { // if no member, delete clan
			if clan._members.Len() == 0 {
				delete(_clans, clanid)
				delete(_clan_names, clan.Name)
				c := db.Collection(COLLECTION)
				err := c.Remove(bson.M{"clanid": clanid})
				if err != nil {
					log.Println(err)
				}
			}
		}()

		clan._members._remove(user_id)
		return true
	}

	return false
}

//------------------------------------------------ get clan ranklist
func Ranklist(clanid uint32) []int32 {
	_lock.Lock()
	defer _lock.Unlock()

	clan := _clans[clanid]
	if clan != nil {
		clan._members.Sort()
		return clan._members.M
	}

	return nil
}

//------------------------------------------------  send message to clan
func Send(obj *IPCObject, clanid uint32) {
	_lock.Lock()
	defer _lock.Unlock()

	clan := _clans[clanid]
	// clan message max
	config := cfg.Get()
	msg_max, err := strconv.Atoi(config["clan_msg_max"])
	if err != nil {
		log.Println("clan_msg_max:", err)
	}

	if clan != nil {
		if len(clan.Messages) >= msg_max {
			clan.Messages = clan.Messages[1:]
		}

		clan.Messages = append(clan.Messages, obj)
		clan.MaxMsgId += 1
		_save(clan)
	}
}

func Recv(lastmsg_id uint32, clanid uint32) []*IPCObject {
	_lock.RLock()
	defer _lock.RUnlock()

	clan := _clans[clanid]
	if clan != nil {
		if lastmsg_id >= clan.MaxMsgId {
			return nil
		}

		count := int(clan.MaxMsgId - lastmsg_id)
		if count > len(clan.Messages) {
			return clan.Messages
		} else {
			return clan.Messages[len(clan.Messages)-count:]
		}
	}

	return nil
}

//------------------------------------------------

//------------------------------------------------ Save to db
func _save(clan *ClanInfo) {
	c := db.Collection(COLLECTION)
	info, err := c.Upsert(bson.M{"clanid": clan.ClanId}, clan)
	if err != nil {
		log.Println(info, err)
	}
}
