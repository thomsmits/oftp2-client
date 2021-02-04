package authentication

import (
	"testing"
)

func TestSECD_RoundTrip(t *testing.T) {
	a1 := SECD{}
	b := a1.Marshal()

	a2 := SECD{}
	err := a2.Parse(b)
	if err != nil {
		t.Error(err)
	}
}
