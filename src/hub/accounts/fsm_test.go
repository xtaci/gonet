package accounts

import (
	"fmt"
	"testing"
	"time"
	. "types"
)

func TestFSM(t *testing.T) {
	user := &User{Id: 1}
	_add_fsm(user)

	v1 := Raid(1)
	v2 := Raid(2)

	fmt.Println(v1, v2)

	if v2 {
		t.Error("Raid consistency failed")
	}
}

func BenchmarkFSM(b *testing.B) {
	for i := 0; i < b.N; i++ {
		user := &User{Id: int32(i)}
		_add_fsm(user)
	}

	for i := 0; i < b.N; i++ {
		Raid(int32(i))
		Free(int32(i))
		Protect(int32(i), time.Now().Unix())
		UnProtect(int32(i))
	}
}
