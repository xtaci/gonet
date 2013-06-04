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
	if !Join(2, clanid) {
		t.Error("cannot join clan")
	}

	if !Leave(2, clanid) {
		t.Error("cannot leave join")
	}
}

var clanid uint32

func init() {
	clanid, _ = Create(1, "tap4funbenchmark")
}

func BenchmarkClan(b *testing.B) {
	fmt.Println("CLANID", clanid)

	for i := 0; i < b.N; i++ {
		user := &User{Id: int32(i)}
		_add_fsm(user)
	}

	for i := 0; i < b.N; i++ {
		Join(int32(i), clanid)
	}

	for i := 0; i < b.N; i++ {
		Leave(int32(i), clanid)
	}
}
