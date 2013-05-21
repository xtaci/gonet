package playerdata

import "testing"
import . "db"
import . "types"

func TestStore(t *testing.T) {
	StartDB()
	data := &PlayerData{}
	Store(0, data)
}
