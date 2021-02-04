package startfile

import (
	"reflect"
	"testing"
)

func TestRTR_RoundTrip(t *testing.T) {

	a1 := RTR{}

	b := a1.Marshal()

	a2 := RTR{}
	err := a2.Parse(b)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(a1, a2) {
		t.Errorf("Roundtrip failed: %v != %v", a1, a2)
	}
}
