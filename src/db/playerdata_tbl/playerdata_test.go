package playerdata_tbl

import "testing"
import . "types"

func TestPlayerData(t *testing.T) {
	data := &PlayerData{}
	Set(0, data)
	Get(0, data)
}
