package event_client

import (
	"fmt"
	"testing"
)

func TestEventFunc(t *testing.T) {
	DialEvent()
	fmt.Println(Ping())
}
