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
	pass := fmt.Sprintf("pass%v", rand.Int())

	New(user, pass)
	if Login(user, pass) == nil {
		t.Error("login failed")
	}

	all := GetAll()
	for _, v := range all {
		fmt.Println(v)
	}
}
