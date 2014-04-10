package gamedata

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

import (
	"cfg"
	. "helper"
)

var _lock sync.RWMutex
var _tables map[string]*Table

//---------------------------------------------------------- info for a level
type Record struct {
	Fields map[string]string
}

//---------------------------------------------------------- Numerical Table for a object
type Table struct {
	Records map[string]*Record
}

func init() {
	Reload()
}

//----------------------------------------------------------- Reload *.csv
func Reload() {
	_lock.Lock()
	defer _lock.Unlock()

	_tables = make(map[string]*Table)

	pattern := os.Getenv("GOPATH") + "/src/gamedata/data/*.csv"

	config := cfg.Get()
	if config["gamedata_dir"] != "" {
		pattern = config["gamedata_dir"] + "/*.csv"
	}

	INFO("Loading GameData From", pattern)
	files, err := filepath.Glob(pattern)

	if err != nil {
		ERR(err)
		panic(err)
	}

	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			ERR("error opening file", err)
			continue
		}

		parse(file)
		file.Close()
	}

	log.Printf("\033[042;1m%v CSV(s) Loaded\033[0m\n", len(_tables))
}

//---------------------------------------------------------- Set Field value
func _set(tblname string, rowname string, fieldname string, value string) {
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

//---------------------------------------------------------- Get Field value
func _get(tblname string, rowname string, fieldname string) string {
	_lock.RLock()
	defer _lock.RUnlock()

	tbl, ok := _tables[tblname]
	if !ok {
		panic(fmt.Sprint("table ", tblname, " not exists!"))
	}

	rec, ok := tbl.Records[rowname]
	if !ok {
		panic(fmt.Sprint("table ", tblname, " row ", rowname, " not exists!"))
	}

	value, ok := rec.Fields[fieldname]
	if !ok {
		panic(fmt.Sprint("table ", tblname, " field ", fieldname, " not exists!"))
	}
	return value
}

//---------------------------------------------------------- Get Field value as Integer
func GetInt(tblname string, rowname string, fieldname string) int32 {
	val := _get(tblname, rowname, fieldname)
	v, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Sprintf("cannot parse integer from gamedata %v %v %v %v\n", tblname, rowname, fieldname, err))
	}

	return int32(v)
}

//---------------------------------------------------------- Get Field value as Float
func GetFloat(tblname string, rowname string, fieldname string) float64 {
	val := _get(tblname, rowname, fieldname)
	if val == "" {
		return 0.0
	}

	f, err := strconv.ParseFloat(val, 32)
	if err != nil {
		panic(fmt.Sprintf("cannot parse float from gamedata %v %v %v %v\n", tblname, rowname, fieldname, err))
	}

	return f
}

//---------------------------------------------------------- Get Field value as string
func GetString(tblname string, rowname string, fieldname string) string {
	return _get(tblname, rowname, fieldname)
}

//---------------------------------------------------------- Get Row Count
func Count(tblname string) int32 {
	tbl := _tables[tblname]

	if tbl == nil {
		return 0
	}

	return int32(len(tbl.Records))
}

//---------------------------------------------------------- Test Field Exists
func IsFieldExists(tblname string, fieldname string) bool {
	_lock.RLock()
	defer _lock.RUnlock()

	tbl := _tables[tblname]

	if tbl == nil {
		return false
	}

	key := ""
	// get one record key
	for k, _ := range tbl.Records {
		key = k
		break
	}

	rec, ok := tbl.Records[key]
	if !ok {
		return false
	}

	_, ok = rec.Fields[fieldname]
	if !ok {
		return false
	}

	return true
}

//---------------------------------------------------------- Load JSON From GameData Directory
func LoadJSON(filename string) ([]byte, error) {
	prefix := os.Getenv("GOPATH") + "/src/gamedata/data"
	config := cfg.Get()
	if config["gamedata_dir"] != "" {
		prefix = config["gamedata_dir"]
	}

	path := prefix + "/" + filename
	return ioutil.ReadFile(path)
}
