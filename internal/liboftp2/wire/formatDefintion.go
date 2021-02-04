package wire

import (
	"fmt"
	"regexp"
	"strconv"
)

// FormatDefinition provides a format definition for a field in an OFTP2 message.
// In contrast to DataFormat, FormatDefinition just wraps the textual definitions
// from the specification and does not parse the data into any machine readable format.
type FormatDefinition struct {
	Mnemonic       string
	FieldName      string
	PossibleValues *[]Value
}

// FormatDefinitionsToDataFormats transforms a list of FormatDefinition objects into a list of DataFormat objects
func FormatDefinitionsToDataFormats(input []FormatDefinition) []DataFormat {
	result := make([]DataFormat, 0)

	for _, f := range input {
		df := f.ToDataFormat()

		if f.PossibleValues != nil {
			df.PossibleValues = f.PossibleValues
		}

		result = append(result, df)
	}

	return result
}

// ToDataFormat transforms this object to a data format object
func (f *FormatDefinition) ToDataFormat() DataFormat {

	r := regexp.MustCompile(`([FV]) ([X9UT])\(([0-9]+|n)\)`)
	matches := r.FindStringSubmatch(f.Mnemonic)

	if len(matches) != 4 {
		panic(fmt.Sprintf("Illegal type format %s specified", f.Mnemonic))
	}

	typeIndicator := matches[1]
	formatIndicator := matches[2]
	lengthIndicator := matches[3]

	result := DataFormat{}

	switch typeIndicator {
	case "F": // A field containing fixed values
		result.Fixed = true
	case "V": // A field with variable values within a defined range.
		result.Fixed = false
	default:
		panic(fmt.Sprintf("Illegal type indicator %s specified", typeIndicator))
	}

	switch formatIndicator {
	case "X": // An alphanumeric field with limited character set
		result.DataType = DataTypeAlphanumeric
	case "9": // A numeric field
		result.DataType = DataTypeNumeric
	case "U": // Binary data
		result.DataType = DataTypeBinary
	case "T": // UTF-8 string. Length is given in BYTES, not characters
		result.DataType = DataTypeUTF8
	default:
		panic(fmt.Sprintf("Illegal format indicator %s specified", formatIndicator))
	}

	if lengthIndicator == "n" {
		result.Length = -1
	} else {
		var err error
		result.Length, err = strconv.Atoi(lengthIndicator)
		if err != nil {
			panic(fmt.Sprintf("Illegal length indicator %s specified", lengthIndicator))
		}
	}

	if f.PossibleValues != nil {
		result.PossibleValues = f.PossibleValues
	}

	return result
}
