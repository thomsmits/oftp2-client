package wire

import "fmt"

// Enum for the different types of data types.
type DataTypes int

const (
	// An alphanumeric field with limited character set
	DataTypeAlphanumeric DataTypes = iota

	// A numeric field
	DataTypeNumeric

	// Binary data
	DataTypeBinary

	// UTF8 text
	DataTypeUTF8
)

// DataFormat represents a data format as it is described in the OFTP2
// specification. In contrast to FormatDefinition, DataFormat is a parsed,
// machine readable form of the data format specification from the RFC.
type DataFormat struct {
	Fixed          bool
	DataType       DataTypes
	Length         int
	PossibleValues *[]Value
}

// Value is a single value in the data format.
type Value struct {
	Name        string
	Description string
}

// IntMapToValues converts a normal go map with return codes as key and
// descriptions as value into an value array. To ensure the correct padding (e.g.
// 1 as "01") of the numeric values, the padding is specified.
func IntMapToValues(input map[int]string, padding int) *[]Value {
	result := make([]Value, 0)
	format := fmt.Sprintf("%%0%dd", padding)
	for key, description := range input {
		keyString := fmt.Sprintf(format, key)
		value := Value{
			Name:        keyString,
			Description: description,
		}
		result = append(result, value)
	}

	return &result
}

// StringMapToValues converts a normal go map with codes as key and descriptions
// as value into an value array.
func StringMapToValues(input map[string]string) *[]Value {
	result := make([]Value, 0)

	for key, description := range input {
		value := Value{
			Name:        key,
			Description: description,
		}
		result = append(result, value)
	}

	return &result
}

// Predefined value for a boolean "Y", "N" selection.
var ValueBooleanYesNo *[]Value = &[]Value{
	{"Y", "Yes"},
	{"N", "No"},
}

// Predefined value for the allowed newline characters from the specification.
var ValueNewline *[]Value = &[]Value{
	{"\x0D", "newline 0D"},
	{"\x8D", "newline 8D"},
}
