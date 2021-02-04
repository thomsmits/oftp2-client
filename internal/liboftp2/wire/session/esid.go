package session

import (
	"fmt"

	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire"
)

/*
5.3.11.  ESID - End Session

   o-------------------------------------------------------------------o
   |       ESID        End Session                                     |
   |                                                                   |
   |       End Session Phase          Speaker ----> Listener           |
   |-------------------------------------------------------------------|
   | Pos | Field     | Description                           | Format  |
   |-----+-----------+---------------------------------------+---------|
   |   0 | ESIDCMD   | ESID Command, 'F'                     | F X(1)  |
   |   1 | ESIDREAS  | Reason Code                           | F 9(2)  |
   |   3 | ESIDREASL | Reason Text Length                    | V 9(3)  |
   |   6 | ESIDREAST | Reason Text                           | V T(n)  |
   |     | ESIDCR    | Carriage Return                       | F X(1)  |
   o-------------------------------------------------------------------o

   ESIDCMD   Command Code                                      Character

      Value: 'F'  ESID Command identifier.

   ESIDREAS  Reason Code                                      Numeric(2)

      Value: '00'  Normal session termination

             '01'  Command not recognised

                   An Exchange Buffer contains an invalid command code
                   (1st octet of the buffer).

             '02'  Protocol violation

                   An Exchange Buffer contains an invalid command for
                   the current state of the receiver.

             '03'  User code not known

                   A Start Session (SSID) command contains an unknown or
                   invalid Identification Code.

             '04'  Invalid password

                   A Start Session (SSID) command contained an invalid
                   password.

             '05'  Local site emergency close down

                   The local site has entered an emergency close down
                   mode.  Communications are being forcibly terminated.

             '06'  Command contained invalid data

                   A field within a Command Exchange Buffer contains
                   invalid data.

             '07'  Exchange Buffer size error

                   The length of the Exchange Buffer as determined by
                   the Stream Transmission Header differs from the
                   length implied by the Command Code.

             '08'  Resources not available

                   The request for connection has been denied due to a
                   resource shortage.  The connection attempt should be
                   retried later.

             '09'  Time out

             '10'  Mode or capabilities incompatible

             '11'  Invalid challenge response

             '12'  Secure authentication requirements incompatible

             '99'  Unspecified Abort code

                   An error was detected for which no specific code is
                   defined.

   ESIDREASL Reason Text Length                               Numeric(3)

             Length in octets of the field ESIDREAST.

             0 indicates that no ESIDREAST field is present.

   ESIDREAST Reason Text                                      [UTF-8](n)

             Reason why session ended in plain text.

             It is encoded using [UTF-8].

             Maximum length of the encoded reason is 999 octets.

             No general structure is defined for this attribute.

   ESIDCR    Carriage Return                                   Character

      Value: Character with hex value '0D' or '8D'.
*/

// ESID terminates the communication session
type ESID struct {
	ReasonCode int
	ReasonText string
}

// ESIDCMD is the command indicator for the ESID command.
const ESIDCMD = "F"

func NewESID(reasonCode int, reasonText string) *ESID {
	return &ESID{ReasonCode: reasonCode, ReasonText: reasonText}
}

func (s *ESID) Command() wire.Command {

	s.ReasonText = wire.TruncateString(s.ReasonText, 999) // TODO: Unicode length

	result := wire.Command{
		Format: s.dataFormat(),
		Data:   []interface{}{ESIDCMD, s.ReasonCode, len([]byte(s.ReasonText)), s.ReasonText, "\n"},
	}

	return result
}

var valuesESIDREAS = map[int]string{
	0:  "Normal session termination",
	1:  "Command not recognised",
	2:  "Protocol violation",
	3:  "User code not known",
	4:  "Invalid password",
	5:  "Local site emergency close down",
	6:  "Command contained invalid data",
	7:  "Exchange Buffer size error",
	8:  "Resources not available",
	9:  "Time out",
	10: "Mode or capabilities incompatible",
	11: "Invalid challenge response",
	12: "Secure authentication requirements incompatible",
	99: "Unspecified Abort code",
}

// formatDefinition returns the format definition as given in the RFC5024
func (s *ESID) formatDefinition() []wire.FormatDefinition {
	return []wire.FormatDefinition{
		{`F X(1)`, `ESIDCMD`, &[]wire.Value{{ESIDCMD, "ESID Command"}}},
		{`F 9(2)`, `ESIDREAS`, wire.IntMapToValues(valuesESIDREAS, 2)},
		{`V 9(3)`, `ESIDREASL`, nil},
		{`V T(n)`, `ESIDREAST`, nil},
		{`F X(1)`, `ESIDCR`, wire.ValueNewline},
	}
}

// dataFormat returns the format definition for this command in a machine
// readable form.
func (s *ESID) dataFormat() []wire.DataFormat {
	return wire.FormatDefinitionsToDataFormats(s.formatDefinition())
}

func (s *ESID) Marshal() []byte {
	c := s.Command()
	return c.Marshal()
}

func (s *ESID) Parse(input []byte) error {
	buffer, err := wire.NewBuffer(&input, ESIDCMD)
	if err != nil {
		return err
	}

	s.ReasonCode = buffer.GetNumInt(2)
	lenReasonText := buffer.GetNumInt(3)
	s.ReasonText = buffer.GetString(lenReasonText)
	_ = buffer.GetString(1)

	return nil
}

func (s *ESID) String() string {
	return fmt.Sprintf("ESID - End Session. Reason Code: %d (%s), Text: %s", s.ReasonCode, valuesESIDREAS[s.ReasonCode], s.ReasonText)
}

func BuildESID(reasonCode int, reasonText string) ESID {
	esid := ESID{
		ReasonCode: reasonCode,
		ReasonText: reasonText,
	}
	return esid
}
