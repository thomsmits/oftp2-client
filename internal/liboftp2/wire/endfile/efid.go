package endfile

import (
	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire"
)

/*
5.3.8.  EFID - End File

   o-------------------------------------------------------------------o
   |       EFID        End File                                        |
   |                                                                   |
   |       End File Phase             Speaker ----> Listener           |
   |-------------------------------------------------------------------|
   | Pos | Field     | Description                           | Format  |
   |-----+-----------+---------------------------------------+---------|
   |   0 | EFIDCMD   | EFID Command, 'T'                     | F X(1)  |
   |   1 | EFIDRCNT  | Record Count                          | V 9(17) |
   |  18 | EFIDUCNT  | Unit Count                            | V 9(17) |
   o-------------------------------------------------------------------o

   EFIDCMD   Command Code                                      Character

      Value: 'T'  EFID Command identifier.

   EFIDRCNT  Record Count                                    Numeric(17)

    Maximum: 99999999999999999

             For SSIDFMT 'F' or 'V', the exact record count.
             For SSIDFMT 'U' or 'T', zeros.

             The count will express the real size of the file (before
             buffer compression, header not included).  The total count
             is always used, even during restart processing.

   EFIDUCNT  Unit Count                                      Numeric(17)

    Maximum: 99999999999999999

             Exact number of units (octets) transmitted.

             The count will express the real size of the file.  The
             total count is always used, even during restart processing.
*/

// EFID ends the transport of a file
type EFID struct {
	RecordCount uint64
	UnitCount   uint64
}

// EFIDCMD is the command indicator for the EFID command.
const EFIDCMD = "T"

// formatDefinition returns the format definition as given in the RFC5024
func (s *EFID) formatDefinition() []wire.FormatDefinition {
	return []wire.FormatDefinition{
		{`F X(1)`, `EFIDCMD`, &[]wire.Value{{EFIDCMD, "EFID Command"}}},
		{`V 9(17)`, `EFIDRCNT`, nil},
		{`V 9(17)`, `EFIDUCNT`, nil},
	}
}

// dataFormat returns the format definition for this command in a machine
// readable form.
func (s *EFID) dataFormat() []wire.DataFormat {
	return wire.FormatDefinitionsToDataFormats(s.formatDefinition())
}

func (s *EFID) Command() wire.Command {
	return wire.Command{
		Format: s.dataFormat(),
		Data:   []interface{}{EFIDCMD, s.RecordCount, s.UnitCount},
	}
}

func (s *EFID) Marshal() []byte {
	c := s.Command()
	return c.Marshal()
}

func (s *EFID) Parse(input []byte) error {
	buffer, err := wire.NewBuffer(&input, EFIDCMD)
	if err != nil {
		return err
	}

	s.RecordCount = buffer.GetBinQWord(17)
	s.UnitCount = buffer.GetBinQWord(17)

	return nil
}
