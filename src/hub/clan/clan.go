package clan

import (
	"sync"
)

type ClanInfo struct {
	Id      int32
	Name    string
	Members []int32
}

var _clans map[int32]*ClanInfo       // id -> claninfo
var _clan_names map[string]*ClanInfo // name-> claninfo

var _lock sync.RWMutex

func init() {
	_clans = make(map[int32]*ClanInfo)
	_clan_names = make(map[string]*ClanInfo)
}

//------------------------------------------------ create clan
func Create(id int32, name string) (int32, bool) {
	_lock.Lock()
	defer _lock.Unlock()

	if _clan_names[name] == nil {
		// TODO: add db ops for ID
		clanid := int32(0)
		clan := &ClanInfo{Id: clanid, Name: name, Members: []int32{id}}
		_clans[clan.Id] = clan
		_clan_names[clan.Name] = clan
		return 0, true
	}

	return -1, false
}

//------------------------------------------------ Join clan
func Join(id, clanid int32) bool {
	_lock.Lock()
	defer _lock.Unlock()

	clan := _clans[clanid]
	if clan != nil {
		var is_added = false
		for _, v := range clan.Members { // check collision
			if v == id {
				is_added = true
				break
			}
		}

		if !is_added {
			clan.Members = append(clan.Members, id)
			return true
		}
	}

	return false
}

//------------------------------------------------ leave clan
func Leave(id, clanid int32) bool {
	_lock.Lock()
	defer _lock.Unlock()

	clan := _clans[clanid]

	if clan != nil {
		defer func() { // if no member, delete clan
			if len(clan.Members) == 0 {
				delete(_clans, clan.Id) // TODO: persistent
				delete(_clan_names, clan.Name)
			}
		}()

		for k, v := range clan.Members { // find & delete
			if v == id {
				clan.Members = append(clan.Members[:k], clan.Members[k+1:]...)
				return true
			}
		}
	}

	return false
}

//------------------------------------------------  Send Message  
func Send(msg string, clanid int32) {
}
