package db

import "fmt"
import "reflect"
import "regexp"
import "github.com/ziutek/mymysql/mysql"

func sql_dump(tbl interface{})(fields []string, values []string) {
	re := regexp.MustCompile(`(\'|\"|\.|\*|\/|\-|\\)`)

	v := reflect.ValueOf(tbl).Elem()
	key := v.Type()
	count := key.NumField()

	fields = make([]string, count)
	values = make([]string, count)

	slice_idx := 0
	for i := 0; i < count; i++ {
		typeok := true
		switch v.Field(i).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			values[slice_idx] = fmt.Sprintf("'%d'",v.Field(i).Int())
		case reflect.Uint,reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			values[slice_idx] = fmt.Sprintf("'%d'",v.Field(i).Uint())
		case reflect.Float32, reflect.Float64:
			values[slice_idx] = fmt.Sprintf("'%f'", v.Field(i).Float())
		case reflect.String:
			tmpstr := re.ReplaceAllString(v.Field(i).String(), `\${1}`)
			values[slice_idx] = fmt.Sprintf("'%s'",tmpstr)
		default:
			typeok = false
		}

		if (typeok) {
			fields[slice_idx] = underscore(key.Field(i).Name)
			slice_idx++
		}
	}

	fields = fields[:slice_idx]
	values = values[:slice_idx]

	return
}

func sql_load(tbl interface{}, row *mysql.Row, res mysql.Result) {
	v := reflect.ValueOf(tbl).Elem()
	for i, field := range res.Fields() {
		f := v.FieldByName(camelcase(field.Name))
		if f.IsValid() {
            if f.CanSet() {
				switch f.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					f.SetInt(int64(row.Int(i)))
				case reflect.Uint,reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					f.SetUint(uint64(row.Uint(i)))
				case reflect.Float32, reflect.Float64:
					f.SetFloat(row.Float(i))
				case reflect.String:
					f.SetString(row.Str(i))
				}
            }
		}
	}
}
