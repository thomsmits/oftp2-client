package wire

import (
	"fmt"
	"time"
)

// TruncateString truncates the given string to a maximum of length characters
func TruncateString(s string, length int) string {
	if length == -1 {
		return s
	} else if len(s) > length {
		return s[0:length]
	} else {
		return s
	}
}

// TruncateAndPadString truncates the given string to a maximum of length characters.
// If the string is shorter than length, the string is padded with spaces at the
// right up to the given length.
func TruncateAndPadString(s string, length int) string {
	if length == -1 {
		return s
	} else if len(s) > length {
		return s[0:length]
	} else if len(s) == length {
		return s
	} else {
		formatString := fmt.Sprintf("%%-%ds", length)
		return fmt.Sprintf(formatString, s)
	}
}

// ParseDateToString splits a Time object into the two strings for date and time
// required by the OFTP2 specification.
func ParseDateToString(dateToParse time.Time) (string, string) {
	d := dateToParse.Format("20060102")
	t := dateToParse.Format("150405")

	// get the milliseconds and first digit of the microseconds
	millis := dateToParse.UnixNano() % int64(time.Second) / (int64(time.Millisecond) / 10)

	// millis can become negative, therefore remove sign
	if millis < 0 {
		millis *= -1
	}

	t += fmt.Sprintf("%04d", millis) // specification requires additional 4 digits as counter, we simply use the millis

	return d, t
}

// ParseStringsToDate combines two strings, for date and time, into a normal Time object.
func ParseStringsToDate(d, t string) time.Time {
	toParse := d + t
	// strip away counter
	toParse = toParse[0:14]
	result, _ := time.Parse("20060102150405", toParse)
	return result
}
