package event_tbl

import (
	"fmt"
	"testing"
)

func TestEvent(t *testing.T) {
	params := []byte("hello")
	Push(1, 2, 3, params)
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
