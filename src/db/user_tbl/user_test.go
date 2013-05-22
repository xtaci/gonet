package user_tbl

import "testing"
import "fmt"

func TestUser(t *testing.T) {
	basic := New("xtaci", "qwer1234")
	fmt.Println(basic)

	all := GetAll()

	for _, v:= range all {
		fmt.Println(v)
	}
}
