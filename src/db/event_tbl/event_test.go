package event_tbl

import (
	"fmt"
	"testing"
	"time"
	. "types"
)

func TestEvent(t *testing.T) {
	params := []byte("hello")
	event := &Event{EventId: 1, UserId: 2, Type: 3, Params: params, Timeout: time.Now().Unix()}
	Add(event)
	obj := Get(1)
	fmt.Println("#", obj, "#")
	if obj == nil {
		t.Fatal("cannot push events to db")
	}

	objs := GetAll()
	if len(objs) == 0 {
		t.Fatal("cannot get events from db")
	}

	fmt.Println("obj count:", len(objs))

	if !Remove(1) {
		t.Fatal("cannot remove event")
	}
}
