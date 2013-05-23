package gamedata

import (
	"strconv"
	"log"
)

//----------------------------------------------- info for a level
type Record struct {
	Fields map[string]string
}

//----------------------------------------------- Numerical Table for a object
type Table struct {
	Records map[int]*Record
}

var _tables map[string]*Table

func init() {
	_tables = make(map[string]*Table)
}

//----------------------------------------------- Set Field value
func Set(tblname string, level int, fieldname string, value string) {
	tbl := _tables[tblname]

	if tbl == nil {
		tbl = &Table{}
		tbl.Records = make(map[int]*Record)
		_tables[tblname] = tbl
	}

	rec := tbl.Records[level]
	if rec == nil {
		rec = &Record{}
		rec.Fields = make(map[string]string)
		tbl.Records[level] = rec
	}

	rec.Fields[fieldname] = value
}

//----------------------------------------------- Get Field value
func _get(tblname string, level int, fieldname string) string {
	tbl := _tables[tblname]

	if tbl == nil {
		return ""
	}

	rec := tbl.Records[level]
	if rec == nil {
		return ""
	}

	return rec.Fields[fieldname]
}


func GetInt(tblname string, level int, fieldname string) int32 {
	val := _get(tblname, level, fieldname)
	if val == "" {
		return ^int32(0)		// return MAX INT
	}

	v, _ := strconv.Atoi(val)

	return int32(v)
}

func GetFloat(tblname string, level int, fieldname string) float32 {
	val := _get(tblname, level, fieldname)
	if val == "" {
		return 0.0
	}

	f, err :=  strconv.ParseFloat(val, 32)
	if err!= nil {
		log.Println(GetFloat, err)
	}

	return float32(f)
}

func FieldNames(tblname string) []string {
	tbl := _tables[tblname]

	if tbl == nil {
		return nil
	}

	rec := tbl.Records[1]
	if rec == nil {
		return nil
	}

	ret := make([]string, len(rec.Fields))
	count :=0
	for k := range rec.Fields {
		ret[count] = k
		count++
	}

	return ret
}

//------------------------------------------------ Get Num of Levels
func NumLevels(tblname string) int {
	tbl := _tables[tblname]

	if tbl == nil {
		return 0
	}

	return len(tbl.Records)
}
