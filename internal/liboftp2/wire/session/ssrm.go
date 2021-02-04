package session

import (
	"errors"
	"fmt"

	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire"
)

/*
   o-------------------------------------------------------------------o
   |       SSRM        Start Session Ready Message                     |
   |                                                                   |
   |       Start Session Phase     Initiator <---- Responder           |
   |-------------------------------------------------------------------|
   | Pos | Field     | Description                           | Format  |
   |-----+-----------+---------------------------------------+---------|
   |   0 | SSRMCMD   | SSRM Command, 'I'                     | F X(1)  |
   |   1 | SSRMMSG   | Ready Message, 'ODETTE FTP READY '    | F X(17) |
   |  18 | SSRMCR    | Carriage Return                       | F X(1)  |
   o-------------------------------------------------------------------o
*/

//  Start Session Ready Message
type SSRM struct {
}

// SSRMCMD is the command indicator for the SSRM command.
const SSRMCMD = "I"

// formatDefinition returns the format definition as given in the RFC5024
func (s *SSRM) formatDefinition() []wire.FormatDefinition {
	return []wire.FormatDefinition{
		{`F X(1)`, `SSRMCMD`, &[]wire.Value{{SSRMCMD, "SSRM Command"}}},
		{`F X(17)`, `SSRMMSG`, &[]wire.Value{{"ODETTE FTP READY ", "HELO"}}},
		{`F X(1)`, ` SSRMCR`, wire.ValueNewline},
	}
}

// dataFormat returns the format definition for this command in a machine
// readable form.
func (s *SSRM) dataFormat() []wire.DataFormat {
	return wire.FormatDefinitionsToDataFormats(s.formatDefinition())
}

func (s *SSRM) Command() wire.Command {
	return wire.Command{
		Format: s.dataFormat(),
		Data:   []interface{}{SSRMCMD, "ODETTE FTP READY ", "\n"},
	}
}

func (s *SSRM) Marshal() []byte {
	c := s.Command()
	return c.Marshal()
}

func (s *SSRM) Parse(input []byte) error {
	buffer, err := wire.NewBuffer(&input, SSRMCMD)
	if err != nil {
		return err
	}

	helo := buffer.GetString(17)
	_ = buffer.GetString(1)

	if helo != "ODETTE FTP READY" {
		return errors.New(fmt.Sprintf("helo %s", helo))
	}

	return nil
}
