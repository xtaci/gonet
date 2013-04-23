package protos

import "misc/packet"
import "reflect"

func pack(tbl interface{}) []byte {
	writer := packet.PacketWriter()
	v := reflect.ValueOf(tbl).Elem()
	key := v.Type()
	count := key.NumField()

	// write code
	writer.WriteU16(Code[reflect.TypeOf(tbl).Name()])

	for i := 0; i < count; i++ {
		f := v.Field(i)
		if f.CanSet() {
			switch f.Type().String() {
			case "byte", "uint8":
				writer.WriteByte(f.Interface().(byte))
			case "uint16":
				writer.WriteU16(f.Interface().(uint16))
			case "uint32":
				writer.WriteU32(f.Interface().(uint32))
			case "uint64":
				writer.WriteU64(f.Interface().(uint64))

			case "int8":
				writer.WriteByte(byte(f.Interface().(int8)))
			case "int16":
				writer.WriteU16(uint16(f.Interface().(int16)))
			case "int32":
				writer.WriteU32(uint32(f.Interface().(int32)))
			case "int64":
				writer.WriteU64(uint64(f.Interface().(int64)))

			case "float32":
				writer.WriteFloat32(f.Interface().(float32))

			case "string":
				writer.WriteString(f.Interface().(string))
			}
		}
	}

	return nil
}
