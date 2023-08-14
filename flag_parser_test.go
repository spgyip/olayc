package olayc

import (
	"reflect"
	"testing"
)

func TestFlagParser(t *testing.T) {
	psr := &flagParser{}
	psr.parse([]string{
		"-foo.id1=123",
		"-foo.id2", "123",
		"--foo.id3=123",
		"--foo.id4", "123",
		"-foo.id5", "-123",
		"-foo.temp", "123.0",
		"-foo.on1",         // Bool value, default true
		"-foo.on2", "true", // Bool value
		"-foo.on3", "false", // Bool value
		"invalid-arg",
	})

	var got = psr.kvs
	var expect = []KV{
		{"foo.id1", uint64(123)},
		{"foo.id2", uint64(123)},
		{"foo.id3", uint64(123)},
		{"foo.id4", uint64(123)},
		{"foo.id5", int64(-123)},
		{"foo.temp", float64(123.0)},
		{"foo.on1", true},
		{"foo.on2", true},
		{"foo.on3", false},
	}

	if len(expect) != len(got) {
		t.Logf("expect(len %v) %v\n", len(expect), expect)
		t.Logf("got(len %v) %v\n", len(got), got)
		t.Fatal("expect != got")
	}

	for i, kvExpect := range expect {
		kvGot := got[i]
		if kvExpect.key != kvGot.key || kvExpect.value != kvGot.value {
			t.Errorf("[%v] expect(%v:%v) != got(%v:%v)", i, expect[i], reflect.TypeOf(expect[i].value),
				got[i], reflect.TypeOf(got[i].value))
		}
	}
}
