package session

import (
	"reflect"
	"testing"
)

func TestSSRM_RoundTrip(t *testing.T) {
	a1 := SSRM{}

	b := a1.Marshal()

	a2 := SSRM{}
	err := a2.Parse(b)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(a1, a2) {
		t.Errorf("Roundtrip failed: %v != %v", a1, a2)
	}
}
