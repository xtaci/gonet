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

func TestEventFunc(t *testing.T) {
	DialHub()
	fmt.Println("testing PING")
	if !Ping() {
		t.Fatal()
	}

	obj := &TMPObj{A: 10, B: 20, C: "test"}
	if !Send(0, 1, 1, obj) {
		t.Error("send")
	}

	if Login(0) {
		t.Error("login")
	}

	info, success := GetInfo(0)
	fmt.Println(info)
	fmt.Println(success)
}
