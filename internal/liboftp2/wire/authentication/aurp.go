package authentication

import (
	"fmt"
	"reflect"

	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire"
)

/*
5.3.18.  AURP - Authentication Response

   o-------------------------------------------------------------------o
   |       AURP        Authentication Response                         |
   |                                                                   |
   |       Start Session Phase     Initiator <---> Responder           |
   |-------------------------------------------------------------------|
   | Pos | Field     | Description                           | Format  |
   |-----+-----------+---------------------------------------+---------|
   |   0 | AURPCMD   | AURP Command, 'S'                     | F X(1)  |
   |   1 | AURPRSP   | Response                              | V U(20) |
   o-------------------------------------------------------------------o

   AURPCMD   Command Code                                      Character

      Value: 'S'  AURP Command identifier.

   AURPRSP   Response                                         Binary(20)

             Contains the decrypted challenge (AUCHCHAL).

   IMPORTANT:

   It is an application implementation issue to validate a received AURP
   to ensure that the response matches the challenge.  This validation
   is extremely important to ensure that a party is correctly
   authenticated.
*/

// AURP contains the response to an authentication challenge (AUCH)
type AURP struct {
	Response []byte
}

// AURPCMD is the command indicator for the AURP command.
const AURPCMD = "S"

// formatDefinition returns the format definition as given in the RFC5024
func (s *AURP) formatDefinition() []wire.FormatDefinition {
	return []wire.FormatDefinition{
		{`F X(1)`, `AURPCMD`, &[]wire.Value{{AURPCMD, "AURP Command"}}},
		{`V U(20)`, `AURPRSP`, nil},
	}
}

// dataFormat returns the format definition for this command in a machine
// readable form.
func (s *AURP) dataFormat() []wire.DataFormat {
	return wire.FormatDefinitionsToDataFormats(s.formatDefinition())
}

func (s *AURP) Command() wire.Command {
	return wire.Command{
		Format: s.dataFormat(),
		Data:   []interface{}{AURPCMD, s.Response},
	}
}

func (s *AURP) Marshal() []byte {
	c := s.Command()
	return c.Marshal()
}

func (s *AURP) Parse(input []byte) error {
	buffer, err := wire.NewBuffer(&input, AURPCMD)
	if err != nil {
		return err
	}

	s.Response = buffer.GetBytes(20)

	return nil
}

func ToAURP(data *interface{}) *AURP {

	var aurp AURP

	switch (*data).(type) {
	case AURP:
		aurp = (*data).(AURP)
		return &aurp
	default:
		panic(fmt.Sprintf("wrong type of data. Expected %s, got %s", "AURP", reflect.TypeOf(data)))
		return nil // never reached
	}
}
