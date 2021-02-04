package authentication

import (
	"bytes"
	"testing"
)

func TestAUCH_RoundTrip(t *testing.T) {
	a1 := AUCH{
		Challenge: []byte("A very secret challenge"),
	}
	b := a1.Marshal()

	a2 := AUCH{}
	err := a2.Parse(b)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(a1.Challenge, a2.Challenge) {
		t.Errorf("Roundtrip failed: %v != %v", a1, a2)
	}
}
