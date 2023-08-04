package olayc

import (
	"testing"
)

func TestFlags(t *testing.T) {
	fl := &flags{}
	fl.parse([]string{
		"--id=123",
		"-name=foo1",
		"-redis.host",
		"redis.cluster",
		"invalid",
		"-switch",
		"--redis.port",
		"8306",
	})

	var got = fl.kvs

	var expect = []kv{
		{"id", "123"},
		{"name", "foo1"},
		{"redis.host", "redis.cluster"},
		{"switch", true},
		{"redis.port", "8306"},
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
