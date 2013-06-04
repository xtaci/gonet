package user_tbl

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	basic := New(fmt.Sprintf("test%v", rand.Int()), fmt.Sprint(rand.Int()))
	fmt.Println("New:", basic)
	fmt.Println("Existing:")
	all := GetAll()
	for _, v := range all {
		fmt.Println(v)
	}
}
