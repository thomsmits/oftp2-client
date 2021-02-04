package endfile

import (
	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire"
)

/*
5.3.9.  EFPA - End File Positive Answer

   o-------------------------------------------------------------------o
   |       EFPA        End File Positive Answer                        |
   |                                                                   |
   |       End File Phase             Speaker <---- Listener           |
   |-------------------------------------------------------------------|
   | Pos | Field     | Description                           | Format  |
   |-----+-----------+---------------------------------------+---------|
   |   0 | EFPACMD   | EFPA Command, '4'                     | F X(1)  |
   |   1 | EFPACD    | Change Direction Indicator, (Y/N)     | F X(1)  |
   o-------------------------------------------------------------------o

   EFPACMD   Command Code                                      Character

      Value: '4'  EFPA Command identifier.

   EFPACD    Change Direction Indicator                        Character

      Value: 'N'  Change direction not requested.
             'Y'  Change direction requested.

             This parameter allows the Listener to request a Change
             Direction (CD) command from the Speaker.
*/

// EFPA signals that the end of file command (EFID) could be handled
// successfully
type EFPA struct {
	ChangeDirection bool
}

// EFPACMD is the command indicator for the EFPA command.
const EFPACMD = "4"

func (s *EFPA) Command() wire.Command {
	return wire.Command{
		Format: s.dataFormat(),
		Data:   []interface{}{EFPACMD, s.ChangeDirection},
	}
}

// formatDefinition returns the format definition as given in the RFC5024
func (s *EFPA) formatDefinition() []wire.FormatDefinition {
	return []wire.FormatDefinition{
		{`F X(1)`, `EFPACMD`, &[]wire.Value{{EFPACMD, "EFPA Command"}}},
		{`F X(1)`, `EFPACD`, wire.ValueBooleanYesNo},
	}
}

// dataFormat returns the format definition for this command in a machine
// readable form.
func (s *EFPA) dataFormat() []wire.DataFormat {
	return wire.FormatDefinitionsToDataFormats(s.formatDefinition())
}

func (s *EFPA) Marshal() []byte {
	c := s.Command()
	return c.Marshal()
}

func (s *EFPA) Parse(input []byte) error {
	buffer, err := wire.NewBuffer(&input, EFPACMD)
	if err != nil {
		return err
	}

	s.ChangeDirection = buffer.GetBool(1)

	return nil
}
