package startfile

import (
	"reflect"
	"testing"
)

func TestSFNA_RoundTrip(t *testing.T) {
	a1 := SFNA{
		ReasonCode:     12,
		RetryIndicator: true,
		ReasonText:     "Some additiona explanation",
	}

	b := a1.Marshal()

	a2 := SFNA{}
	err := a2.Parse(b)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(a1, a2) {
		t.Errorf("Roundtrip failed: %v != %v", a1, a2)
	}
}
