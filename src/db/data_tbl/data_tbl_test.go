package data_tbl

import (
	"fmt"
	"testing"
	"time"
	"types/estates"
)

func TestEstate(t *testing.T) {
	data := &estates.Manager{}
	e1 := &estates.Estate{}
	data.Append(1024, e1)
	cd1 := &estates.CD{OID: 100, Timeout: time.Now().Unix()}
	data.AppendCD(1, cd1)
	data.UserId = 1

	Set("B_DEFENSIVE", data)

	value := &estates.Manager{}
	Get("B_DEFENSIVE", 1, value)
	fmt.Println("VALUE:", value)
	for k := range value.CDs {
		fmt.Println("CD:", value.CDs[k])
	}
}
