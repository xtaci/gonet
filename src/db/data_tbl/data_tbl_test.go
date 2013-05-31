package data_tbl

import (
	"fmt"
	"testing"
	"time"
	"types/estates"
	"misc/naming"
)

func TestEstate(t *testing.T) {
	data := &estates.Manager{}
	e1 := &estates.Estate{TYPE:naming.FNV1a("工人小屋")}
	data.Append(1024, e1)
	cd1 := &estates.CD{OID: 100, Timeout: time.Now().Unix()}
	data.AppendCD(1, cd1)
	data.UserId = 1

	Set(estates.COLLECTION, data)

	value := &estates.Manager{}
	Get(estates.COLLECTION, 1, value)
	fmt.Println("VALUE:", value)
	for k := range value.CDs {
		fmt.Println("CD:", value.CDs[k])
	}
}
