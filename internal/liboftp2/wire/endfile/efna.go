package endfile

import (
	"fmt"

	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire"
)

/*
5.3.10.  EFNA - End File Negative Answer

   o-------------------------------------------------------------------o
   |       EFNA        End File Negative Answer                        |
   |                                                                   |
   |       End File Phase             Speaker <---- Listener           |
   |-------------------------------------------------------------------|
   | Pos | Field     | Description                           | Format  |
   |-----+-----------+---------------------------------------+---------|
   |   0 | EFNACMD   | EFNA Command, '5'                     | F X(1)  |
   |   1 | EFNAREAS  | Answer Reason                         | F 9(2)  |
   |   3 | EFNAREASL | Answer Reason Text Length             | V 9(3)  |
   |   6 | EFNAREAST | Answer Reason Text                    | V T(n)  |
   o-------------------------------------------------------------------o

   EFNACMD   Command Code                                      Character

      Value: '5'  EFNA Command identifier.

   EFNAREAS  Answer Reason                                    Numeric(2)

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
             '21'  Invalid file signature.
             '22'  File decryption failure.
             '23'  File decompression failure.
             '99'  Unspecified reason.

             Reason why transmission failed.

   EFNAREASL Answer Reason Text Length                        Numeric(3)

             Length in octets of the field EFNAREAST.

             0 indicates that no EFNAREAST field follows.

   EFNAREAST Answer Reason Text                               [UTF-8](n)

             Reason why transmission failed in plain text.

             It is encoded using [UTF-8].

             Maximum length of the encoded reason is 999 octets.

             No general structure is defined for this attribute.
*/

// EFNA signals that the end of the file command (EFID) could not be handled correctly
type EFNA struct {
	ReasonCode int
	AnswerText string
}

// EFNACMD is the command indicator for the EFNA command.
const EFNACMD = "5"

func (s *EFNA) Command() wire.Command {
	s.AnswerText = wire.TruncateString(s.AnswerText, 999) // TODO: Unicode length

	result := wire.Command{
		Format: s.dataFormat(),
		Data:   []interface{}{EFNACMD, s.ReasonCode, len([]byte(s.AnswerText)), s.AnswerText},
	}

	return result
}

var valuesEFNAREAS = map[int]string{
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
	21: "Invalid file signature.",
	22: "File decryption failure.",
	23: "File decompression failure.",
}

// formatDefinition returns the format definition as given in the RFC5024
func (s *EFNA) formatDefinition() []wire.FormatDefinition {
	return []wire.FormatDefinition{
		{`F X(1)`, `EFNACMD`, &[]wire.Value{{EFNACMD, "EFNA Command"}}},
		{`F 9(2)`, `EFNAREAS`, wire.IntMapToValues(valuesEFNAREAS, 2)},
		{`V 9(3)`, `EFNAREASL`, nil},
		{`V T(n)`, `EFNAREAST`, nil},
	}
}

// dataFormat returns the format definition for this command in a machine
// readable form.
func (s *EFNA) dataFormat() []wire.DataFormat {
	return wire.FormatDefinitionsToDataFormats(s.formatDefinition())
}

func (s *EFNA) Marshal() []byte {
	c := s.Command()
	return c.Marshal()
}

func (s *EFNA) Parse(input []byte) error {
	buffer, err := wire.NewBuffer(&input, EFNACMD)
	if err != nil {
		return err
	}

	s.ReasonCode = buffer.GetNumInt(2)
	lenAnswerText := buffer.GetNumInt(3)
	s.AnswerText = buffer.GetString(lenAnswerText)

	return nil
}

func (s *EFNA) String() string {
	return fmt.Sprintf("EFNA - End File Negative Answer. Reason Code: %d (%s), Text: %s", s.ReasonCode, valuesEFNAREAS[s.ReasonCode], s.AnswerText)
}
