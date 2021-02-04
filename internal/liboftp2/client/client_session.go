package client

import (
	"errors"
	"fmt"

	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire/session"
)

// StartSession opens a session with the server
func (s *OFTP2Client) StartSession(password string, compression, restart, authentication bool) error {

	ssid := session.SSID{
		Id:             s.OdetteId,
		Password:       password,
		BufferSize:     1024,
		Capability:     "S",
		Compress:       compression,
		Restart:        restart,
		Special:        authentication,
		Credit:         999,
		Authentication: false,
		UserData:       "",
	}

	// Send out SSID to the server
	err := s.write(ssid.Marshal())
	if err != nil {
		return err
	}

	// Read the server's SSID
	buffer, err := s.read()
	if err != nil {
		return err
	}

	answer, t, err := DetermineMessageType(buffer)
	if err != nil {
		return err
	}

	if t == "ESID" {
		return errors.New(fmt.Sprintf("server terminated helo: %v", answer))
	} else if t != "SSID" {
		return errors.New(fmt.Sprintf("server send unexpected answer: %v", answer))
	}

	serverSSID := session.SSID{}
	err = serverSSID.Parse(buffer)
	if err != nil {
		return err
	}

	// negotiation of security is not allowed
	if serverSSID.Authentication != authentication {
		_ = s.EndSession() // ignore error, we are anyhow lost
		return errors.New("cannot agree on security features")
	}

	s.serverId = serverSSID.Id
	s.serverPassword = serverSSID.Password
	s.serverBufferSize = serverSSID.BufferSize
	s.serverCapability = serverSSID.Capability
	s.serverCompress = serverSSID.Compress
	s.serverRestartSupported = serverSSID.Restart
	s.serverSpecial = serverSSID.Special
	s.serverCredit = serverSSID.Credit
	s.serverAuthenticationSupported = serverSSID.Authentication
	s.serverUserData = serverSSID.UserData

	return nil
}

// EndSession closes the session
func (s *OFTP2Client) EndSession() error {

	esid := session.ESID{
		ReasonCode: 0,
		ReasonText: "OK",
	}

	// Send out ESID to the server
	err := s.write(esid.Marshal())
	if err != nil {
		return err
	}

	return nil
}
