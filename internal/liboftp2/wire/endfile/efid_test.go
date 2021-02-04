package endfile

import (
	"testing"
)

func TestEFID_RoundTrip(t *testing.T) {
	a1 := EFID{
		RecordCount: 10,
		UnitCount:   5,
	}
	b := a1.Marshal()

	a2 := EFID{}
	err := a2.Parse(b)
	if err != nil {
		t.Error(err)
	}

	if a1.RecordCount != a2.RecordCount || a1.UnitCount != a2.UnitCount {
		t.Errorf("Roundtrip failed: %v != %v", a1, a2)
	}
}
