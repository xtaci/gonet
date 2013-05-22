package estate_tbl

import "testing"
import "types/estate"
import "fmt"

func TestEstate(t *testing.T) {
	data := &estate.Manager{}
	data.Estates = make([]estate.Estate, 2)
	Set(0, data)
	data = Get(0)
	fmt.Println(data)
}
