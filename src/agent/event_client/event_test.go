package event_client

import (
	"testing"
	"fmt"
)

func TestEventFunc(t *testing.T) {
	DialEvent()
	fmt.Println(Ping())
}
