package client

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire/authentication"
)

// AnswerChallenge answers a challenge received from the server and then sends a
// new challenge to the server for mutual authentication
func (s *OFTP2Client) AnswerChallenge(answer, ownChallenge, expectedResult []byte) error {

	if s.serverAuthenticationSupported == false {
		return errors.New("server does not support authentication")
	}

	// Send out answer to challenge
	aurp := authentication.AURP{
		Response: answer,
	}

	err := s.write(aurp.Marshal())
	if err != nil {
		return err
	}

	// Read answer from communication partner
	buffer, err := s.read()

	_, t, err := DetermineMessageType(buffer)
	if err != nil {
		return err
	}

	if t == "ESID" {
		return errors.New("wrong answer to challenge")
	} else if t != "SECD" {
		return errors.New(fmt.Sprintf("unknown answer. Expected ESID or SECD, got %s", t))
	}

	// Send out our challenge
	auch := authentication.AUCH{
		Challenge: ownChallenge,
	}

	err = s.write(auch.Marshal())
	if err != nil {
		return err
	}

	// Read answer from communication partner
	buffer, err = s.read()
	if err != nil {
		return err
	}

	serverAurp := authentication.AURP{}
	err = serverAurp.Parse(buffer)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(serverAurp.Response, expectedResult) {
		return errors.New(fmt.Sprintf("answer from server does not match. %v != %v", serverAurp.Response, expectedResult))
	}

	return nil
}
