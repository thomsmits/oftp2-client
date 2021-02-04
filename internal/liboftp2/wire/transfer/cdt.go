package transfer

import (
	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire"
)

/*
5.3.7.  CDT - Set Credit

   o-------------------------------------------------------------------o
   |       CDT         Set Credit                                      |
   |                                                                   |
   |       Data Transfer Phase        Speaker <---- Listener           |
   |-------------------------------------------------------------------|
   | Pos | Field     | Description                           | Format  |
   |-----+-----------+---------------------------------------+---------|
   |   0 | CDTCMD    | CDT Command, 'C'                      | F X(1)  |
   |   1 | CDTRSV1   | Reserved                              | F X(2)  |
   o-------------------------------------------------------------------o

   CDTCMD    Command Code                                      Character
      Value: 'C'  CDT Command identifier.

   CDTRSV1   Reserved                                          String(2)
             This field is reserved for future use.
*/

// CDT sets credits (flow control)
type CDT struct {
}

// CDTCMD is the command indicator for the CDT command.
const CDTCMD = "C"

// formatDefinition returns the format definition as given in the RFC5024
func (s *CDT) formatDefinition() []wire.FormatDefinition {
	return []wire.FormatDefinition{
		{`F X(1)`, `CDTCMD`, &[]wire.Value{{CDTCMD, "CDT Command"}}},
		{`F X(2)`, `CDTRSV1`, &[]wire.Value{{"  ", "Reserved"}}},
	}
}

// dataFormat returns the format definition for this command in a machine
// readable form.
func (s *CDT) dataFormat() []wire.DataFormat {
	return wire.FormatDefinitionsToDataFormats(s.formatDefinition())
}

func (s *CDT) Command() wire.Command {
	return wire.Command{
		Format: s.dataFormat(),
		Data:   []interface{}{CDTCMD, "  "},
	}
}

func (s *CDT) Marshal() []byte {
	c := s.Command()
	return c.Marshal()
}

func (s *CDT) Parse(input []byte) error {
	_, err := wire.NewBuffer(&input, CDTCMD)
	if err != nil {
		return err
	}

	return nil
}
