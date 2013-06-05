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

	clan := Clan(clanid)
	fmt.Println(clanid, succ)

	for i := 0; i < 100; i++ {
		user := &User{Id: int32(i), Score: int32(i)}
		_add_rank(user)

		clan.Join(int32(i))
	}
	fmt.Println(_clans[clanid]._members.M)

	rl := clan.Ranklist()

	if rl[0] != 99 || rl[99] != 0 {
		t.Error("clan ranklist failed")
	}

	for i := 0; i < 100; i++ {
		clan.Leave(int32(i))
	}

	fmt.Println("testing send & recv")

	for i := 0; i < 200; i++ {
		clan.Push(nil)
	}
}

var clanid int32

func init() {
	clanid, _ = Create(1, "tap4funbenchmark")
}

func BenchmarkClan(b *testing.B) {
	fmt.Println("CLANID", clanid)
	clan := Clan(clanid)

	for i := 0; i < b.N; i++ {
		user := &User{Id: int32(i), Score: int32(i)}
		_add_rank(user)
	}

	for i := 0; i < b.N; i++ {
		clan.Join(int32(i))
	}

	clan.Ranklist()

	for i := 0; i < b.N; i++ {
		clan.Leave(int32(i))
	}

}
