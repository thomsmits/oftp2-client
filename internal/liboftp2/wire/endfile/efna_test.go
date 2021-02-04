package endfile

import (
	"reflect"
	"testing"
)

func TestEFNA_RoundTrip(t *testing.T) {
	a1 := EFNA{
		ReasonCode: 9,
		AnswerText: "A very long and reasonable text explaining what happened",
	}
	b := a1.Marshal()

	a2 := EFNA{}
	err := a2.Parse(b)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(a1, a2) {
		t.Errorf("Roundtrip failed: %v != %v", a1, a2)
	}
}
