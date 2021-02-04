package wire

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

// Command represents a single OFTP2 command send over the wire
type Command struct {
	Format []DataFormat
	Data   []interface{}
}

// Marshal marshals the command into a byte sequence
func (c *Command) Marshal() []byte {

	formatDefs := c.Format
	dataEntries := c.Data

	if len(formatDefs) != len(dataEntries) {
		panic(fmt.Sprintf("Data definition and data have different length: %d != %d", len(formatDefs), len(dataEntries)))
	}

	result := make([]byte, 0)

	for i, format := range formatDefs {
		data := dataEntries[i]
		var truncated string
		var bytes []byte

		switch format.DataType {
		case DataTypeAlphanumeric:
			switch data.(type) {
			case string:
				s := data.(string)
				truncated = TruncateAndPadString(s, format.Length)
			case bool:
				if data.(bool) {
					truncated = "Y"
				} else {
					truncated = "N"
				}
			case fmt.Stringer:
				s := data.(fmt.Stringer)
				truncated = TruncateAndPadString(s.String(), format.Length)

			default:
				panic(fmt.Sprintf("data type does not match with definition. Expecting alphanumeric, got %s", reflect.TypeOf(data).String()))
			}

			// TODO: Spec says space is not allowed as an embedded character but uses it in the SSRM message
			check := regexp.MustCompile(`^[0-9A-Z/\-.&() ]*\s*$`)

			if !check.MatchString(truncated) {
				//  panic(fmt.Sprintf("wrong format for alphanumeric string"))
			}

			bytes = []byte(truncated)

		case DataTypeNumeric:
			switch data.(type) {
			case uint16:
				formatString := fmt.Sprintf("%%0%dd", format.Length)
				formatted := fmt.Sprintf(formatString, data.(uint16))
				truncated = TruncateString(formatted, format.Length)

			case int:
				formatString := fmt.Sprintf("%%0%dd", format.Length)
				formatted := fmt.Sprintf(formatString, data.(int))
				truncated = TruncateString(formatted, format.Length)

			case uint32:
				formatString := fmt.Sprintf("%%0%dd", format.Length)
				formatted := fmt.Sprintf(formatString, data.(uint32))
				truncated = TruncateString(formatted, format.Length)

			case uint64:
				formatString := fmt.Sprintf("%%0%dd", format.Length)
				formatted := fmt.Sprintf(formatString, data.(uint64))
				truncated = TruncateString(formatted, format.Length)

			case string:
				s := data.(string)
				i, err := strconv.Atoi(s)
				if err != nil {
					panic(fmt.Sprintf("Cannot convert %s to int", s))
				}

				formatString := fmt.Sprintf("%%0%dd", format.Length)
				formatted := fmt.Sprintf(formatString, i)
				truncated = TruncateString(formatted, format.Length)

			default:
				panic("Data not numeric")
			}

			bytes = []byte(truncated)

		case DataTypeBinary:
			var v uint64
			var numeric = true

			switch data.(type) {
			case int:
				v = uint64(data.(int))

			case uint16:
				v = uint64(data.(uint16))

			case uint32:
				v = uint64(data.(uint32))

			case uint64:
				v = data.(uint64)

			case []byte:
				bytes = data.([]byte)
				numeric = false
			}

			if numeric {
				switch format.Length {
				case 2: // 16 Bit
					s := uint16(v)
					b := make([]byte, 2)
					binary.BigEndian.PutUint16(b, s)

					bytes = b

				case 4: // 32 Bit
					s := uint32(v)
					b := make([]byte, 4)
					binary.BigEndian.PutUint32(b, s)

					bytes = b

				case 8: // 64 Bit
					s := uint64(v)
					b := make([]byte, 8)
					binary.BigEndian.PutUint64(b, s)

					bytes = b
				}
			}

		case DataTypeUTF8:
			switch data.(type) {
			case string:
				s := data.(string)
				truncated = TruncateAndPadString(s, format.Length)
			case fmt.Stringer:
				s := data.(fmt.Stringer)
				truncated = TruncateAndPadString(s.String(), format.Length)
			default:
				panic(fmt.Sprintf("data type does not match with definition. Expecting UTF8, got %s", reflect.TypeOf(data).String()))
			}

			bytes = []byte(truncated)

		default:
			panic(fmt.Sprintf("Unknown data type %d", format.DataType))
		}

		result = append(result, bytes...)
	}

	return result
}
