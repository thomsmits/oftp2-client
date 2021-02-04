package startfile

import (
	"fmt"
	"time"

	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire"
)

/*
5.3.14.  NERP - Negative End Response

   o-------------------------------------------------------------------o
   |       NERP        Negative End Response                           |
   |                                                                   |
   |       Start File Phase           Speaker ----> Listener           |
   |       End File Phase             Speaker ----> Listener           |
   |-------------------------------------------------------------------|
   | Pos | Field     | Description                           | Format  |
   |-----+-----------+---------------------------------------+---------|
   |   0 | NERPCMD   | NERP Command, 'N'                     | F X(1)  |
   |   1 | NERPDSN   | Virtual File Dataset Name             | V X(26) |
   |  27 | NERPRSV1  | Reserved                              | F X(6)  |
   |  33 | NERPDATE  | Virtual File Date stamp, (CCYYMMDD)   | V 9(8)  |
   |  41 | NERPTIME  | Virtual File Time stamp, (HHMMSScccc) | V 9(10) |
   |  51 | NERPDEST  | Destination                           | V X(25) |
   |  76 | NERPORIG  | Originator                            | V X(25) |
   | 101 | NERPCREA  | Creator of NERP                       | V X(25) |
   | 126 | NERPREAS  | Reason code                           | F 9(2)  |
   | 128 | NERPREASL | Reason text length                    | V 9(3)  |
   | 131 | NERPREAST | Reason text                           | V T(n)  |
   |     | NERPHSHL  | Virtual File hash length              | V U(2)  |
   |     | NERPHSH   | Virtual File hash                     | V U(n)  |
   |     | NERPSIGL  | NERP signature length                 | V U(2)  |
   |     | NERPSIG   | NERP signature                        | V U(n)  |
   o-------------------------------------------------------------------o

   NERPCMD   Command Code                                      Character

      Value: 'N'  NERP Command identifier.

   NERPDSN   Virtual File Dataset Name                        String(26)

             Dataset name of the Virtual File being transferred,
             assigned by bilateral agreement.

             No general structure is defined for this attribute.

             See Virtual Files - Identification (Section 1.5.2)

   NERPRSV1  Reserved                                          String(6)

             This field is reserved for future use.

   NERPDATE  Virtual File Date stamp                          Numeric(8)

     Format: 'CCYYMMDD'  8 decimal digits representing the century,
             year, month, and day, respectively.

             Date stamp assigned by the Virtual File's Originator
             indicating when the file was made available for
             transmission.

             See Virtual Files - Identification (Section 1.5.2)

   NERPTIME  Virtual File Time stamp                         Numeric(10)

     Format: 'HHMMSScccc'  10 decimal digits representing hours,
             minutes, seconds, and a counter (0001-9999), which gives
             higher resolution.

             Time stamp assigned by the Virtual File's Originator
             indicating when the file was made available for
             transmission.

             See Virtual Files - Identification (Section 1.5.2)

   NERPDEST  Destination                                      String(25)

     Format: See Identification Code (Section 5.4)

             Originator of the Virtual File.

             This is the location that created the data for
             transmission.

   NERPORIG  Originator                                       String(25)

     Format: See Identification Code (Section 5.4)

             The Final Recipient of the Virtual File.

             This is the location that will look into the Virtual File
             content and perform mapping functions.

   NERPCREA  Creator of the NERP                              String(25)

     Format: See Identification Code (Section 5.4)

             It is the location that created the NERP.

   NERPREAS  Reason code                                      Numeric(2)

             This attribute will specify why transmission cannot proceed
             or why processing of the file failed.

             "SFNA(RETRY=N)" below should be interpreted as "EFNA or
             SFNA(RETRY=N)" where appropriate.

      Value  '03'  ESID received with reason code '03'
                    (user code not known)
             '04'  ESID received with reason code '04'
                    (invalid password)
             '09'  ESID received with reason code '99'
                    (unspecified reason)
             '11'  SFNA(RETRY=N) received with reason code '01'
                    (invalid file name)
             '12'  SFNA(RETRY=N) received with reason code '02'
                    (invalid destination)
             '13'  SFNA(RETRY=N) received with reason code '03'
                    (invalid origin)
             '14'  SFNA(RETRY=N) received with reason code '04'
                    (invalid storage record format)
             '15'  SFNA(RETRY=N) received with reason code '05'
                    (maximum record length not supported)
             '16'  SFNA(RETRY=N) received with reason code '06'
                    (file size too big)
             '20'  SFNA(RETRY=N) received with reason code '10'
                    (invalid record count)
             '21'  SFNA(RETRY=N) received with reason code '11'
                    (invalid byte count)
             '22'  SFNA(RETRY=N) received with reason code '12'
                    (access method failure)
             '23'  SFNA(RETRY=N) received with reason code '13'
                    (duplicate file)
             '24'  SFNA(RETRY=N) received with reason code '14'
                    (file direction refused)
             '25'  SFNA(RETRY=N) received with reason code '15'
                    (cipher suite not supported)
             '26'  SFNA(RETRY=N) received with reason code '16'
                    (encrypted file not allowed)
             '27'  SFNA(RETRY=N) received with reason code '17'
                    (unencrypted file not allowed)
             '28'  SFNA(RETRY=N) received with reason code '18'
                    (compression not allowed)
             '29'  SFNA(RETRY=N) received with reason code '19'
                    (signed file not allowed)
             '30'  SFNA(RETRY=N) received with reason code '20'
                    (unsigned file not allowed)
             '31'  File signature not valid.
             '32'  File decompression failed.
             '33'  File decryption failed.
             '34'  File processing failed.
             '35'  Not delivered to recipient.
             '36'  Not acknowledged by recipient.
             '50'  Transmission stopped by the operator.
             '90'  File size incompatible with recipient's
                    protocol version.
             '99'  Unspecified reason.

   NERPREASL Reason Text Length                              Numeric(3)

             Length in octets of the field NERPREAST.

             0 indicates that no NERPREAST field follows.

   NERPREAST Reason Text                                     [UTF-8](n)

             Reason why transmission cannot proceed in plain text.

             It is encoded using [UTF-8].

             Maximum length of the encoded reason is 999 octets.

             No general structure is defined for this attribute.

   NERPHSHL  Virtual File hash length                          Binary(2)

             Length in octets of the field NERPHSH.

             A binary value of 0 indicates that no hash is present.
             This is always the case if the NERP is not signed.

   NERPHSH   Virtual File hash                                 Binary(n)

             Hash of the Virtual File being transmitted.

             The algorithm used is determined by the bilaterally agreed
             cipher suite specified in the SFIDCIPH.

   NERPSIGL  NERP Signature length                             Binary(2)

             0 indicates that this NERP has not been signed.

             Any other value indicates the length of NERPSIG in octets
             and indicates that this NERP has been signed.

   NERPSIG   NERP Signature                                    Binary(n)

             Contains the [CMS] enveloped signature of the NERP.

             Signature = Sign{NERPDSN
                              NERPDATE
                              NERPTIME
                              NERPDEST
                              NERPORIG
                              NERPCREA
                              NERPHSH}

             Each field is taken in its entirety, including any padding.
             The envelope must contain the original data, not just the
             signature.

             The [CMS] content type used is SignedData.

             The encapsulated content type used is id-data.

             It is an application issue to validate the signature with
             the contents of the NERP.
*/

// NERP is a negative end to end response
type NERP struct {
	VirtualDataSetName string
	VirtualFileDate    time.Time
	Destination        string
	Originator         string
	CreatorOfNERP      string
	ReasonCode         int
	ReasonText         string
	FileHash           []byte
	Signature          []byte
}

// NERPCMD is the command indicator for the NERP command.
const NERPCMD = "N"

var valuesNERPREAS = map[int]string{
	3:  "ESID received with reason code '03'",
	4:  "ESID received with reason code '04'",
	9:  "ESID received with reason code '99'",
	11: "SFNA(RETRY=N) received with reason code '01'",
	12: "SFNA(RETRY=N) received with reason code '02'",
	13: "SFNA(RETRY=N) received with reason code '03'",
	14: "SFNA(RETRY=N) received with reason code '04'",
	15: "SFNA(RETRY=N) received with reason code '05'",
	16: "SFNA(RETRY=N) received with reason code '06'",
	20: "SFNA(RETRY=N) received with reason code '10'",
	21: "SFNA(RETRY=N) received with reason code '11'",
	22: "SFNA(RETRY=N) received with reason code '12'",
	23: "SFNA(RETRY=N) received with reason code '13'",
	24: "SFNA(RETRY=N) received with reason code '14'",
	25: "SFNA(RETRY=N) received with reason code '15'",
	26: "SFNA(RETRY=N) received with reason code '16'",
	27: "SFNA(RETRY=N) received with reason code '17'",
	28: "SFNA(RETRY=N) received with reason code '18'",
	29: "SFNA(RETRY=N) received with reason code '19'",
	30: "SFNA(RETRY=N) received with reason code '20'",
	31: "File signature not valid.",
	32: "File decompression failed.",
	33: "File decryption failed.",
	34: "File processing failed.",
	35: "Not delivered to recipient.",
	36: "Not acknowledged by recipient.",
	50: "Transmission stopped by the operator.",
	90: "File size incompatible with recipient's protocol version.",
	99: "Unspecified reason.",
}

// formatDefinition returns the format definition as given in the RFC5024
func (s *NERP) formatDefinition() []wire.FormatDefinition {
	return []wire.FormatDefinition{
		{`F X(1)`, `NERPCMD `, &[]wire.Value{{NERPCMD, "NERP Command"}}},
		{`V X(26)`, `NERPDSN `, nil},
		{`F X(6)`, `NERPRSV1`, &[]wire.Value{{"      ", "Reserved"}}},
		{`V 9(8)`, `NERPDATE`, nil},
		{`V 9(10)`, `NERPTIME`, nil},
		{`V X(25)`, `NERPDEST`, nil},
		{`V X(25)`, `NERPORIG`, nil},
		{`V X(25)`, `NERPCREA`, nil},
		{`F 9(2)`, `NERPREAS`, wire.IntMapToValues(valuesNERPREAS, 2)},
		{`V 9(3)`, `NERPREASL`, nil},
		{`V T(n)`, `NERPREAST`, nil},
		{`V U(2)`, `NERPHSHL`, nil},
		{`V U(n)`, `NERPHSH`, nil},
		{`V U(2)`, `NERPSIGL`, nil},
		{`V U(n)`, `NERPSIG`, nil},
	}
}

// dataFormat returns the format definition for this command in a machine
// readable form.
func (s *NERP) dataFormat() []wire.DataFormat {
	return wire.FormatDefinitionsToDataFormats(s.formatDefinition())
}

func (s *NERP) Command() wire.Command {

	fileDate, fileTime := wire.ParseDateToString(s.VirtualFileDate)

	// TODO: Make Enumeration for reasonCode

	result := wire.Command{
		Format: s.dataFormat(),
		Data: []interface{}{
			NERPCMD,
			s.VirtualDataSetName,
			"",
			fileDate,
			fileTime,
			s.Destination,
			s.Originator,
			s.CreatorOfNERP,
			s.ReasonCode,
			len(s.ReasonText),
			s.ReasonText,
			len(s.FileHash),
			s.FileHash,
			len(s.Signature),
			s.Signature,
		},
	}

	return result
}

func (s *NERP) Marshal() []byte {
	c := s.Command()
	return c.Marshal()
}

func (s *NERP) Parse(input []byte) error {
	buffer, err := wire.NewBuffer(&input, NERPCMD)
	if err != nil {
		return err
	}

	s.VirtualDataSetName = buffer.GetString(26)
	_ = buffer.GetString(6)
	fileDate := buffer.GetString(8)
	fileTime := buffer.GetString(10)
	s.Destination = buffer.GetString(25)
	s.Originator = buffer.GetString(25)
	s.CreatorOfNERP = buffer.GetString(25)
	s.ReasonCode = buffer.GetNumInt(2)
	lenReasonText := buffer.GetNumInt(3)
	s.ReasonText = buffer.GetString(lenReasonText)
	lenFileHash := buffer.GetBinWord(2)
	s.FileHash = buffer.GetBytes(int(lenFileHash))
	lenSignature := buffer.GetBinWord(2)
	s.Signature = buffer.GetBytes(int(lenSignature))

	s.VirtualFileDate = wire.ParseStringsToDate(fileDate, fileTime)
	return nil
}

func (s *NERP) String() string {
	return fmt.Sprintf("NERP - Negative End Response. Reason Code: %d (%s), Text: %s", s.ReasonCode, valuesNERPREAS[s.ReasonCode], s.ReasonText)
}
