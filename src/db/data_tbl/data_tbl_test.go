package data_tbl

import (
	"fmt"
	"misc/naming"
	"testing"
	"types/estates"
)

func TestEstate(t *testing.T) {
	data := &estates.Manager{}
	e1 := &estates.Estate{TYPE: naming.FNV1a("工人小屋"), Ready: true}
	data.Append(100, e1)
	data.UserId = 1

	fmt.Println("Set")
	Set(estates.COLLECTION, data)

	var value *estates.Manager
	Get(estates.COLLECTION, 1, &value)
	fmt.Println("Get:", value)

	var all []estates.Manager
	GetAll(estates.COLLECTION, &all)
	fmt.Println("GetAll")
	fmt.Println(value)
}
