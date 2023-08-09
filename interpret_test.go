package olayc

import (
	"reflect"
	"testing"
)

func TestTypeInterpret(t *testing.T) {
	for i, test := range []struct {
		s      string
		expect reflect.Kind
	}{
		{"hellowrold", reflect.String},
		{"123", reflect.Uint64},
		{"-123", reflect.Int64},
		{"123.0", reflect.Float64},
		{"true", reflect.Bool},
		{"True", reflect.Bool},
		{"false", reflect.Bool},
		{"False", reflect.Bool},
	} {
		got := typeInterpret(test.s)
		if got != test.expect {
			t.Errorf("[%v] got(%v) != expect(%v)\n", i, got, test.expect)
		}
	}
}

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
