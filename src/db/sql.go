package db

import "fmt"
import "reflect"
import "regexp"
import "github.com/ziutek/mymysql/mysql"
import "time"
import "utils"

var escape_regexp * regexp.Regexp

func init() {
	escape_regexp = regexp.MustCompile(`(\'|\"|\.|\*|\/|\-|\\)`)
}

func SQL_escape(v string) string {
	return escape_regexp.ReplaceAllString(v, `\${1}`)
}

func SQL_dump(tbl interface{}, excludes ...string) (fields []string, values []string) {

	v := reflect.ValueOf(tbl).Elem()
	key := v.Type()
	count := key.NumField()

	fields = make([]string, count)
	values = make([]string, count)

	slice_idx := 0
	for i := 0; i < count; i++ {
		typeok := true
		switch v.Field(i).Type().String() {
		case "int", "int8", "int16","int32","int64":
			values[slice_idx] = fmt.Sprintf("'%d'", v.Field(i).Interface())
		case "uint", "uint8", "uint16","uint32","uint64":
			values[slice_idx] = fmt.Sprintf("'%d'", v.Field(i).Interface())
		case "float32", "float64":
			values[slice_idx] = fmt.Sprintf("'%f'", v.Field(i).Interface())
		case "string":
			tmpstr := SQL_escape(v.Field(i).Interface().(string))
			values[slice_idx] = fmt.Sprintf("'%s'", tmpstr)
		case "time.Time":
			values[slice_idx] = fmt.Sprintf("'%s'", v.Field(i).Interface().(time.Time).Format("2006-01-02 15:04:05"))
		default:
			typeok = false
		}

		if typeok {
			// kickout excluded fields
			fieldname := utils.UnderScore(key.Field(i).Name)
			for ei := range excludes {
				if excludes[ei] == fieldname {
					goto L
				}
			}
			fields[slice_idx] = fieldname
			slice_idx++
		}
	L:
	}

	fields = fields[:slice_idx]
	values = values[:slice_idx]

	return
}

func SQL_load(tbl interface{}, row *mysql.Row, res mysql.Result) {
	v := reflect.ValueOf(tbl).Elem()
	for i, field := range res.Fields() {
		f := v.FieldByName(utils.CamelCase(field.Name))
		if f.IsValid() {
			if f.CanSet() {
				switch f.Type().String() {
				case "int", "int8", "int16","int32","int64":
					f.SetInt(int64(row.Int(i)))
				case "uint", "uint8", "uint16","uint32","uint64":
					f.SetUint(uint64(row.Uint(i)))
				case "float32", "float64":
					f.SetFloat(row.Float(i))
				case "string":
					f.SetString(row.Str(i))
				case "time.Time":
					t,_ := time.Parse("2006-01-02 15:04:05", row.Str(i))
					f.Set(reflect.ValueOf(t))
				}
			}
		}
	}
}

func SQL_set_clause(fields []string, values []string) []string{
	changes := make([]string, len(fields))
	for i:= range fields {
		changes[i] = fields[i] + "=" + values[i]
	}

	return changes
}
