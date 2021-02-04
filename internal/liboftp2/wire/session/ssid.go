package session

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire"
)

/*
5.3.2.  SSID - Start Session

   o-------------------------------------------------------------------o
   |       SSID        Start Session                                   |
   |                                                                   |
   |       Start Session Phase     Initiator <---> Responder           |
   |-------------------------------------------------------------------|
   | Pos | Field     | Description                           | Format  |
   |-----+-----------+---------------------------------------+---------|
   |   0 | SSIDCMD   | SSID Command 'X'                      | F X(1)  |
   |   1 | SSIDLEV   | Protocol Release Level                | F 9(1)  |
   |   2 | SSIDCODE  | Initiator's Identification Code       | V X(25) |
   |  27 | SSIDPSWD  | Initiator's Password                  | V X(8)  |
   |  35 | SSIDSDEB  | Data Exchange Buffer Size             | V 9(5)  |
   |  40 | SSIDSR    | Send / Receive Capabilities (S/R/B)   | F X(1)  |
   |  41 | SSIDCMPR  | Buffer Compression Indicator (Y/N)    | F X(1)  |
   |  42 | SSIDREST  | Restart Indicator (Y/N)               | F X(1)  |
   |  43 | SSIDSPEC  | Special Logic Indicator (Y/N)         | F X(1)  |
   |  44 | SSIDCRED  | Credit                                | V 9(3)  |
   |  47 | SSIDAUTH  | Secure Authentication (Y/N)           | F X(1)  |
   |  48 | SSIDRSV1  | Reserved                              | F X(4)  |
   |  52 | SSIDUSER  | User Data                             | V X(8)  |
   |  60 | SSIDCR    | Carriage Return                       | F X(1)  |
   o-------------------------------------------------------------------o

      SSIDCMD   Command Code
      Character

      Value: 'X'  SSID Command identifier.

   SSIDLEV   Protocol Release Level                           Numeric(1)

             Used to specify the level of the ODETTE-FTP protocol

      Value: '1' for Revision 1.2
             '2' for Revision 1.3
             '4' for Revision 1.4
             '5' for Revision 2.0

             Future release levels will have higher numbers.  The
             protocol release level is negotiable, with the lowest level
             being selected.

             Note: ODETTE File Transfer Protocol 1.3 (RFC 2204)
                   specifies '1' for the release level, despite adhering
                   to revision 1.3.

   SSIDCODE  Initiator's Identification Code                  String(25)

    Format:  See Identification Code (Section 5.4)

             Uniquely identifies the Initiator (sender) participating in
             the ODETTE-FTP session.

             It is an application implementation issue to link the
             expected [X.509] certificate to the SSIDCODE provided.

   SSIDPSWD  Initiator's Password                              String(8)

             Key to authenticate the sender.  Assigned by bilateral
             agreement.

   SSIDSDEB  Data Exchange Buffer Size                        Numeric(5)

    Minimum: 128
    Maximum: 99999

             The length, in octets, of the largest Data Exchange Buffer
             that can be accepted by the location.  The length includes
             the command octet but does not include the Stream
             Transmission Header.

             After negotiation, the smallest size will be selected.

   SSIDSR    Send / Receive Capabilities                       Character

      Value: 'S'  Location can only send files.
             'R'  Location can only receive files.
             'B'  Location can both send and receive files.

             Sending and receiving will be serialised during the
             session, so parallel transmissions will not take place in
             the same session.

             An error occurs if adjacent locations both specify the send
             or receive capability.

   SSIDCMPR  Buffer Compression Indicator                      Character

      Value: 'Y'  The location can handle OFTP data buffer compression
             'N'  The location cannot handle OFTP buffer compression

             Compression is only used if supported by both locations.

             The compression mechanism referred to here applies to each
             individual OFTP data buffer.  This is different from the
             file compression mechanism in OFTP, which involves the
             compression of whole files.

   SSIDREST  Restart Indicator                                 Character

      Value: 'Y'  The location can handle the restart of a partially
                  transmitted file.
             'N'  The location cannot restart a file.

   SSIDSPEC  Special Logic Indicator                           Character

      Value: 'Y'  Location can handle Special Logic
             'N'  Location cannot handle Special Logic

             Special Logic is only used if supported by both locations.

             The Special Logic extensions are only useful to access an
             X.25 network via an asynchronous entry and are not
             supported for TCP/IP connections.

   SSIDCRED  Credit                                           Numeric(3)

    Maximum: 999

             The number of consecutive Data Exchange Buffers sent by the
             Speaker before it must wait for a Credit (CDT) command from
             the Listener.

             The credit value is only applied to Data flow in the Data
             Transfer phase.

             The Speaker's available credit is initialised to SSIDCRED
             when it receives a Start File Positive Answer (SFPA)
             command from the Listener.  It is zeroed by the End File
             (EFID) command.

             After negotiation, the smallest size must be selected in
             the answer of the Responder, otherwise a protocol error
             will abort the session.

             Negotiation of the "credit-window-size" parameter.

             Window Size m  -- SSID ------------>
                            <------------ SSID --  Window Size n
                                                  (n less than or
                                                   equal to m)
             Note: negotiated value will be "n".

   SSIDAUTH  Secure Authentication                             Character

      Value: 'Y'  The location requires secure authentication.  'N'  The
             location does not require secure authentication.

             Secure authentication is only used if agreed by both
             locations.

             If the answer of the Responder does not match with the
             authentication requirements of the Initiator, then the
             Initiator must abort the session.

             No negotiation of authentication is allowed.

             authentication p  -- SSID ------------>
                               <------------ SSID --  authentication q

             p == q -> continue.
             p != q -> abort.

   SSIDRSV1  Reserved                                          String(4)

             This field is reserved for future use.

   SSIDUSER  User Data                                         String(8)

             May be used by ODETTE-FTP in any way.  If unused, it should
             be initialised to spaces.  It is expected that a bilateral
             agreement exists as to the meaning of the data.

   SSIDCR    Carriage Return                                   Character

      Value: Character with hex value '0D' or '8D'.
*/

// SSID starts the session
type SSID struct {
	Id             string
	Password       string
	BufferSize     uint32
	Capability     string
	Compress       bool
	Restart        bool
	Special        bool
	Credit         uint32
	Authentication bool
	UserData       string
}

var valuesSSIDSR = map[string]string{
	"S": "Location can only send files.",
	"R": "Location can only receive files.",
	"B": "Location can both send and receive files.",
}

var valuesSSIDLEV = map[string]string{
	"1": "Revision 1.2",
	"2": "Revision 1.3",
	"4": "Revision 1.4",
	"5": "Revision 2.0",
}

// SSIDCMD is the command indicator for the SSID command.
const SSIDCMD = "X"

// formatDefinition returns the format definition as given in the RFC5024
func (s *SSID) formatDefinition() []wire.FormatDefinition {
	return []wire.FormatDefinition{
		{`F X(1)`, `SSIDCMD`, &[]wire.Value{{SSIDCMD, "SSID Command"}}},
		{`F 9(1)`, `SSIDLEV`, wire.StringMapToValues(valuesSSIDLEV)},
		{`V X(25)`, `SSIDCODE`, nil},
		{`V X(8)`, `SSIDPSWD`, nil},
		{`V 9(5)`, `SSIDSDEB`, nil},
		{`F X(1)`, `SSIDSR`, wire.StringMapToValues(valuesSSIDSR)},
		{`F X(1)`, `SSIDCMPR`, wire.ValueBooleanYesNo},
		{`F X(1)`, `SSIDREST`, wire.ValueBooleanYesNo},
		{`F X(1)`, `SSIDSPEC`, wire.ValueBooleanYesNo},
		{`V 9(3)`, `SSIDCRED`, nil},
		{`F X(1)`, `SSIDAUTH`, wire.ValueBooleanYesNo},
		{`F X(4)`, `SSIDRSV1`, nil},
		{`V X(8)`, `SSIDUSER`, nil},
		{`F X(1)`, `SSIDCR`, wire.ValueNewline},
	}
}

// dataFormat returns the format definition for this command in a machine
// readable form.
func (s *SSID) dataFormat() []wire.DataFormat {
	return wire.FormatDefinitionsToDataFormats(s.formatDefinition())
}

func (s *SSID) Marshal() []byte {
	cmd := s.Command()
	return cmd.Marshal()
}

func (s *SSID) Command() wire.Command {
	return wire.Command{
		Format: s.dataFormat(),
		Data: []interface{}{
			SSIDCMD,
			"5",
			s.Id,
			s.Password,
			s.BufferSize,
			s.Capability,
			s.Compress,
			s.Restart,
			s.Special,
			s.Credit,
			s.Authentication,
			"",
			s.UserData,
			"\n",
		},
	}
}

func (s *SSID) Parse(input []byte) error {
	buffer, err := wire.NewBuffer(&input, SSIDCMD)
	if err != nil {
		return err
	}

	version := buffer.GetString(1)

	if version != "5" {
		return errors.New(fmt.Sprintf("wrong version %s, expected 5", version))
	}

	s.Id = buffer.GetString(25)
	s.Password = buffer.GetString(8)
	s.BufferSize = buffer.GetBinDWord(5)
	s.Capability = buffer.GetString(1)
	s.Compress = buffer.GetBool(1)
	s.Restart = buffer.GetBool(1)
	s.Special = buffer.GetBool(1)
	s.Credit = buffer.GetBinDWord(3)
	s.Authentication = buffer.GetBool(1)
	_ = buffer.GetString(4)
	s.UserData = buffer.GetString(8)

	return nil
}

func (s *SSID) String() string {
	result := ""
	result += fmt.Sprintf("OdetteId                : %s\n", s.Id)
	result += fmt.Sprintf("Max buffer size         : %d\n", s.BufferSize)
	result += fmt.Sprintf("Capability              : %s\n", valuesSSIDSR[s.Capability])
	result += fmt.Sprintf("Compression supported   : %t\n", s.Compress)
	result += fmt.Sprintf("Restart supported       : %t\n", s.Restart)
	result += fmt.Sprintf("Special logic supported : %t\n", s.Special)
	result += fmt.Sprintf("Authentication supported: %t\n", s.Authentication)
	result += fmt.Sprintf("User data               : %s\n", s.UserData)

	return result
}

func ToSSID(data *interface{}) *SSID {
	switch (*data).(type) {
	case SSID:
		ssid := (*data).(SSID)
		return &ssid
	case *SSID:
		ssid := (*data).(*SSID)
		return ssid
	default:
		panic(fmt.Sprintf("wrong type of data. Expected %s, got %s", "SSID", reflect.TypeOf(data)))
		return nil // never reached
	}
}
