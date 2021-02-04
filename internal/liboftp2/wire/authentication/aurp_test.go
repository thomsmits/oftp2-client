package authentication

import (
	"bytes"
	"testing"
)

func TestAURP_RoundTrip(t *testing.T) {
	a1 := AURP{
		Response: []byte("A very secret respon"),
	}
	b := a1.Marshal()

	a2 := AURP{}
	err := a2.Parse(b)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(a1.Response, a2.Response) {
		t.Errorf("Roundtrip failed: %v != %v", a1, a2)
	}
}
