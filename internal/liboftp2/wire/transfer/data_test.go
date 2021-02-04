package transfer

import (
	"reflect"
	"testing"
)

func TestDATA_RoundTrip(t *testing.T) {
	a1 := DATA{
		Length: 11,
		Buffer: []byte{0xca, 0xff, 0xee, 0xba, 0xba, 0xbe, 0xde, 0xad, 0xbe, 0xee, 0xff},
	}

	b := a1.Marshal()

	a2 := DATA{
		Length: 11,
	}

	err := a2.Parse(b)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(a1, a2) {
		t.Errorf("Roundtrip failed: %v != %v", a1, a2)
	}
}
