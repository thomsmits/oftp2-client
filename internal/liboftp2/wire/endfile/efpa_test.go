package endfile

import (
	"reflect"
	"testing"
)

func TestEFPA_RoundTrip(t *testing.T) {
	a1 := EFPA{
		ChangeDirection: false,
	}
	b := a1.Marshal()

	a2 := EFPA{}
	err := a2.Parse(b)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(a1, a2) {
		t.Errorf("Roundtrip failed: %v != %v", a1, a2)
	}
}
