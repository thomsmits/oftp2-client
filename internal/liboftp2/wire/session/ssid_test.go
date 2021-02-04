package session

import (
	"reflect"
	"testing"
)

func TestSSID_RoundTrip(t *testing.T) {
	a1 := SSID{
		Id:             "O1818181DDD",
		Password:       "PASSWORD",
		BufferSize:     10,
		Capability:     "X",
		Compress:       true,
		Restart:        true,
		Special:        true,
		Credit:         10,
		Authentication: true,
		UserData:       "USERDATA",
	}
	b := a1.Marshal()

	a2 := SSID{}
	err := a2.Parse(b)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(a1, a2) {
		t.Errorf("Roundtrip failed: %v != %v", a1, a2)
	}
}
