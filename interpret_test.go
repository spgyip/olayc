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
		{"123", reflect.Int},
		{"123.0", reflect.Float64},
		{"\"123.0\"", reflect.String},
		{"true", reflect.Bool},
		{"false", reflect.Bool},
		{"\"true\"", reflect.String},
		{"\"false\"", reflect.String},
		//{"1,2,3", reflect.Slice},
	} {
		valI := interpret(test.val)
		if valI == nil {
			t.Errorf("interpret got nil, value %v\n", test.val)
			continue
		}
		got := reflect.TypeOf(valI).Kind()
		if got != test.expect {
			t.Errorf("got(%v)!=expect(%v), value %v\n", got, test.expect, test.val)
		}
	}
}
