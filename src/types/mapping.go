package types

import (
	"types/estates"
)

var Mapping map[string]interface{} = map[string]interface{} {
	"USER":User{},
	"B_BASIC":estates.BasManager{},
	"B_DEFENSIVE":estates.DefManager{},
	"B_OFFENSIVE":estates.OffManager{},
	"B_RESOURCE":estates.ResManager{},
}
