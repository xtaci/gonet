package ipc

import (
	"testing"
	"fmt"
)

func TestEventFunc(t *testing.T) {
	DialHub()
	fmt.Println(Ping())
}
