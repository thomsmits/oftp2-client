package session

import (
	"reflect"
	"testing"
)

func TestESID_RoundTrip(t *testing.T) {
	a1 := ESID{
		ReasonCode: 2,
		ReasonText: "A text explaining, what went wrong and why - or maybe not",
	}
	b := a1.Marshal()

	a2 := ESID{}
	err := a2.Parse(b)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(a1, a2) {
		t.Errorf("Roundtrip failed: %v != %v", a1, a2)
	}
}
