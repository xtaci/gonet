package ipc

import (
	"fmt"
	"testing"
)

type TMPObj struct {
	A int32
	B int64
	C string
}

func init() {
	DialHub()
}

func TestPing(t *testing.T) {
	fmt.Println("testing PING")
	if !Ping() {
		t.Fatal()
	}
}

func TestLogin(t *testing.T) {
	if !Login(1) {
		t.Error("cannot login")
		fmt.Println("please run benchmark first")
	}

	if Login(1) {
		t.Error("login twice")
	}

	if !Logout(1) {
		t.Error("logout")
	}

	if !Login(1) {
		t.Error("login again")
	}
}

func BenchmarkForward(b *testing.B) {
	obj := &TMPObj{A: 10, B: 20, C: "test"}
	for i := 0; i < b.N; i++ {
		Send(0, 2, 1, obj)
	}
}
