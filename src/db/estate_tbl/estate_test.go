package estate_tbl

import "testing"
import "types/estate"
import "fmt"

func TestEstate(t *testing.T) {
	data := &estate.Manager{}
	data.Estates = make([]estate.Estate, 2)
	data.Id = 1
	Set(data)
	data = Get(1)
	fmt.Println(data)
}
