package clan

import (
	"fmt"
	"testing"
)

func TestClan(t *testing.T) {
	clanid, succ := Create(1, "tap4fun")
	if !succ {
		t.Error("cannot create clan")
	}
	fmt.Println(clanid, succ)
	if !Join(2, clanid) {
		t.Error("cannot join clan")
	}

	if !Leave(2, clanid) {
		t.Error("cannot leave join")
	}
}
