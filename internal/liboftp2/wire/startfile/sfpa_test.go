package startfile

import (
	"reflect"
	"testing"
)

func TestSFPA_RoundTrip(t *testing.T) {
	a1 := SFPA{
		AnswerCount: 2,
	}

	b := a1.Marshal()

	a2 := SFPA{}
	err := a2.Parse(b)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(a1, a2) {
		t.Errorf("Roundtrip failed: %v != %v", a1, a2)
	}
}
