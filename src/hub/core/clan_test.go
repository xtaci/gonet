package core

import (
	"fmt"
	"testing"
	. "types"
)

func TestClan(t *testing.T) {
	clanid, succ := Create(1, "tap4fun")
	if !succ {
		t.Error("cannot create clan")
	}

	fmt.Println(clanid, succ)

	for i := 0; i < 100; i++ {
		user := &User{Id: int32(i), Score: int32(i)}
		_add_rank(user)

		if !Join(int32(i), clanid) {
			t.Error("cannot join clan")
		}
	}
	fmt.Println(_clans[clanid]._members.M)

	rl := Ranklist(clanid)

	if rl[0] != 99 || rl[99] != 0 {
		t.Error("clan ranklist failed")
	}

	for i := 0; i < 100; i++ {
		if !Leave(int32(i), clanid) {
			t.Error("cannot leave join")
		}
	}

	fmt.Println("testing send & recv")

	for i := 0; i < 200; i++ {
		Send(nil, 1)
	}

	result := Recv(195, 1)

	fmt.Println(len(result))

	if len(result) != 5 {
		t.Error("send recv failed on size")
	}
}

var clanid int32

func init() {
	clanid, _ = Create(1, "tap4funbenchmark")
}

func BenchmarkClan(b *testing.B) {
	fmt.Println("CLANID", clanid)

	for i := 0; i < b.N; i++ {
		user := &User{Id: int32(i), Score: int32(i)}
		_add_rank(user)
	}

	for i := 0; i < b.N; i++ {
		Join(int32(i), clanid)
	}

	Ranklist(clanid)

	for i := 0; i < b.N; i++ {
		Leave(int32(i), clanid)
	}

}
