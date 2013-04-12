package db

import "fmt"
import "reflect"

func sql_dump(tbl interface{})(fields []string, values []string){
	v := reflect.ValueOf(tbl).Elem()
	key := v.Type()
	count := key.NumField()

	fields = make([]string, count)
	values = make([]string, count)

	slice_idx := 0
	for i := 0; i < count; i++ {
		typeok := true
		switch v.Field(i).Kind() {
		case reflect.Int:
			values[slice_idx] = fmt.Sprintf("'%d'",v.Field(i).Int())
		case reflect.String:
			values[slice_idx] = fmt.Sprintf("'%s'",v.Field(i).String())
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
