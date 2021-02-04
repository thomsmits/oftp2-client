package client

import (
	"fmt"
	"net"
	"strings"

	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire"
)

// OFTP2Client represents a communication facility to speak OFTP2 with a server.
type OFTP2Client struct {
	// ServerHost is the hostname of the server to connect to
	ServerHost string

	// ServerPort is the port of the server to connect to
	ServerPort int

	// OdetteId is the clients Odette ID (see GenerateOdetteId())
	OdetteId string

	// Verbose sets verbose output during communication with the server
	Verbose bool

	// Fuzzer is a function that gets each data package before it is sent to the
	// server and can perform changes on it. Please note that the data shown to the
	// fuzzer in data is the raw data that goes on the wire. Therefore, it contains
	// all OFTP2 TCP headers, all sub record headers and so on. Simply writing to a
	// location may invalidate the data from the OFTP2 protocol point of view, which
	// may (or may not) be what you want.
	Fuzzer                        func(data []byte) []byte
	con                           *net.Conn // Network connection
	serverId                      string    // Odette ID of the server we are talking to
	serverPassword                string    // Server password
	serverBufferSize              uint32    // Server's maximum buffer size
	serverCapability              string    // Capabilities of the server
	serverCompress                bool      // Server supports compression
	serverRestartSupported        bool      // Server supports restart
	serverSpecial                 bool      // Server supports special commands
	serverCredit                  uint32    // Number of data buffers, server accepts before CDT command
	serverAuthenticationSupported bool      // Server supports authentication
	serverUserData                string    // User data string send by the server
}

// OFTP2FileFormat specifies the file formats supported by the protocol
type OFTP2FileFormat string

const (
	// FileFormatFixedBinary is a binary file with fixed record structure
	FileFormatFixedBinary OFTP2FileFormat = "F"

	// FileFormatVariable is a file with a variable structure
	FileFormatVariable OFTP2FileFormat = "V"

	// FileFormatUnstructured is an unstructured file without any substructure
	FileFormatUnstructured OFTP2FileFormat = "U"

	// FileFormatText is a text file
	FileFormatText OFTP2FileFormat = "T"
)

// OFTP2SecurityLevel Security Levels
type OFTP2SecurityLevel int

const (
	// SecurityLevelNone indicates no security level
	SecurityLevelNone OFTP2SecurityLevel = iota

	// SecurityLevelEncrypted indicates an encrypted file
	SecurityLevelEncrypted

	// SecurityLevelSigned indicates a signed file
	SecurityLevelSigned

	// SecurityLevelSignedAndEncrypted indicates a file which is encrypted and signed
	SecurityLevelSignedAndEncrypted
)

func (o OFTP2FileFormat) String() string {
	return string(o)
}

// GenerateOdetteId generates a RFC-compliant ODETTE id for the given
// international code, organization code and computer address
func GenerateOdetteId(intCode int, orgCode, subAddress string) string {
	return strings.Trim(fmt.Sprintf("O%04d%s%s", intCode, wire.TruncateAndPadString(orgCode, 14), wire.TruncateAndPadString(subAddress, 6)), " ")
}
