package ipc

import (
	"fmt"
	"testing"
)

func TestEventFunc(t *testing.T) {
	DialHub()
	fmt.Println(Ping())
}
