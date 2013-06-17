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
	COLLECTION             = "GROUPS"
	GROUP_GEN              = "GROUP_GEN"
	DEFAULT_GROUP_MESSAGES = 128
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

//------------------------------------------------ Group
type GroupInfo struct {
	GroupId  int32
	Leader   int32
	Name     string
	Desc     string
	Messages []*IPCObject
	MaxMsgId uint32

	// runtime
	_members MemberSlice
}

var (
	_groups      map[int32]*GroupInfo  // id -> groupinfo
	_group_names map[string]*GroupInfo // name-> groupinfo
	_lock        sync.RWMutex
)

func init() {
	_groups = make(map[int32]*GroupInfo)
	_group_names = make(map[string]*GroupInfo)
}

//------------------------------------------------ create group
func Create(creator_id int32, groupname string) (groupid int32, succ bool) {
	_lock.Lock()
	defer _lock.Unlock()

	if _group_names[groupname] == nil {
		groupid := db.NextVal(GROUP_GEN)
		group := &GroupInfo{GroupId: groupid, Name: groupname, Leader: creator_id}
		group._members._add(creator_id)

		// index
		_groups[groupid] = group
		_group_names[groupname] = group

		// save
		_save(group)

		return groupid, true
	}

	return 0, false
}

func Group(groupid int32) *GroupInfo {
	_lock.Lock()
	defer _lock.Unlock()
	return _groups[groupid]
}

func (group *GroupInfo) Members() []int32 {
	_lock.RLock()
	defer _lock.RUnlock()

	m := make([]int32, len(group._members.M))
	copy(m, group._members.M)

	return m
}

//------------------------------------------------ Join group
func (group *GroupInfo) Join(user_id int32) {
	_lock.Lock()
	defer _lock.Unlock()

	group._members._add(user_id)
	_save(group)
}

//------------------------------------------------ leave group
func (group *GroupInfo) Leave(user_id int32) {
	_lock.Lock()
	defer _lock.Unlock()

	defer func() { // if no member, delete group
		if group._members.Len() == 0 {
			delete(_groups, group.GroupId)
			delete(_group_names, group.Name)
			c := db.Collection(COLLECTION)
			err := c.Remove(bson.M{"groupid": group.GroupId})
			if err != nil {
				log.Println(err)
			}
		}
	}()

	group._members._remove(user_id)
}

//------------------------------------------------ get group ranklist
func (group *GroupInfo) Ranklist() []int32 {
	_lock.RLock()
	defer _lock.RUnlock()

	group._members.Sort()
	m := make([]int32, len(group._members.M))
	copy(m, group._members.M)

	return m
}

//------------------------------------------------  push message to group
func (group *GroupInfo) Push(obj *IPCObject) {
	_lock.Lock()
	defer _lock.Unlock()

	// group message max
	config := cfg.Get()
	msg_max, err := strconv.Atoi(config["group_msg_max"])
	if err != nil {
		log.Println("group_msg_max:", err)
		msg_max = DEFAULT_GROUP_MESSAGES
	}

	if len(group.Messages) >= msg_max {
		group.Messages = group.Messages[1:]
	}

	group.Messages = append(group.Messages, obj)
	group.MaxMsgId += 1
	_save(group)
}

func (group *GroupInfo) Recv(lastmsg_id uint32) []*IPCObject {
	_lock.RLock()
	defer _lock.RUnlock()

	if lastmsg_id >= group.MaxMsgId {
		return nil
	}

	count := int(group.MaxMsgId - lastmsg_id)
	if count > len(group.Messages) {
		return group.Messages
	} else {
		return group.Messages[len(group.Messages)-count:]
	}

	return nil
}

//------------------------------------------------ Save to db
func _save(group *GroupInfo) {
	c := db.Collection(COLLECTION)
	info, err := c.Upsert(bson.M{"groupid": group.GroupId}, group)
	if err != nil {
		log.Println(info, err)
	}
}
