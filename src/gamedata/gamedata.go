package gamedata

import (
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

import (
	"misc/naming"
)

var _lock sync.RWMutex
var _tables map[uint32]*Table
var _hashtbl map[uint32]string // hash->string

//---------------------------------------------------------- info for a level
type Record struct {
	Fields map[uint32]string
}

//---------------------------------------------------------- Numerical Table for a object
type Table struct {
	Records map[uint32]*Record
}

func init() {
	Reload()
}

//----------------------------------------------------------- Reload *.csv
func Reload() {
	_lock.Lock()
	defer _lock.Unlock()

	log.Println("Loading GameData...")
	defer log.Println("GameData Loaded.")

	_tables = make(map[uint32]*Table)
	_hashtbl = make(map[uint32]string)

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

//----------------------------------------------------------- Query a name by hash
func Query(hash uint32) string {
	return _hashtbl[hash]
}

//---------------------------------------------------------- Set Field value
func _set(tblname string, rowname string, fieldname string, value string) {
	// store hashing
	h_rowname := naming.FNV1a(rowname)
	h_fieldname := naming.FNV1a(fieldname)
	h_tblname := naming.FNV1a(tblname)
	_hashtbl[h_rowname] = rowname
	_hashtbl[h_fieldname] = fieldname
	_hashtbl[h_tblname] = tblname

	//
	tbl := _tables[h_tblname]

	if tbl == nil {
		tbl = &Table{}
		tbl.Records = make(map[uint32]*Record)
		_tables[h_tblname] = tbl
	}

	rec := tbl.Records[h_rowname]
	if rec == nil {
		rec = &Record{}
		rec.Fields = make(map[uint32]string)
		tbl.Records[h_rowname] = rec
	}

	rec.Fields[h_fieldname] = value
}

//---------------------------------------------------------- Get Field value
func _gethash(h_tblname uint32, h_rowname uint32, h_fieldname uint32) string {
	_lock.RLock()
	defer _lock.RUnlock()

	tbl := _tables[h_tblname]

	if tbl == nil {
		return ""
	}

	rec := tbl.Records[h_rowname]
	if rec == nil {
		return ""
	}

	return rec.Fields[h_fieldname]
}

func _get(tblname string, rowname string, fieldname string) string {
	return _gethash(naming.FNV1a(tblname), naming.FNV1a(rowname), naming.FNV1a(fieldname))
}

//---------------------------------------------------------- Get Field value as Integer
func GetInt(tblname string, rowname string, fieldname string) int32 {
	val := _get(tblname, rowname, fieldname)
	v, err := strconv.Atoi(val)
	if err != nil {
		v = math.MaxInt32
		log.Println("cannot parse integer from gamedata", err)
	}

	return int32(v)
}

//---------------------------------------------------------- Get Field value as Float
func GetFloat(tblname string, rowname string, fieldname string) float32 {
	val := _get(tblname, rowname, fieldname)
	f, err := strconv.ParseFloat(val, 32)
	if err != nil {
		f = math.MaxFloat32
		log.Println("cannot parse float from gamedata", err)
	}

	return float32(f)
}

//---------------------------------------------------------- Get Field value as string
func GetString(tblname string, rowname string, fieldname string) string {
	return _get(tblname, rowname, fieldname)
}
