package olayc

import (
	"testing"
)

func TestFlags(t *testing.T) {
	fp := &flagParser{}
	fp.parse([]string{
		"-name1=123",
		"-name2", "123",
		"--name3=123",
		"--name4", "123",
		"-id1", "123.0",
		"-on1",         // Bool value, default true
		"-on2", "true", // Bool value
		"-on3", "false", // Bool value
		"invalid-arg",
	})

	var got = fp.kvs
	var expect = []kv{
		{"name1", 123},
		{"name2", 123},
		{"name3", 123},
		{"name4", 123},
		{"id1", 123.0},
		{"on1", true},
		{"on2", true},
		{"on3", false},
	}

	if len(expect) != len(got) {
		t.Logf("expect %v\n", expect)
		t.Logf("got %v\n", got)
		t.Fatal("expect != got")
	}

	for i, kvExpect := range expect {
		kvGot := got[i]
		if kvExpect.key != kvGot.key || kvExpect.value != kvGot.value {
			t.Errorf("expect[%v](%v) != got[%v](%v)", i, expect[i], i, got[i])
		}
	}
}
