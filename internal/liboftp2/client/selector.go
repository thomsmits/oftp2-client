package client

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire"
	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire/authentication"
	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire/endfile"
	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire/session"
	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire/startfile"
	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire/transfer"
)

// DetermineMessageType determines the type of message found in the input buffer and returns
// the corresponding data structure and the message type as a string.
func DetermineMessageType(input []byte) (wire.Protocol, string, error) {

	indicator := string(input[0:1])

	var p wire.Protocol
	var err error

	switch indicator {
	case session.ESIDCMD: // "F"
		p = &session.ESID{}
		err = p.Parse(input)
	case session.SSIDCMD: // "X"
		p = &session.SSID{}
		err = p.Parse(input)
	case session.SSRMCMD: // "I"
		p = &session.SSRM{}
		err = p.Parse(input)
	case authentication.AUCHCMD: // "A"
		p = &authentication.AUCH{}
		err = p.Parse(input)
	case authentication.AURPCMD: // "S"
		p = &authentication.AURP{}
		err = p.Parse(input)
	case authentication.SECDCMD: // "J"
		p = &authentication.SECD{}
		err = p.Parse(input)
	case startfile.EERPCMD: // "E"
		p = &startfile.EERP{}
		err = p.Parse(input)
	case startfile.SFNACMD: // "3"
		p = &startfile.SFNA{}
		err = p.Parse(input)
	case startfile.NERPCMD: // "N"
		p = &startfile.NERP{}
		err = p.Parse(input)
	case startfile.RTRCMD: // "P"
		p = &startfile.RTR{}
		err = p.Parse(input)
	case startfile.SFIDCMD: // "H"
		p = &startfile.SFID{}
		err = p.Parse(input)
	case startfile.SFPACMD: // "2"
		p = &startfile.SFPA{}
		err = p.Parse(input)
	case transfer.DATACMD: // "D"
		p = &transfer.DATA{}
		err = p.Parse(input)
	case transfer.CDTCMD: // "C"
		p = &transfer.CDT{}
		err = p.Parse(input)
	case endfile.EFIDCMD: // "T"
		p = &endfile.EFID{}
		err = p.Parse(input)
	case endfile.EFNACMD: // "5"
		p = &endfile.EFNA{}
		err = p.Parse(input)
	case endfile.EFPACMD: // "4"
		p = &endfile.EFPA{}
		err = p.Parse(input)
	case wire.CDCMD: // "R"
		p = &wire.CD{}
		err = p.Parse(input)
	default:
		return nil, "", errors.New(fmt.Sprintf("unknown start of message %s", indicator))
	}

	r := regexp.MustCompile(`.*\.(.+)`)
	typeName := reflect.TypeOf(p).String()
	typeName = r.FindStringSubmatch(typeName)[1]
	return p, typeName, err
}
