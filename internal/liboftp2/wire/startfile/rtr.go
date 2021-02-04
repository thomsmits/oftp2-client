package startfile

import (
	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire"
)

/*
5.3.15.  RTR - Ready To Receive

   o-------------------------------------------------------------------o
   |       RTR         Ready To Receive                                |
   |                                                                   |
   |       Start File Phase         Initiator <---- Responder          |
   |       End File Phase           Initiator <---- Responder          |
   |-------------------------------------------------------------------|
   | Pos | Field     | Description                           | Format  |
   |-----+-----------+---------------------------------------+---------|
   |   0 | RTRCMD    | RTR Command, 'P'                      | F X(1)  |
   o-------------------------------------------------------------------o

   RTRCMD    Command Code                                      Character

      Value: 'P'  RTR Command identifier.
*/

// RTR signals that the communication partner is ready to Ready To Receive
type RTR struct {
}

// RTRCMD is the command indicator for the RTR command.
const RTRCMD = "P"

// formatDefinition returns the format definition as given in the RFC5024
func (s *RTR) formatDefinition() []wire.FormatDefinition {
	return []wire.FormatDefinition{
		{`F X(1)`, `RTRCMD`, &[]wire.Value{{RTRCMD, "RTR Command"}}},
	}
}

// dataFormat returns the format definition for this command in a machine
// readable form.
func (s *RTR) dataFormat() []wire.DataFormat {
	return wire.FormatDefinitionsToDataFormats(s.formatDefinition())
}

func (s *RTR) Command() wire.Command {
	return wire.Command{
		Format: s.dataFormat(),
		Data:   []interface{}{RTRCMD},
	}
}

func (s *RTR) Marshal() []byte {
	c := s.Command()
	return c.Marshal()
}

func (s *RTR) Parse(input []byte) error {
	_, err := wire.NewBuffer(&input, RTRCMD)
	if err != nil {
		return err
	}

	return nil
}
