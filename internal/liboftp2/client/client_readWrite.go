package client

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// TCP_HEADER_LENGTH defines the number of bytes, OFTP2 uses in the TCP data
// stream to represent header information
const TCP_HEADER_LENGTH = 4

// oftp2TcpMagicWord is the magic word that marks the beginning of an OFTP2
// message via TCP
const oftp2TcpMagicWord = 0b0010000

// Read bytes from the connection and return the result as an byte array
func (s *OFTP2Client) read() ([]byte, error) {

	if s.con == nil {
		panic("Open connection first")
	}

	// Read Header with length indicator
	header := make([]byte, 4)
	_, err := (*s.con).Read(header)

	if err != nil {
		return nil, err
	}

	// Zero out fist byte, because it does not contribute to length information and
	// is just an OFTP2 specific bit sequence
	header[0] = 0

	// Get length of the rest of the message
	var length uint32
	length = binary.BigEndian.Uint32(header) - TCP_HEADER_LENGTH

	// Read remaining data from connection
	buff := make([]byte, length)
	_, err = (*s.con).Read(buff)
	if err != nil {
		return nil, err
	}

	if s.Verbose {
		fmt.Printf("<-- %s\n", strings.Trim(string(buff), "\n"))
	}

	return buff, nil
}

// Sends the given data to the connection, adding OFTP2 specific header
// information
func (s *OFTP2Client) write(input []byte) error {

	// The OFTP2 protocol requieres a 4 Byte header if with TCP as the transport protocol. The header consts of a

	if s.con == nil {
		panic("Open connection first")
	}

	// Make buffer from input plus 4 bytes for OFTP TCP header
	buffer := make([]byte, TCP_HEADER_LENGTH+len(input))
	length := uint32(len(input) + TCP_HEADER_LENGTH)
	binary.BigEndian.PutUint32(buffer, length)
	buffer[0] = byte(oftp2TcpMagicWord) // header as given by specification

	// Build final message
	copy(buffer[TCP_HEADER_LENGTH:], input)

	if s.Fuzzer != nil {
		buffer = s.Fuzzer(buffer)
	}

	if string(buffer[4]) == "D" && s.Verbose {
		// Hack for data message logging
		fmt.Printf("--> D...(%d bytes)\n", len(buffer[4:]))
	} else if s.Verbose {
		fmt.Printf("--> %s\n", strings.Trim(string(buffer[4:]), "\n"))
	}

	_, err := (*s.con).Write(buffer)
	if err != nil {
		return err
	}

	return nil
}

// Maximum length of a sub record
const maxSubRecordLength = 63 // Byte, from specification 6 Bit available to indicate the length

// Determine the maximum size of a file read buffer which fits into one
// data exchange buffer
func (s *OFTP2Client) maxReadBufferSize() uint32 {

	maxSubRecordCount := s.serverBufferSize / maxSubRecordLength

	// one byte per sub record, therefore we cannot fit maxSubRecordCount sub
	// records into the buffer but less. Therefore we subtract the maximum overhead
	// that may occur and calculate the record count again
	overhead := maxSubRecordCount * 1 // overhead bytes
	maxSubRecordCount = (s.serverBufferSize - overhead) / maxSubRecordLength
	resultBufferLength := maxSubRecordCount*maxSubRecordLength - 1 // minus one byte for "D" header

	return resultBufferLength
}

// Split a raw binary buffer into the sub record format needed by OFTP2. We
// assume that the provided buffer sourceBuffer completely fits into the send
// buffer of the OFTP2 protocol (see serverBufferSize in OFTP2Client struct). Due
// to the protocol overhead of the subrecords handling, the sourceBuffer cannot
// be serverBufferSize long but the overhead has to be subtracted first.
// Therefore, the buffer size bust be smaller or equal to the value returned by
// the maxReadBufferSize function.
func (s *OFTP2Client) splitBufferIntoSubRecords(sourceBuffer []byte, lastBuffer bool) []byte {
	targetBuffer := make([]byte, s.serverBufferSize)

	// now loop over the provided buffer in chunks and add the OFTP magic sub record
	// headers to the result
	targetBufferPos := 0
	sourceBufferPos := 0

	for true {
		remaining := len(sourceBuffer) - sourceBufferPos

		var bytesToCopy int
		var end bool
		var eof = 0

		if remaining == 0 {
			break
		}

		if remaining < maxSubRecordLength {
			// last chunk in the source buffer
			bytesToCopy = remaining
			end = true

			if lastBuffer {
				// this is the last record of the file
				eof = 1
			}

		} else {
			bytesToCopy = maxSubRecordLength
		}

		// add the header
		const compression = 0

		// The specification contains an error regarding the format of the subrecord
		// headers (7.2. Data Exchange Buffer Format). It states that the format is
		//
		//   0   1   2   3   4   5   6   7
		// o-------------------------------o
		// | E | C |                       |
		// | o | F | C O U N T             |
		// | R |   |                       |
		// o-------------------------------o
		//
		// but ist should be
		//
		//   0   1   2   3   4   5   6   7
		// o-------------------------------o
		// |                       | C | E |
		// | C O U N T             | F | o |
		// |                       |   | R |
		// o-------------------------------o
		//
		header := bytesToCopy<<0 +
			compression<<6 +
			eof<<7

		targetBuffer[targetBufferPos] = byte(header)
		targetBufferPos++

		for i := bytesToCopy; i > 0; i-- {
			targetBuffer[targetBufferPos] = sourceBuffer[sourceBufferPos]
			targetBufferPos++
			sourceBufferPos++
		}

		if end {
			break
		}
	}

	return targetBuffer
}
