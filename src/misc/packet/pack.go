package packet

import (
	"reflect"
)

//----------------------------------------------- export struct fields with packet writer.
func Pack(tos int16, tbl interface{}, writer *Packet) []byte {
	if writer == nil {
		writer = Writer()
	}

	// write code
	if tos != -1 {
		writer.WriteU16(uint16(tos))
	}

	// is nil?
	v := reflect.ValueOf(tbl)
	if !v.IsValid() {
		return writer.Data()
	}

	// deal with pointers
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		v = v.Elem()
	}
	count := v.NumField()

	for i := 0; i < count; i++ {
		f := v.Field(i)
		switch f.Type().Kind() {
		case reflect.Slice, reflect.Array:
			writer.WriteU16(uint16(f.Len()))
			for a := 0; a < f.Len(); a++ {
				if _is_primitive(f.Index(a)) {
					_write_primitive(f.Index(a), writer)
				} else {
					elem := f.Index(a).Interface()
					Pack(-1, elem, writer)
				}
			}
		case reflect.Struct:
			Pack(-1, f.Interface(), writer)
		default:
			_write_primitive(f, writer)
		}
	}

	return writer.Data()
}

//----------------------------------------------- test whether the field is primitive type
func _is_primitive(f reflect.Value) bool {
	switch f.Type().Kind() {
	case reflect.Bool,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Float32,
		reflect.Float64,
		reflect.String:
		return true
	}
	return false
}

//----------------------------------------------- write a primitive field
func _write_primitive(f reflect.Value, writer *Packet) {
	switch f.Type().Kind() {
	case reflect.Bool:
		writer.WriteBool(f.Interface().(bool))
	case reflect.Uint8:
		writer.WriteByte(f.Interface().(byte))
	case reflect.Uint16:
		writer.WriteU16(f.Interface().(uint16))
	case reflect.Uint32:
		writer.WriteU32(f.Interface().(uint32))
	case reflect.Uint64:
		writer.WriteU64(f.Interface().(uint64))

	case reflect.Int:
		writer.WriteU32(uint32(f.Interface().(int)))
	case reflect.Int8:
		writer.WriteByte(byte(f.Interface().(int8)))
	case reflect.Int16:
		writer.WriteU16(uint16(f.Interface().(int16)))
	case reflect.Int32:
		writer.WriteU32(uint32(f.Interface().(int32)))
	case reflect.Int64:
		writer.WriteU64(uint64(f.Interface().(int64)))

	case reflect.Float32:
		writer.WriteFloat32(f.Interface().(float32))

	case reflect.Float64:
		writer.WriteFloat64(f.Interface().(float64))

	case reflect.String:
		writer.WriteString(f.Interface().(string))
	}
}
