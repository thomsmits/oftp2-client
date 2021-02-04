package wire

import (
	"reflect"
	"testing"
)

func TestCD_RoundTrip(t *testing.T) {
	a1 := CD{}

	b := a1.Marshal()

	a2 := CD{}

	err := a2.Parse(b)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(a1, a2) {
		t.Errorf("Roundtrip failed: %v != %v", a1, a2)
	}
}
