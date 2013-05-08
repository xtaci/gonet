package types

import (
	"fmt"
	"strings"
)

type Building struct {
	TYPE uint8
	X    uint8
	Y    uint8
	LV   uint8
}

const (
	FMT = "%v:%v:%v:%v"
)

func Marshal(list []Building) string {
	var build_strs []string
	for k := range list {
		b := &list[k]
		build_strs = append(build_strs, fmt.Sprintf(FMT, b.TYPE, b.X,b.Y,b.LV))
	}

	return strings.Join(build_strs,"#")
}

func Unmarshal(list_str string) []Building {
	list := strings.Split(list_str,"#")
	var buildings []Building
	for k := range list {
		b := &Building{}
		fmt.Sscanf(list[k], FMT, &b.TYPE, &b.X, &b.Y, &b.LV)
		buildings = append(buildings, *b)
	}

	return buildings
}
