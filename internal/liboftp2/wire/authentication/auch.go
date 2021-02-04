package authentication

import (
	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire"
)

/*
5.3.17.  AUCH - Authentication Challenge

  o-------------------------------------------------------------------o
  |       AUCH        Authentication Challenge                        |
  |                                                                   |
  |       Start Session Phase     Initiator <---> Responder           |
  |-------------------------------------------------------------------|
  | Pos | Field     | Description                           | Format  |
  |-----+-----------+---------------------------------------+---------|
  |   0 | AUCHCMD   | AUCH Command, 'A'                     | F X(1)  |
  |   1 | AUCHCHLL  | Challenge Length                      | V U(2)  |
  |   3 | AUCHCHAL  | Challenge                             | V U(n)  |
  o-------------------------------------------------------------------o

AUCHCMD   Command Code                                      Character

      Value: 'A'  AUCH Command identifier.

   AUCHCHLL  Challenge length                                  Binary(2)

             Indicates the length of AUCHCHAL in octets.

             The length is expressed as an unsigned binary number using
             network byte order.

   AUCHCHAL  Challenge                                         Binary(n)

             A [CMS] encrypted 20-byte random number uniquely generated
             each time an AUCH is sent.

   NOTE:

   Any encryption algorithm that is available through a defined cipher
   suite (Section 10.2) may be used.  See Section 10.1 regarding the
   choice of a cipher suite.
*/

// AUCH presents an authentication challenge to the communication partner
type AUCH struct {
	Challenge []byte
}

// AUCHCMD is the command indicator for the AUCH command.
const AUCHCMD = "A"

// formatDefinition returns the format definition as given in the RFC5024
func (s *AUCH) formatDefinition() []wire.FormatDefinition {
	return []wire.FormatDefinition{
		{`F X(1)`, `AUCHCMD`, &[]wire.Value{{AUCHCMD, "AUCH Command"}}},
		{`V U(2)`, `AUCHCHLL`, nil},
		{`V U(n)`, `AUCHCHAL`, nil},
	}
}

// dataFormat returns the format definition for this command in a machine
// readable form.
func (s *AUCH) dataFormat() []wire.DataFormat {
	return wire.FormatDefinitionsToDataFormats(s.formatDefinition())
}

func (s *AUCH) Command() wire.Command {
	return wire.Command{ // TODO: Enforce 20 Byte limit?
		Format: s.dataFormat(),
		Data:   []interface{}{AUCHCMD, uint16(len(s.Challenge)), s.Challenge},
	}
}

func (s *AUCH) Marshal() []byte {
	c := s.Command()
	return c.Marshal()
}

func (s *AUCH) Parse(input []byte) error {
	buffer, err := wire.NewBuffer(&input, AUCHCMD)
	if err != nil {
		return err
	}

	lenChallenge := buffer.GetBinWord(2)
	s.Challenge = buffer.GetBytes(int(lenChallenge))

	return nil
}
