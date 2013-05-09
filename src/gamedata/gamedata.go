package gamedata

//----------------------------------------------- info for a level
type Record struct {
	Fields map[string]interface{}
}

//----------------------------------------------- Building Upgrade Info 
type Table struct {
	Records map[int]*Record
}

var _tables map[string]*Table

func init() {
	_tables = make(map[string]*Table)
}

//----------------------------------------------- Set Field value
func Set(tblname string, row int, fieldname string, value interface{}) {
	tbl := _tables[tblname]

	if tbl == nil {
		tbl = &Table{}
		tbl.Records = make(map[int]*Record)
		_tables[tblname] = tbl
	}

	rec := tbl.Records[row]
	if rec == nil {
		rec = &Record{}
		rec.Fields = make(map[string]interface{})
		tbl.Records[row] = rec
	}

	rec.Fields[fieldname] = value
}

//----------------------------------------------- Get Field value
func Get(tblname string, row int, fieldname string) interface{} {
	tbl := _tables[tblname]

	if tbl == nil {
		return nil
	}

	rec := tbl.Records[row]
	if rec == nil {
		return nil
	}

	return rec.Fields[fieldname]
}
