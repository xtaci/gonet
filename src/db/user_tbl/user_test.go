package user_tbl

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	user := fmt.Sprintf("test%v", rand.Int())
	mac := fmt.Sprintf("mac:%v", rand.Int())

	New(user, mac)
	if LoginMac(user, mac) == nil {
		t.Error("login failed")
	}

	all := GetAll()
	for _, v := range all {
		fmt.Println(v)
	}
}

func BenchmarkCreateUser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		user := fmt.Sprintf("test%v", i)
		mac := fmt.Sprintf("mac%v", i)
		New(user, mac)
	}
}
