package startfile

import (
	"reflect"
	"testing"
	"time"

	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire"
)

func TestEERP_RoundTrip(t *testing.T) {

	date, tme := wire.ParseDateToString(time.Now())
	now := wire.ParseStringsToDate(date, tme)

	a1 := EERP{
		VirtualDataSetName: "DATASETNAME",
		VirtualFileDate:    now,
		UserData:           "WHATEVER",
		Destination:        "JDHDJD",
		Originator:         "JDJDIUII",
		FileHash:           []byte{0xca, 0xff, 0xee, 0xba, 0xba, 0xbe},
		Signature:          []byte{0xde, 0xad, 0xbe, 0xee, 0xff},
	}

	b := a1.Marshal()

	a2 := EERP{}
	err := a2.Parse(b)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(a1, a2) {
		t.Errorf("Roundtrip failed: %v != %v", a1, a2)
	}
}
