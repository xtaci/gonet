package defensive_tbl

import "testing"
import "types/defensive"
import "fmt"

func TestEstate(t *testing.T) {
	data := &defensive.Manager{}
	e1 := &defensive.Defensive{}
	data.AppendDefensive(e1)
	cd1 := &defensive.CD{}
	data.AppendCD(1, cd1)
	data.Id = 1
	Set(data)
	data = Get(1)
	fmt.Println(data)
	for k := range data.CDs {
		fmt.Println(data.CDs[k])
	}
}
