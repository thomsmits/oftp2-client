package client

import (
	"errors"
	"fmt"

	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire/session"
)

func (s *OFTP2Client) internalQueryServerCapabilities(auth bool) (session.SSID, error) {

	err := s.Connect()
	if err != nil {
		return session.SSID{}, err
	}

	ssid := session.SSID{
		Id:             s.OdetteId,
		Password:       "",
		BufferSize:     102400,
		Capability:     "S",
		Compress:       true,
		Restart:        true,
		Special:        true,
		Credit:         999,
		Authentication: auth,
		UserData:       "",
	}

	// Send out SSID to the server
	err = s.write(ssid.Marshal())
	if err != nil {
		return session.SSID{}, err
	}

	// Read the server's SSID
	buffer, err := s.read()
	if err != nil {
		return session.SSID{}, err
	}

	answer, t, err := DetermineMessageType(buffer)
	if err != nil {
		return session.SSID{}, err
	}

	if t == "ESID" && !auth {
		return session.SSID{}, errors.New(fmt.Sprintf("server terminated helo: %v", answer))
	} else if t == "ESID" && auth {
		// server does not support authentication
		ssid.Authentication = false
		return ssid, nil
	} else if t != "SSID" {
		return session.SSID{}, errors.New(fmt.Sprintf("server send unexpected answer: %v", answer))
	}

	serverSSID := session.SSID{}
	err = serverSSID.Parse(buffer)
	if err != nil {
		return session.SSID{}, err
	}

	// close session
	err = s.EndSession()
	if err != nil {
		// ignore error here
	}

	err = s.Close()
	if err != nil {
		// ignore error here
	}

	// next try with authentication set to true
	return serverSSID, nil
}

// QueryServerCapabilities connects to the server and tries to find out what
// capabilities the server supports. To do this, a session is initiated and the
// server is presented with a client that seems to support all OFTP2 features.
// The server then has to answer with its features. After that the session is closed.
//
// This methods opens and closes the connection to the server, therefore it is
// not necessary to call Connect before using this method.
func (s *OFTP2Client) QueryServerCapabilities() (session.SSID, error) {

	// we cannot test secure authentication in the first shot because the server will
	// answer with an ESID(r=12) in case of an security mismatch. Therefore, we start
	// with no authentication and try to check it in the next step
	ssid, err := s.internalQueryServerCapabilities(false)
	if err != nil {
		return session.SSID{}, err
	}

	// next try with authentication set to true
	secSsid, err := s.internalQueryServerCapabilities(true)
	if err != nil {
		return session.SSID{}, err
	}

	// copy the authentication setting
	ssid.Authentication = secSsid.Authentication

	return ssid, nil
}
