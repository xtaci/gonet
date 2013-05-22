package estate_tbl

import "testing"
import "types/estate"
import "fmt"

func TestEstate(t *testing.T) {
	data := &estate.EstateManager{}
	Set(0, data)
	data = Get(0)
	fmt.Println(data)
}
