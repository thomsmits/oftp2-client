package startfile

import (
	"fmt"

	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire"
)

/*
5.3.5.  SFNA - Start File Negative Answer

   o-------------------------------------------------------------------o
   |       SFNA        Start File Negative Answer                      |
   |                                                                   |
   |       Start File Phase           Speaker <---- Listener           |
   |-------------------------------------------------------------------|
   | Pos | Field     | Description                           | Format  |
   |-----+-----------+---------------------------------------+---------|
   |   0 | SFNACMD   | SFNA Command, '3'                     | F X(1)  |
   |   1 | SFNAREAS  | Answer Reason                         | F 9(2)  |
   |   3 | SFNARRTR  | Retry Indicator, (Y/N)                | F X(1)  |
   |   4 | SFNAREASL | Answer Reason Text Length             | V 9(3)  |
   |   7 | SFNAREAST | Answer Reason Text                    | V T(n)  |
   o-------------------------------------------------------------------o

   SFNACMD   Command Code                                      Character

      Value: '3'  SFNA Command identifier.

   SFNAREAS  Answer Reason                                    Numeric(2)

      Value: '01'  Invalid filename.
             '02'  Invalid destination.
             '03'  Invalid origin.
             '04'  Storage record format not supported.
             '05'  Maximum record length not supported.
             '06'  File size is too big.
             '10'  Invalid record count.
             '11'  Invalid byte count.
             '12'  Access method failure.
             '13'  Duplicate file.
             '14'  File direction refused.
             '15'  Cipher suite not supported.
             '16'  Encrypted file not allowed.
             '17'  Unencrypted file not allowed.
             '18'  Compression not allowed.
             '19'  Signed file not allowed.
             '20'  Unsigned file not allowed.
             '99'  Unspecified reason.

             Reason why transmission cannot proceed.

   SFNARRTR  Retry Indicator                                   Character

      Value: 'N'  Transmission should not be retried.
             'Y'  The transmission may be retried later.

             This parameter is used to advise the Speaker if it should
             retry at a later time due to a temporary condition at the
             Listener site, such as a lack of storage space.  It should
             be used in conjunction with the Answer Reason code
             (SFNAREAS).

             An invalid file name error code may be the consequence of a
             problem in the mapping of the Virtual File on to a real
             file.  Such problems cannot always be resolved immediately.
             It is therefore recommended that when an SFNA with Retry =
             Y is received the User Monitor attempts to retransmit the
             relevant file in a subsequent session.

   SFNAREASL Answer Reason Text Length                        Numeric(3)

             Length in octets of the field SFNAREAST.

             0 indicates that no SFNAREAST field follows.

   SFNAREAST Answer Reason Text                               [UTF-8](n)

             Reason why transmission cannot proceed in plain text.

             It is encoded using [UTF-8].

             Maximum length of the encoded reason is 999 octets.

             No general structure is defined for this attribute.
*/

// SFNA Start File Negative Answer
type SFNA struct {
	ReasonCode     int
	RetryIndicator bool
	ReasonText     string
}

// SFNACMD is the command indicator for the SFNA command.
const SFNACMD = "3"

var valuesSFNAREAS = map[int]string{
	1:  "Invalid filename.",
	2:  "Invalid destination.",
	3:  "Invalid origin.",
	4:  "Storage record format not supported.",
	5:  "Maximum record length not supported.",
	6:  "File size is too big.",
	10: "Invalid record count.",
	11: "Invalid byte count.",
	12: "Access method failure.",
	13: "Duplicate file.",
	14: "File direction refused.",
	15: "Cipher suite not supported.",
	16: "Encrypted file not allowed.",
	17: "Unencrypted file not allowed.",
	18: "Compression not allowed.",
	19: "Signed file not allowed.",
	20: "Unsigned file not allowed.",
	99: "Unspecified reason.",
}

// formatDefinition returns the format definition as given in the RFC5024
func (s *SFNA) formatDefinition() []wire.FormatDefinition {
	return []wire.FormatDefinition{
		{`F X(1)`, `SFNACMD`, &[]wire.Value{{SFNACMD, "SFNA Command"}}},
		{`F 9(2)`, `SFNAREAS`, wire.IntMapToValues(valuesSFNAREAS, 2)},
		{`F X(1)`, `SFNARRTR`, wire.ValueBooleanYesNo},
		{`V 9(3)`, `SFNAREASL`, nil},
		{`V T(n)`, `SFNAREAST`, nil},
	}
}

// dataFormat returns the format definition for this command in a machine
// readable form.
func (s *SFNA) dataFormat() []wire.DataFormat {
	return wire.FormatDefinitionsToDataFormats(s.formatDefinition())
}

func (s *SFNA) Command() wire.Command {

	s.ReasonText = wire.TruncateString(s.ReasonText, 999) // TODO: unicode length

	result := wire.Command{
		Format: s.dataFormat(),
		Data: []interface{}{
			SFNACMD,
			s.ReasonCode,
			s.RetryIndicator,
			len([]byte(s.ReasonText)),
			s.ReasonText,
		},
	}

	return result
}

func (s *SFNA) Marshal() []byte {
	c := s.Command()
	return c.Marshal()
}

func (s *SFNA) Parse(input []byte) error {
	buffer, err := wire.NewBuffer(&input, SFNACMD)
	if err != nil {
		return err
	}

	s.ReasonCode = buffer.GetNumInt(2)
	s.RetryIndicator = buffer.GetBool(1)
	lenAnswerText := buffer.GetNumInt(3)
	s.ReasonText = buffer.GetString(lenAnswerText)

	return nil
}

func (s *SFNA) String() string {
	return fmt.Sprintf("SFNA - Start File Negative Answer. Reason Code: %d (%s), Text: %s", s.ReasonCode, valuesSFNAREAS[s.ReasonCode], s.ReasonText)
}
