package transfer

import (
	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire"
)

/*
5.3.6.  DATA - Data Exchange Buffer

   o-------------------------------------------------------------------o
   |       DATA        Data Exchange Buffer                            |
   |                                                                   |
   |       Data Transfer Phase        Speaker ----> Listener           |
   |-------------------------------------------------------------------|
   | Pos | Field     | Description                           | Format  |
   |-----+-----------+---------------------------------------+---------|
   |   0 | DATACMD   | DATA Command, 'D'                     | F X(1)  |
   |   1 | DATABUF   | Data Exchange Buffer payload          | V U(n)  |
   o-------------------------------------------------------------------o

   DATACMD   Command Code                                      Character

      Value: 'D'  DATA Command identifier.

   DATABUF   Data Exchange Buffer payload                      Binary(n)

             Variable-length buffer containing the data payload.  The
             Data Exchange Buffer is described in Section 7.
*/

// DATA is used to exchange data (Data Exchange Buffer in OFTP lingo)
type DATA struct {
	Length uint64
	Buffer []byte
}

// DATACMD is the command indicator for the DATA command.
const DATACMD = "D"

// formatDefinition returns the format definition as given in the RFC5024
func (s *DATA) formatDefinition() []wire.FormatDefinition {
	return []wire.FormatDefinition{
		{`F X(1)`, `DATACMD`, &[]wire.Value{{DATACMD, "DATA Command"}}},
		{`V U(n)`, `DATABUF`, nil},
	}
}

// dataFormat returns the format definition for this command in a machine
// readable form.
func (s *DATA) dataFormat() []wire.DataFormat {
	return wire.FormatDefinitionsToDataFormats(s.formatDefinition())
}

func (s *DATA) Command() wire.Command {
	return wire.Command{
		Format: s.dataFormat(),
		Data:   []interface{}{DATACMD, s.Buffer},
	}
}

func (s *DATA) Marshal() []byte {
	c := s.Command()
	return c.Marshal()
}

// TODO: Handling of buffer length not correct
func (s *DATA) Parse(input []byte) error {
	buffer, err := wire.NewBuffer(&input, DATACMD)
	if err != nil {
		return err
	}

	s.Buffer = buffer.GetBytes(int(s.Length))

	return nil
}
