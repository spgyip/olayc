package olayc

import (
	"reflect"
	"testing"
)

func TestInterpret(t *testing.T) {
	for _, test := range []struct {
		val    string
		expect reflect.Kind
	}{
		{"foo", reflect.String},
		{"\"123\"", reflect.String},
		{"\"123.0\"", reflect.String},
		{"\"true\"", reflect.String},
		{"\"false\"", reflect.String},

		{"123", reflect.Int},

		{"123.0", reflect.Float64},

		{"true", reflect.Bool},
		{"True", reflect.Bool},
		{"false", reflect.Bool},
		{"False", reflect.Bool},
		//{"1,2,3", reflect.Slice},
	} {
		gotVal := interpret(test.val)
		if gotVal == nil {
			t.Errorf("interpret fail: %v\n", test.val)
			continue
		}
		got := reflect.TypeOf(gotVal).Kind()
		if got != test.expect {
			t.Errorf("got(%v)!=expect(%v), value %v\n", got, test.expect, test.val)
		}
	}
}
