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
	Send(0, 1, 1, obj)
}
