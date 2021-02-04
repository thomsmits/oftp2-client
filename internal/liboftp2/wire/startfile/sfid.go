package startfile

import (
	"time"

	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire"
)

/*
5.3.3.  SFID - Start File

   o-------------------------------------------------------------------o
   |       SFID        Start File                                      |
   |                                                                   |
   |       Start File Phase           Speaker ----> Listener           |
   |-------------------------------------------------------------------|
   | Pos | Field     | Description                           | Format  |
   |-----+-----------+---------------------------------------+---------|
   |   0 | SFIDCMD   | SFID Command, 'H'                     | F X(1)  |
   |   1 | SFIDDSN   | Virtual File Dataset Name             | V X(26) |
   |  27 | SFIDRSV1  | Reserved                              | F X(3)  |
   |  30 | SFIDDATE  | Virtual File Date stamp, (CCYYMMDD)   | V 9(8)  |
   |  38 | SFIDTIME  | Virtual File Time stamp, (HHMMSScccc) | V 9(10) |
   |  48 | SFIDUSER  | User Data                             | V X(8)  |
   |  56 | SFIDDEST  | Destination                           | V X(25) |
   |  81 | SFIDORIG  | Originator                            | V X(25) |
   | 106 | SFIDFMT   | File Format (F/V/U/T)                 | F X(1)  |
   | 107 | SFIDLRECL | Maximum Record Size                   | V 9(5)  |
   | 112 | SFIDFSIZ  | File Size, 1K blocks                  | V 9(13) |
   | 125 | SFIDOSIZ  | Original File Size, 1K blocks         | V 9(13) |
   | 138 | SFIDREST  | Restart Position                      | V 9(17) |
   | 155 | SFIDSEC   | Security Level                        | F 9(2)  |
   | 157 | SFIDCIPH  | Cipher suite selection                | F 9(2)  |
   | 159 | SFIDCOMP  | File compression algorithm            | F 9(1)  |
   | 160 | SFIDENV   | File enveloping format                | F 9(1)  |
   | 161 | SFIDSIGN  | Signed EERP request                   | F X(1)  |
   | 162 | SFIDDESCL | Virtual File Description length       | V 9(3)  |
   | 165 | SFIDDESC  | Virtual File Description              | V T(n)  |
   o-------------------------------------------------------------------o

   SFIDCMD   Command Code                                      Character

      Value: 'H'  SFID Command identifier.

   SFIDDSN   Virtual File Dataset Name                        String(26)

             Dataset name of the Virtual File being transferred,
             assigned by bilateral agreement.

             No general structure is defined for this attribute.

             See Virtual Files - Identification (Section 1.5.2)

   SFIDRSV1  Reserved                                          String(3)

             This field is reserved for future use.

   SFIDDATE  Virtual File Date stamp                          Numeric(8)

     Format: 'CCYYMMDD'  8 decimal digits representing the century,
             year, month, and day.

             Date stamp assigned by the Virtual File's Originator
             indicating when the file was made available for
             transmission.

             See Virtual Files - Identification (Section 1.5.2)

   SFIDTIME  Virtual File Time stamp                         Numeric(10)

     Format: 'HHMMSScccc'  10 decimal digits representing hours,
             minutes, seconds, and a counter (0001-9999), which gives
             higher resolution.

             Time stamp assigned by the Virtual File's Originator
             indicating when the file was made available for
             transmission.

             See Virtual Files - Identification (Section 1.5.2)

   SFIDUSER  User Data                                         String(8)

             May be used by ODETTE-FTP in any way.  If unused, it should
             be initialised to spaces.  It is expected that a bilateral
             agreement exists as to the meaning of the data.

   SFIDDEST  Destination                                      String(25)

     Format: See Identification Code (Section 5.4)

             The Final Recipient of the Virtual File.

             This is the location that will look into the Virtual File
             content and perform mapping functions.  It is also the
             location that creates the End to End Response (EERP)
             command for the received file.

   SFIDORIG  Originator                                       String(25)

     Format: See Identification Code (Section 5.4)

             Originator of the Virtual File.

             It is the location that created (mapped) the data for
             transmission.

   SFIDFMT   File Format                                       Character

      Value: 'F'  Fixed format binary file
             'V'  Variable format binary file
             'U'  Unstructured binary file
             'T'  Text

             Virtual File format.  Used to calculate the restart
             position (Section 1.5.4).

             Once a file has been signed, compressed, and/or encrypted,
             in file format terms it becomes unstructured, format U.
             The record boundaries are no longer discernable until the
             file is decrypted, decompressed, and/or verified.  SFID
             File Format Field in this scenario indicates the format of
             the original file, and the transmitted file must be treated
             as U format.

   SFIDLRECL Maximum Record Size                              Numeric(5)

    Maximum: 99999

             Length in octets of the longest logical record that may be
             transferred to a location.  Only user data is included.

             If SFIDFMT is 'T' or 'U', then this attribute must be set
             to '00000'.

             If SFIDFMT is 'V' and the file is compressed, encrypted, or
             signed, then the maximum value of SFIDRECL is '65536'.

   SFIDFSIZ  Transmitted File Size                           Numeric(13)

    Maximum: 9999999999999

             Space in 1K (1024 octet) blocks required at the Originator
             location to store the actual Virtual File that is to be
             transmitted.

             For example, if a file is compressed before sending, then
             this is the space required to store the compressed file.

             This parameter is intended to provide only a good estimate
             of the Virtual File size.

             Using 13 digits allows for a maximum file size of
             approximately 9.3 PB (petabytes) to be transmitted.

   SFIDOSIZ  Original File Size                              Numeric(13)

    Maximum: 9999999999999

             Space in 1K (1024 octet) blocks required at the Originator
             location to store the original before it was signed,
             compressed, and/or encrypted.

             If no security or compression services have been used,
             SFIDOSIZ should contain the same value as SFIDFSIZ.

             If the original file size is not known, the value zero
             should be used.

             This parameter is intended to provide only a good estimate
             of the original file size.

             The sequence of events in file exchange are:

              (a) raw data file ready to be sent
                   SFIDOSIZ = Original File Size

              (b) signing/compression/encryption

              (c) transmission
                   SFIDFSIZ = Transmitted File Size

              (d) decryption/decompression/verification

              (e) received raw data file for in-house applications
                   SFIDOSIZ = Original File Size

             The Transmitted File Size at (c) indicates to the receiver
             how much storage space is needed to receive the file.

             The Original File Size at (e) indicates to the in-house
             application how much storage space is needed to process the
             file.

   SFIDREST  Restart Position                                Numeric(17)

    Maximum: 99999999999999999

             Virtual File restart position.

             The count represents the:
                - Record Number if SSIDFMT is 'F' or 'V'.
                - File offset in 1K (1024 octet) blocks if SFIDFMT is
                  'U' or 'T'.

             The count will express the transmitted user data (i.e.,
             before ODETTE-FTP buffer compression, header not included).

             After negotiation between adjacent locations,
             retransmission will start at the lowest value.

             Once a file has been signed, compressed, and/or encrypted,
             in file format terms, it has become unstructured, like
             format U.  The file should be treated as format U for the
             purposes of restart, regardless of the actual value in
             SFIDFMT.

   SFIDSEC   Security Level                                   Numeric(2)

      Value: '00'  No security services
             '01'  Encrypted
             '02'  Signed
             '03'  Encrypted and signed

             Indicates whether the file has been signed and/or encrypted
             before transmission. (See Section 6.2.)

   SFIDCIPH  Cipher suite selection                           Numeric(2)

      Value: '00'  No security services
             '01'  See Section 10.2

             Indicates the cipher suite used to sign and/or encrypt the
             file and also to indicate the cipher suite that should be
             used when a signed EERP or NERP is requested.

   SFIDCOMP  File compression algorithm                       Numeric(1)

      Value: '0'  No compression
             '1'  Compressed with [ZLIB] algorithm

             Indicates the algorithm used to compress the file.
             (See Section 6.4.)

   SFIDENV   File enveloping format                           Numeric(1)

      Value: '0'  No envelope
             '1'  File is enveloped using [CMS]

             Indicates the enveloping format used in the file.

             If the file is encrypted/signed/compressed or is an
             enveloped file for the exchange and revocation of
             certificates, this field must be set accordingly.

   SFIDSIGN  Signed EERP request                               Character

      Value: 'Y'  The EERP returned in acknowledgement of the file
                  must be signed
             'N'  The EERP must not be signed

             Requests whether the EERP returned for the file must be
             signed.

   SFIDDESCL Virtual File Description length                  Numeric(3)

             Length in octets of the field SFIDDESC.

             A value of 0 indicates that no description is present.

   SFIDDESC  Virtual File Description                         [UTF-8](n)

             May be used by ODETTE-FTP in any way.  If not used,
             SFIDDESCL should be set to zero.

             No general structure is defined for this attribute, but it
             is expected that a bilateral agreement exists as to the
             meaning of the data.

             It is encoded using [UTF-8] to support a range of national
             languages.

             Maximum length of the encoded value is 999 octets.
*/

// SFID starts file transfer
type SFID struct {
	DatasetName            string
	FileDateTime           time.Time
	UserData               string
	Destination            string
	Originator             string
	FileFormat             string
	MaxRecordSize          int
	FileSizeInK            uint64
	OriginalFileSizeInK    uint64
	RestartPosition        uint64
	SecurityLevel          int
	CipherSuite            int
	Compression            int
	Envelope               int
	SigningRequired        bool
	VirtualFileDescription string
}

// SFIDCMD is the command indicator for the SFID command.
const SFIDCMD = "H"

// formatDefinition returns the format definition as given in the RFC5024
func (s *SFID) formatDefinition() []wire.FormatDefinition {
	return []wire.FormatDefinition{
		{`F X(1)`, `SFIDCMD`, &[]wire.Value{{SFIDCMD, "SFID Command"}}},
		{`V X(26)`, `SFIDDSN`, nil},
		{`F X(3)`, `SFIDRSV1`, &[]wire.Value{{"   ", "Reserved"}}},
		{`V 9(8)`, `SFIDDATE`, nil},
		{`V 9(10)`, `SFIDTIME`, nil},
		{`V X(8)`, `SFIDUSER`, nil},
		{`V X(25)`, `SFIDDEST`, nil},
		{`V X(25)`, `SFIDORIG`, nil},
		{`F X(1)`, `SFIDFMT`, &[]wire.Value{
			{"F", "Fixed format binary file"},
			{"V", "Variable format binary file"},
			{"U", "Unstructured binary file"},
			{"T", "Text"},
		}},
		{`V 9(5)`, `SFIDLRECL`, nil},
		{`V 9(13)`, `SFIDFSIZ`, nil},
		{`V 9(13)`, `SFIDOSIZ`, nil},
		{`V 9(17)`, `SFIDREST`, nil},
		{`F 9(2)`, `SFIDSEC`, &[]wire.Value{
			{"00", "No, security, services"},
			{"01", "Encrypted"},
			{"02", "Signed"},
			{"03", "Encrypted and signed"},
		}},
		{`F 9(2)`, `SFIDCIPH`, &[]wire.Value{
			{"00", "No security services"},
			{"01", " See Section 10.2"},
		}},
		{`F 9(1)`, `SFIDCOMP`, &[]wire.Value{
			{"0", "No compression"},
			{"1", "Compressed with [ZLIB] algorithm"},
		}},
		{`F 9(1)`, `SFIDENV`, &[]wire.Value{
			{"0", "No envelope"},
			{"1", "File is enveloped using [CMS]"},
		}},
		{`F X(1)`, `SFIDSIGN`, wire.ValueBooleanYesNo},
		{`V 9(3)`, `SFIDDESCL`, nil},
		{`V T(n)`, `SFIDDESC`, nil},
	}
}

// dataFormat returns the format definition for this command in a machine
// readable form.
func (s *SFID) dataFormat() []wire.DataFormat {
	return wire.FormatDefinitionsToDataFormats(s.formatDefinition())
}

func (s *SFID) Command() wire.Command {

	fileDate, fileTime := wire.ParseDateToString(s.FileDateTime)

	// truncated virtualFileDescription to maximum 999 characters
	s.VirtualFileDescription = wire.TruncateString(s.VirtualFileDescription, 999) // TODO: Length in case of UTF-8 not correct

	result := wire.Command{
		Format: s.dataFormat(),
		Data: []interface{}{
			SFIDCMD,
			s.DatasetName,
			"",
			fileDate,
			fileTime,
			s.UserData,
			s.Destination,
			s.Originator,
			s.FileFormat,
			s.MaxRecordSize,
			s.FileSizeInK,
			s.OriginalFileSizeInK,
			s.RestartPosition,
			s.SecurityLevel,
			s.CipherSuite,
			s.Compression,
			s.Envelope,
			s.SigningRequired,
			len([]byte(s.VirtualFileDescription)),
			s.VirtualFileDescription,
		},
	}

	return result
}

func (s *SFID) Marshal() []byte {
	c := s.Command()
	return c.Marshal()
}

func (s *SFID) Parse(input []byte) error {
	buffer, err := wire.NewBuffer(&input, SFIDCMD)
	if err != nil {
		return err
	}

	s.DatasetName = buffer.GetString(26)
	_ = buffer.GetString(3)
	fileDate := buffer.GetString(8)
	fileTime := buffer.GetString(10)
	s.UserData = buffer.GetString(8)
	s.Destination = buffer.GetString(25)
	s.Originator = buffer.GetString(25)
	s.FileFormat = buffer.GetString(1)
	s.MaxRecordSize = buffer.GetNumInt(5)
	s.FileSizeInK = buffer.GetBinQWord(13)
	s.OriginalFileSizeInK = buffer.GetBinQWord(13)
	s.RestartPosition = buffer.GetBinQWord(17)
	s.SecurityLevel = buffer.GetNumInt(2)
	s.CipherSuite = buffer.GetNumInt(2)
	s.Compression = buffer.GetNumInt(1)
	s.Envelope = buffer.GetNumInt(1)
	s.SigningRequired = buffer.GetBool(1)
	fileDescriptionLen := buffer.GetNumInt(3)
	s.VirtualFileDescription = buffer.GetString(fileDescriptionLen)

	s.FileDateTime = wire.ParseStringsToDate(fileDate, fileTime)

	return nil
}
