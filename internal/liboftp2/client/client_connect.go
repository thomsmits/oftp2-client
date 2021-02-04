package client

import (
	"net"
	"strconv"
	"strings"

	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire/session"
)

// Connect connects to the OFTP2 server
func (s *OFTP2Client) Connect() error {

	// open TCP connection to server
	addr := strings.Join([]string{s.ServerHost, strconv.Itoa(s.ServerPort)}, ":")
	connection, err := net.Dial("tcp", addr)

	if err != nil {
		return err
	}

	s.con = &connection

	// Read Odette MessageenvelopeIndicator
	buff, err := s.read()
	if err != nil {
		return err
	}

	ssrm := session.SSRM{}
	err = ssrm.Parse(buff)
	if err != nil {
		return err
	}
	return nil
}

// Close closes the connection to the server
func (s *OFTP2Client) Close() error {
	err := (*s.con).Close()
	if err != nil {
		return err
	}
	s.con = nil
	return nil
}
