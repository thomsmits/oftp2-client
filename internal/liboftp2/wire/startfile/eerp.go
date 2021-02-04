package startfile

import (
	"time"

	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire"
)

/*
5.3.13.  EERP - End to End Response

   o-------------------------------------------------------------------o
   |       EERP        End to End Response                             |
   |                                                                   |
   |       Start File Phase           Speaker ----> Listener           |
   |       End File Phase             Speaker ----> Listener           |
   |-------------------------------------------------------------------|
   | Pos | Field     | Description                           | Format  |
   |-----+-----------+---------------------------------------+---------|
   |   0 | EERPCMD   | EERP Command, 'E'                     | F X(1)  |
   |   1 | EERPDSN   | Virtual File Dataset Name             | V X(26) |
   |  27 | EERPRSV1  | Reserved                              | F X(3)  |
   |  30 | EERPDATE  | Virtual File Date stamp, (CCYYMMDD)   | V 9(8)  |
   |  38 | EERPTIME  | Virtual File Time stamp, (HHMMSScccc) | V 9(10) |
   |  48 | EERPUSER  | User Data                             | V X(8)  |
   |  56 | EERPDEST  | Destination                           | V X(25) |
   |  81 | EERPORIG  | Originator                            | V X(25) |
   | 106 | EERPHSHL  | Virtual File hash length              | V U(2)  |
   | 108 | EERPHSH   | Virtual File hash                     | V U(n)  |
   |     | EERPSIGL  | EERP signature length                 | V U(2)  |
   |     | EERPSIG   | EERP signature                        | V U(n)  |
   o-------------------------------------------------------------------o

   EERPCMD   Command Code                                      Character

      Value: 'E'  EERP Command identifier.

   EERPDSN   Virtual File Dataset Name                        String(26)

             Dataset name of the Virtual File being transferred,
             assigned by bilateral agreement.

             No general structure is defined for this attribute.

             See Virtual Files - Identification (Section 1.5.2)

   EERPRSV1  Reserved                                          String(3)

             This field is reserved for future use.

   EERPDATE  Virtual File Date stamp                          Numeric(8)

     Format: 'CCYYMMDD'  8 decimal digits representing the century,
             year, month, and day, respectively.

             Date stamp assigned by the Virtual File's Originator
             indicating when the file was made available for
             transmission.

             See Virtual Files - Identification (Section 1.5.2)

   EERPTIME  Virtual File Time stamp                         Numeric(10)

     Format: 'HHMMSScccc'  10 decimal digits representing hours,
             minutes, seconds, and a counter (0001-9999), which gives
             higher resolution.

             Time stamp assigned by the Virtual File's Originator
             indicating when the file was made available for
             transmission.

             See Virtual Files - Identification (Section 1.5.2)

   EERPUSER  User Data                                         String(8)

             May be used by ODETTE-FTP in any way.  If unused, it should
             be initialised to spaces.  It is expected that a bilateral
             agreement exists as to the meaning of the data.

   EERPDEST  Destination                                      String(25)

     Format: See Identification Code (Section 5.4)

             Originator of the Virtual File.

             This is the location that created the data for
             transmission.

   EERPORIG  Originator                                       String(25)

     Format: See Identification Code (Section 5.4)

             Final Recipient of the Virtual File.

             This is the location that will look into the Virtual File
             content and process it accordingly.  It is also the
             location that creates the EERP for the received file.

   EERPHSHL  Virtual File hash length                          Binary(2)

             Length in octets of the field EERPHSH.

             A binary value of 0 indicates that no hash is present.
             This is always the case if the EERP is not signed.

   EERPHSH   Virtual File hash                                 Binary(n)

             Hash of the transmitted Virtual File, i.e., not the hash of
             the original file.

             The algorithm used is determined by the bilaterally agreed
             cipher suite specified in the SFIDCIPH.

             It is an application implementation issue to validate the
             EERPHSH to ensure that the EERP is acknowledging the exact
             same file as was originally transmitted.

   EERPSIGL  EERP signature length                             Binary(2)

             0 indicates that this EERP has not been signed.

             Any other value indicates the length of EERPSIG in octets
             and indicates that this EERP has been signed.

   EERPSIG   EERP signature                                    Binary(n)

             Contains the [CMS] enveloped signature of the EERP.

             Signature = Sign{EERPDSN
                              EERPDATE
                              EERPTIME
                              EERPDEST
                              EERPORIG
                              EERPHSH}

             Each field is taken in its entirety, including any padding.
             The envelope must contain the original data, not just the
             signature.

             The [CMS] content type used is SignedData.

             The encapsulated content type used is id-data.

             It is an application issue to validate the signature with
             the contents of the EERP.
*/

// EERP represents an end to end response
type EERP struct {
	VirtualDataSetName string
	VirtualFileDate    time.Time
	UserData           string
	Destination        string
	Originator         string
	FileHash           []byte
	Signature          []byte
}

// EERPCMD is the command indicator for the EERP command.
const EERPCMD = "E"

// formatDefinition returns the format definition as given in the RFC5024
func (s *EERP) formatDefinition() []wire.FormatDefinition {
	return []wire.FormatDefinition{
		{`F X(1)`, `EERPCMD`, &[]wire.Value{{EERPCMD, "EERP Command"}}},
		{`V X(26)`, `EERPDSN`, nil},
		{`F X(3)`, `EERPRSV1`, &[]wire.Value{{"   ", "Reserved"}}},
		{`V 9(8)`, `EERPDATE`, nil},
		{`V 9(10)`, `EERPTIME`, nil},
		{`V X(8)`, `EERPUSER`, nil},
		{`V X(25)`, `EERPDEST`, nil},
		{`V X(25)`, `EERPORIG`, nil},
		{`V U(2)`, `EERPHSHL`, nil},
		{`V U(n)`, `EERPHSH`, nil},
		{`V U(2)`, `EERPSIGL`, nil},
		{`V U(n)`, `EERPSIG`, nil},
	}
}

// dataFormat returns the format definition for this command in a machine
// readable form.
func (s *EERP) dataFormat() []wire.DataFormat {
	return wire.FormatDefinitionsToDataFormats(s.formatDefinition())
}

func (s *EERP) Command() wire.Command {

	fileDate, fileTime := wire.ParseDateToString(s.VirtualFileDate)

	result := wire.Command{
		Format: s.dataFormat(),
		Data: []interface{}{
			EERPCMD,
			s.VirtualDataSetName,
			"   ",
			fileDate,
			fileTime,
			s.UserData,
			s.Destination,
			s.Originator,
			len(s.FileHash),
			s.FileHash,
			len(s.Signature),
			s.Signature,
		},
	}

	return result
}

func (s *EERP) Marshal() []byte {
	c := s.Command()
	return c.Marshal()
}

func (s *EERP) Parse(input []byte) error {
	buffer, err := wire.NewBuffer(&input, EERPCMD)
	if err != nil {
		return err
	}

	s.VirtualDataSetName = buffer.GetString(26)
	_ = buffer.GetString(3)
	fileDate := buffer.GetString(8)
	fileTime := buffer.GetString(10)
	s.UserData = buffer.GetString(8)
	s.Destination = buffer.GetString(25)
	s.Originator = buffer.GetString(25)
	lenFileHash := buffer.GetBinWord(2)
	s.FileHash = buffer.GetBytes(int(lenFileHash))
	lenSignature := buffer.GetBinWord(2)
	s.Signature = buffer.GetBytes(int(lenSignature))

	s.VirtualFileDate = wire.ParseStringsToDate(fileDate, fileTime)

	return nil
}
