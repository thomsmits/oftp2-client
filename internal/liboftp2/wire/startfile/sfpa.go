package startfile

import (
	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire"
)

/*
5.3.4.  SFPA - Start File Positive Answer

   o-------------------------------------------------------------------o
   |       SFPA        Start File Positive Answer                      |
   |                                                                   |
   |       Start File Phase           Speaker <---- Listener           |
   |-------------------------------------------------------------------|
   | Pos | Field     | Description                           | Format  |
   |-----+-----------+---------------------------------------+---------|
   |   0 | SFPACMD   | SFPA Command, '2'                     | F X(1)  |
   |   1 | SFPAACNT  | Answer Count                          | V 9(17) |
   o-------------------------------------------------------------------o

   SFPACMD   Command Code                                      Character

      Value: '2'  SFPA Command identifier.

   SFPAACNT  Answer Count                                    Numeric(17)

             The Listener must enter a count lower than or equal to the
             restart count specified by the Speaker in the Start File
             (SFID) command.  The count expresses the received user
             data.  If restart facilities are not available, a count of
             zero must be specified.
*/

// Start File Positive Answer
type SFPA struct {
	AnswerCount uint64
}

// SFPACMD is the command indicator for the SFPA command.
const SFPACMD = "2"

// formatDefinition returns the format definition as given in the RFC5024
func (s *SFPA) formatDefinition() []wire.FormatDefinition {
	return []wire.FormatDefinition{
		{`F X(1)`, `SFPACMD`, &[]wire.Value{{SFPACMD, "SFPA Command"}}},
		{`V 9(17)`, `SFPAACNT`, nil},
	}
}

// dataFormat returns the format definition for this command in a machine
// readable form.
func (s *SFPA) dataFormat() []wire.DataFormat {
	return wire.FormatDefinitionsToDataFormats(s.formatDefinition())
}

func (s *SFPA) Command() wire.Command {
	return wire.Command{
		Format: s.dataFormat(),
		Data:   []interface{}{SFPACMD, s.AnswerCount},
	}
}

func (s *SFPA) Marshal() []byte {
	c := s.Command()
	return c.Marshal()
}

func (s *SFPA) Parse(input []byte) error {
	buffer, err := wire.NewBuffer(&input, SFPACMD)
	if err != nil {
		return err
	}

	s.AnswerCount = buffer.GetBinQWord(17)

	return nil
}
