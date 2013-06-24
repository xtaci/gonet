package user_tbl

import (
	"fmt"
	"testing"
)

func TestUser(t *testing.T) {
	user := fmt.Sprintf("test%v", 1000000)
	mac := fmt.Sprintf("mac:%v", 1000000)

	New(user, mac)
	if LoginMac(user, mac) == nil {
		t.Error("login failed")
	}

	all := GetAll()
	for _, v := range all {
		fmt.Println(v)
	}

	u := Query(user)
	if u == nil {
		t.Error("cannot query by name")
	}
}

func BenchmarkCreateUser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		user := fmt.Sprintf("test%v", i)
		mac := fmt.Sprintf("mac%v", i)
		New(user, mac)
	}
}
