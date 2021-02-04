package wire

import (
	"testing"
	"time"
)

func TestTruncateString(t *testing.T) {
	s := "123456789012345678901234567890" // 30 Chars
	r := TruncateString(s, 20)
	e := "12345678901234567890"

	if r != e {
		t.Errorf("expected %s, got %s", e, r)
	}

	r = TruncateString(s, 31)
	e = s

	if r != e {
		t.Errorf("expected %s, got %s", e, r)
	}

	r = TruncateString(s, 30)
	e = s

	if r != e {
		t.Errorf("expected %s, got %s", e, r)
	}

	r = TruncateString(s, -1)
	e = s

	if r != e {
		t.Errorf("expected %s, got %s", e, r)
	}
}

func TestTruncateAndPadString(t *testing.T) {
	s := "123456789012345678901234567890" // 30 Chars
	r := TruncateAndPadString(s, 20)
	e := "12345678901234567890"

	if r != e {
		t.Errorf("expected %s, got %s", e, r)
	}

	r = TruncateAndPadString(s, 33)
	e = s + "   "

	if r != e {
		t.Errorf("expected %s, got %s", e, r)
	}

	r = TruncateAndPadString(s, 30)
	e = s

	if r != e {
		t.Errorf("expected %s, got %s", e, r)
	}

	r = TruncateAndPadString(s, -1)
	e = s

	if r != e {
		t.Errorf("expected %s, got %s", e, r)
	}
}

func TestParseDateToString(t *testing.T) {
	dateTime := time.Date(2020, 12, 17, 10, 22, 34, 345678912, time.UTC)
	dt, tm := ParseDateToString(dateTime)
	edt := "20201217"
	etm := "1022343456"

	if dt != edt {
		t.Errorf("expected date %s, got %s", edt, dt)
	}

	if tm != etm {
		t.Errorf("expected time %s, got %s", etm, tm)
	}
}

func TestParseStringsToDate(t *testing.T) {
	ds := "20201217"
	ts := "1022343456"

	dateTime := ParseStringsToDate(ds, ts)

	if dateTime.Year() != 2020 || dateTime.Month() != 12 || dateTime.Day() != 17 {
		t.Errorf("wrong date, got %v expected %s", dateTime, ds)
	}

	if dateTime.Hour() != 10 || dateTime.Minute() != 22 || dateTime.Second() != 34 {
		t.Errorf("wrong time, got %v expected %s", dateTime, ts)
	}
}
