package report

import (
	"encoding"
	"encoding/binary"
	"reflect"
	"strconv"

	"github.com/CursedHardware/go-nsdp"
)

func UnmarshalReport(message *nsdp.Message, report any) error {
	v := reflect.ValueOf(report).Elem()
	t := v.Type()
	var value []byte
	typeByteArray := reflect.TypeOf(([]byte)(nil))
	typeBinary := reflect.TypeOf((encoding.BinaryUnmarshaler)(nil))
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if parsed, err := strconv.ParseUint(field.Tag.Get("nsdp-scan"), 16, 16); err != nil {
			return err
		} else if _, ok := message.Tags[nsdp.Tag(parsed)]; !ok {
			continue
		} else {
			value = message.Tags[nsdp.Tag(parsed)]
		}
		switch fieldType := field.Type; fieldType.Kind() {
		case reflect.String:
			v.Field(i).SetString(string(value))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			parsed, _ := binary.Varint(value)
			v.Field(i).SetInt(parsed)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			parsed, _ := binary.Uvarint(value)
			v.Field(i).SetUint(parsed)
		case reflect.Bool:
			v.Field(i).SetBool(value[0] != 0)
		default:
			switch {
			case fieldType.ConvertibleTo(typeByteArray):
				v.Field(i).SetBytes(value)
			case fieldType.ConvertibleTo(typeBinary):
				_ = v.Field(i).Interface().(encoding.BinaryUnmarshaler).UnmarshalBinary(value)
			}
		}
	}
	return nil
}
