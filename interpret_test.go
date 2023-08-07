package olayc

import (
	"testing"
)

func TestInterpret(t *testing.T) {
	for _, test := range []struct {
		val    string
		expect any
	}{
		{"foo", string("foo")},
		{"\"123\"", string("\"123\"")},
		{"\"123.0\"", string("\"123.0\"")},
		{"\"true\"", string("\"true\"")},
		{"\"false\"", string("\"false\"")},

		{"123", uint64(123)},
		{"-123", int64(-123)},

		{"123.0", float64(123.0)},
		{"-123.0", float64(-123.0)},

		{"true", bool(true)},
		{"True", bool(true)},
		{"false", bool(false)},
		{"False", bool(false)},
	} {
		got := interpret(test.val)
		if got != test.expect {
			t.Errorf("got(%v)!=expect(%v), value \"%v\"\n", got, test.expect, test.val)
		}
	}
}
