package res_tbl

import "testing"
import "fmt"
import . "types"

func TestRes(t *testing.T) {
	res := &Res{Id:1}
	Set(res)
	data := Get(1)
	fmt.Println(data)
}
