package wire

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Buffer represents data received over the network. It offers method to parse
// data from the buffer into other formats. Especially the exotic OFTP2 encodings
// are supported.
type Buffer struct {
	data *[]byte
	pos  int
}

// NewBuffer creates a new buffer for the given data and checks if the buffer
// starts with the provided marker. If not, no buffer is created but an error is
// returned.
func NewBuffer(data *[]byte, marker string) (*Buffer, error) {
	b := &Buffer{data: data, pos: 0}
	err := b.checkMarker(marker)
	if err != nil {
		return nil, err
	} else {
		return b, nil
	}
}

// checkMarker tests if the byte the buffer currently points at is equal to
// marker. If not, an error is returned, otherwise nil.
func (b *Buffer) checkMarker(marker string) error {
	m := b.GetString(1)
	if marker != marker {
		return errors.New(fmt.Sprintf("wrong marker; found %s, expected %s", m, marker))
	} else {
		return nil
	}
}

// GetString gets a string from the buffer at the current position with the length len.
// The position is afterwards incremented, so that it points to the next data
// portion in the input array.
func (b *Buffer) GetString(len int) string {
	result := b.GetBytes(len)
	return strings.Trim(string(result), " ")
}

// GetNumInt gets an int from the buffer at the current position with the length len. The
// position is afterwards incremented, so that it points to the next data portion
// in the input array. The integer value is expected to be encoded in ASCII, e.g.
// "142" for 142.
func (b *Buffer) GetNumInt(len int) int {
	result := b.GetString(len)
	intValue, _ := strconv.Atoi(result)
	return intValue
}

// GetBinWord gets an uint16 from the buffer at the current position with the length len.
// The position is afterwards incremented, so that it points to the next data
// portion in the input array. The integer value is expected to be encoded in
// network byte order, e.g. 0xba 0xbe for 0xbabe.
func (b *Buffer) GetBinWord(len int) uint16 {
	bytes := b.GetBytes(len)
	return binary.BigEndian.Uint16(bytes)
}

// GetBinDWord gets an uint32 from the buffer at the current position with the length len.
// The position is afterwards incremented, so that it points to the next data
// portion in the input array. The integer value is expected to be encoded in
// network byte order, e.g. 0xca 0xfe 0xba 0xbe for 0xcafebabe.
func (b *Buffer) GetBinDWord(len int) uint32 {
	result := b.GetString(len)
	intValue, _ := strconv.ParseUint(result, 10, 32)
	return uint32(intValue)
}

// GetBinQWord gets an uint64 from the buffer at the current position with the length len.
// The position is afterwards incremented, so that it points to the next data
// portion in the input array. The integer value is expected to be encoded in
// network byte order, e.g. 0xca 0xfe 0xba 0xbe 0xca 0xfe 0xba 0xbe for
// 0xcafebabecafebabe.
func (b *Buffer) GetBinQWord(len int) uint64 {
	result := b.GetString(len)
	intValue, _ := strconv.ParseUint(result, 10, 64)
	return intValue
}

// GetBool gets an boolean from the buffer at the current position with the length len.
// The position is afterwards incremented, so that it points to the next data
// portion in the input array. The boolean value is expected to be encoded as "Y"
// and "N" for true and false.
func (b *Buffer) GetBool(len int) bool {
	result := b.GetString(len)
	return result == "Y"
}

// GetBytes gets the raw bytes from the input byte array at position pos of length
// len. The position is afterwards incremented, so that it points to the next
// data portion in the input array.
func (b *Buffer) GetBytes(len int) []byte {
	result := (*b.data)[b.pos : b.pos+len]
	b.pos += len
	return result
}
