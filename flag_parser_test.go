package olayc

import (
	"reflect"
	"testing"
)

func TestFlags(t *testing.T) {
	fp := &flagParser{}
	fp.parse([]string{
		"-id1=123",
		"-id2", "123",
		"--id3=123",
		"--id4", "123",
		"-id5", "-123",
		"-temp", "123.0",
		"-on1",         // Bool value, default true
		"-on2", "true", // Bool value
		"-on3", "false", // Bool value
		"invalid-arg",
	})

	var got = fp.kvs
	var expect = []kv{
		{"id1", uint64(123)},
		{"id2", uint64(123)},
		{"id3", uint64(123)},
		{"id4", uint64(123)},
		{"id5", int64(-123)},
		{"temp", float64(123.0)},
		{"on1", true},
		{"on2", true},
		{"on3", false},
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
