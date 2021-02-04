package authentication

import (
	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire"
)

/*
5.3.16.  SECD - Security Change Direction

   o-------------------------------------------------------------------o
   |       SECD        Security Change Direction                       |
   |                                                                   |
   |       Start Session Phase     Initiator <---> Responder           |
   |-------------------------------------------------------------------|
   | Pos | Field     | Description                           | Format  |
   |-----+-----------+---------------------------------------+---------|
   |   0 | SECDCMD   | SECD Command, 'J'                     | F X(1)  |
   o-------------------------------------------------------------------o

   SECDCMD   Command Code                                      Character

      Value: 'J'  SECD Command identifier.
*/

// SECD Security Change Direction
type SECD struct {
}

// SECDCMD is the command indicator for the SECD command.
const SECDCMD = "J"

// formatDefinition returns the format definition as given in the RFC5024
func (s *SECD) formatDefinition() []wire.FormatDefinition {
	return []wire.FormatDefinition{
		{`F X(1)`, `SECDCMD`, &[]wire.Value{{SECDCMD, "SECD Command"}}},
	}
}

// dataFormat returns the format definition for this command in a machine
// readable form.
func (s *SECD) dataFormat() []wire.DataFormat {
	return wire.FormatDefinitionsToDataFormats(s.formatDefinition())
}

func (s *SECD) Command() wire.Command {
	return wire.Command{
		Format: s.dataFormat(),
		Data:   []interface{}{SECDCMD},
	}
}

func (s *SECD) Marshal() []byte {
	c := s.Command()
	return c.Marshal()
}

func (s *SECD) Parse(input []byte) error {
	_, err := wire.NewBuffer(&input, SECDCMD)
	if err != nil {
		return err
	}

	return nil
}
