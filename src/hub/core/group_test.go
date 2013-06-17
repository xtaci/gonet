package core

import (
	"fmt"
	"testing"
	. "types"
)

func TestGroup(t *testing.T) {
	groupid, succ := Create(1, "t4f")
	if !succ {
		t.Error("cannot create group")
	}

	group := Group(groupid)
	fmt.Println(groupid, succ)

	for i := 0; i < 100; i++ {
		user := &User{Id: int32(i), Score: int32(i)}
		_add_rank(user)

		group.Join(int32(i))
	}
	fmt.Println(_groups[groupid]._members.M)

	rl := group.Ranklist()

	if rl[0] != 99 || rl[99] != 0 {
		t.Error("group ranklist failed")
	}

	for i := 0; i < 100; i++ {
		group.Leave(int32(i))
	}

	fmt.Println("testing send & recv")

	for i := 0; i < 5; i++ {
		obj := &IPCObject{}
		group.Push(obj)
	}

	result := group.Recv(0)

	fmt.Println(len(result))

	if len(result) != 5 {
		t.Error("send recv failed on size")
	}
}

func init() {
	for i := 0; i < 1000000; i++ {
		user := &User{Id: int32(i), Score: int32(i)}
		_add_rank(user)
	}
}

func BenchmarkGroupRankList(b *testing.B) {
	LoadGroups()
	group := Group(1)

	for i := 0; i < b.N; i++ {
		group.Join(int32(i))
	}

	group.Ranklist()

	for i := 0; i < b.N; i++ {
		group.Leave(int32(i))
	}
}
