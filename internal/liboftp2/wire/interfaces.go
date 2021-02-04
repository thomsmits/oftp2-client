package wire

// Protocol is the interface all OFTP2 commands (e.g. AUCH) have to implement.
type Protocol interface {

	// Parse parses the given byte array into the struct.
	Parse([]byte) error

	// Converts transforms the OFTP2 command into a general command structure. The
	// command structure can then be used to marshal the OFTP2 command into a byte
	// array.
	Command() Command

	// Marshal converts the OFTP2 command into a byte array.
	Marshal() []byte
}
