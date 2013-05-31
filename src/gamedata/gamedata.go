package gamedata

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var _tables map[string]*Table

//----------------------------------------------- info for a level
type Record struct {
	Fields map[string]string
}

//----------------------------------------------- Numerical Table for a object
type Table struct {
	Records map[string]*Record
}

func init() {
	_tables = make(map[string]*Table)

	pattern := os.Getenv("GOPATH") + "/src/gamedata/data/*.csv"
	files, err := filepath.Glob(pattern)

	if err != nil {
		log.Println(err)
		panic(err)
	}

	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			log.Println("error opening file %v\n", err)
			continue
		}

		parse(file)
		file.Close()
	}
}

//----------------------------------------------- Set Field value
func Set(tblname string, rowname string, fieldname string, value string) {
	tbl := _tables[tblname]

	if tbl == nil {
		tbl = &Table{}
		tbl.Records = make(map[string]*Record)
		_tables[tblname] = tbl
	}

	rec := tbl.Records[rowname]
	if rec == nil {
		rec = &Record{}
		rec.Fields = make(map[string]string)
		tbl.Records[rowname] = rec
	}

	rec.Fields[fieldname] = value
}

//----------------------------------------------- Get Field value
func _get(tblname string, rowname string, fieldname string) string {
	tbl := _tables[tblname]

	if tbl == nil {
		return ""
	}

	rec := tbl.Records[rowname]
	if rec == nil {
		return ""
	}

	return rec.Fields[fieldname]
}

func GetInt(tblname string, rowname string, fieldname string) int32 {
	val := _get(tblname, rowname, fieldname)
	if val == "" {
		return ^int32(0) // return MAX INT
	}

	v, _ := strconv.Atoi(val)

	return int32(v)
}

func GetFloat(tblname string, rowname string, fieldname string) float32 {
	val := _get(tblname, rowname, fieldname)
	if val == "" {
		return 0.0
	}

	f, err := strconv.ParseFloat(val, 32)
	if err != nil {
		log.Println(GetFloat, err)
	}

	return float32(f)
}

//----------------------------------------------- Get Field value as string
func GetString(tblname string, rowname string, fieldname string) string {
	return _get(tblname, rowname, fieldname)
}
