package user_tbl

import (
	. "db"
	. "types"
)

import (
	"fmt"
	"strings"
	"errors"
	"encoding/base64"
)

const (
	FMT = "%v:%v:%v:%v"
)


func Load(user_id int32) (list []Building, bitmap *Map, err error) {
	stmt := "SELECT list, map FROM buildings where user_id ='%v' LIMIT 1"

	db := <-DBCH
	defer func() { DBCH <- db }()

	rows, _, err := db.Query(stmt, user_id)
	CheckErr(err)

	if len(rows) > 0 {
		list = _unpack(rows[0].Str(0))
		bitmap = _decode_map(rows[0].Str(1))
	}

	err = errors.New(fmt.Sprint("cannot find building belongs to id:%v", user_id))
	return

}

func Store(user_id int32, list []Building, bitmap *Map) {
	str_list := _pack(list)
	str_bitmap := _encode_map(bitmap.Bitset)

	stmt := "UPDATE buildings SET list='%v', bitmap='%v' WHERE user_id = %v"

	db := <-DBCH
	defer func() { DBCH <- db }()
	_, _, err := db.Query(stmt, str_list, str_bitmap, user_id)

	CheckErr(err)
}

func _pack(list []Building) string {
	var build_strs []string
	for k := range list {
		b := &list[k]
		build_strs = append(build_strs, fmt.Sprintf(FMT, b.TYPE, b.X,b.Y,b.LV))
	}

	return strings.Join(build_strs,"#")
}

func _unpack(list_str string) []Building {
	list := strings.Split(list_str,"#")
	var buildings []Building
	for k := range list {
		b := &Building{}
		fmt.Sscanf(list[k], FMT, &b.TYPE, &b.X, &b.Y, &b.LV)
		buildings = append(buildings, *b)
	}

	return buildings
}

func _decode_map(mapstr string) *Map {
	bitmap := &Map{}
	mapdata, err := base64.StdEncoding.DecodeString(mapstr)
	bitmap.Bitset = mapdata

	if err != nil {
		return nil
	}

	return bitmap
}

func _encode_map(bitmap []byte) string {
	return base64.StdEncoding.EncodeToString(bitmap)
}
