package startfile

import (
	"reflect"
	"testing"
	"time"

	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire"
)

func TestSFID_RoundTrip(t *testing.T) {

	date, tme := wire.ParseDateToString(time.Now())
	now := wire.ParseStringsToDate(date, tme)

	a1 := SFID{
		DatasetName:            "DATASET",
		FileDateTime:           now,
		UserData:               "USEDATA",
		Destination:            "O292929",
		Originator:             "O181811",
		FileFormat:             "F",
		MaxRecordSize:          100,
		FileSizeInK:            1000,
		OriginalFileSizeInK:    10000,
		RestartPosition:        10,
		SecurityLevel:          1,
		CipherSuite:            1,
		Compression:            1,
		Envelope:               1,
		SigningRequired:        true,
		VirtualFileDescription: "VIRTUAL DESCR",
	}

	b := a1.Marshal()

	a2 := SFID{}
	err := a2.Parse(b)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(a1, a2) {
		t.Errorf("Roundtrip failed: %v != %v", a1, a2)
	}
}
