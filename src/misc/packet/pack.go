package packet

import (
	"log"
	"reflect"
)

type FastPack interface {
	Pack(w *Packet)
}

//----------------------------------------------- export struct fields with packet writer.
func Pack(tos int16, tbl interface{}, writer *Packet) []byte {
	// create writer if not specified
	if writer == nil {
		writer = Writer()
	}

	// write protocol number
	if tos != -1 {
		writer.WriteU16(uint16(tos))
	}

	// is the table nil?
	if tbl == nil {
		return writer.Data()
	}

	// fastpack
	if fastpack, ok := tbl.(FastPack); ok {
		fastpack.Pack(writer)
		return writer.Data()
	}

	// pack by reflection
	_pack(reflect.ValueOf(tbl), writer)

	// return byte array
	return writer.Data()
}

//----------------------------------------------- export struct fields with packet writer.
func _pack(v reflect.Value, writer *Packet) {
	switch v.Kind() {
	case reflect.Bool:
		writer.WriteBool(v.Bool())
	case reflect.Uint8:
		writer.WriteByte(byte(v.Uint()))
	case reflect.Uint16:
		writer.WriteU16(uint16(v.Uint()))
	case reflect.Uint32:
		writer.WriteU32(uint32(v.Uint()))
	case reflect.Uint64:
		writer.WriteU64(uint64(v.Uint()))

	case reflect.Int16:
		writer.WriteS16(int16(v.Int()))
	case reflect.Int32:
		writer.WriteS32(int32(v.Int()))
	case reflect.Int64:
		writer.WriteS64(int64(v.Int()))

	case reflect.Float32:
		writer.WriteFloat32(float32(v.Float()))
	case reflect.Float64:
		writer.WriteFloat64(float64(v.Float()))

	case reflect.String:
		writer.WriteString(v.String())
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return
		}
		_pack(v.Elem(), writer)
	case reflect.Slice:
		l := v.Len()
		writer.WriteU16(uint16(l))
		for i := 0; i < l; i++ {
			_pack(v.Index(i), writer)
		}
	case reflect.Struct:
		numFields := v.NumField()
		for i := 0; i < numFields; i++ {
			_pack(v.Field(i), writer)
		}
	default:
		log.Println("cannot pack type:", v)
	}
}
