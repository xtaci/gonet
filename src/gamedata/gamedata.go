package gamedata

//----------------------------------------------- info for a level
type Record struct {
	Fields map[string]interface{}
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
func Set(tblname string, level int, fieldname string, value interface{}) {
	tbl := _tables[tblname]

	if tbl == nil {
		tbl = &Table{}
		tbl.Records = make(map[int]*Record)
		_tables[tblname] = tbl
	}

	rec := tbl.Records[level]
	if rec == nil {
		rec = &Record{}
		rec.Fields = make(map[string]interface{})
		tbl.Records[level] = rec
	}

	rec.Fields[fieldname] = value
}

//----------------------------------------------- Get Field value
func Get(tblname string, level int, fieldname string) interface{} {
	tbl := _tables[tblname]

	if tbl == nil {
		return nil
	}

	rec := tbl.Records[level]
	if rec == nil {
		return nil
	}

	return rec.Fields[fieldname]
}

//------------------------------------------------ Get Num of Fields
func NumField(tblname string) int {
	tbl := _tables[tblname]

	if tbl == nil {
		return 0
	}
	rec := tbl.Records[0]
	if rec == nil {
		return 0
	}

	return len(tbl.Records[0].Fields)
}

//------------------------------------------------ Get Num of Levels
func NumLevel(tblname string) int {
	tbl := _tables[tblname]

	if tbl == nil {
		return 0
	}

	return len(tbl.Records)
}

//TODO : GetAsXXX 
