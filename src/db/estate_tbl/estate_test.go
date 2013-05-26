package estate_tbl

import "testing"
import "types/estate"
import "fmt"

func TestEstate(t *testing.T) {
	data := &estate.Manager{}
	e1 := &estate.Estate{}
	data.AppendEstate(e1)
	cd1  := &estate.CD{}
	data.AppendCD(1, cd1)
	data.Id = 1
	Set(data)
	data = Get(1)
	fmt.Println(data)
	for k:= range data.CDs {
		fmt.Println(data.CDs[k])
	}
}
