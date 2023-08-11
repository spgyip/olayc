package olayc

import (
	"fmt"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// Value represents a configure value, it can be scalar node or sub-tree node.
type Value struct {
	v any
}

// Return if it's nil value
func (v *Value) IsNil() bool {
	return v.v == nil
}

// Get string value, return "" if it doesn't exist.
func (v *Value) String() string {
	if v.v == nil {
		return ""
	}

	var s = ""
	switch x := v.v.(type) {
	case string:
		s = x
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64,
		bool:
		s = fmt.Sprintf("%v", x)
	}
	return s
}

// Get int value, return 0 if fails.
func (v *Value) Int() int {
	if v.v == nil {
		return 0
	}

	var i = 0
	switch x := v.v.(type) {
	case int:
		i = int(x)
	case int8:
		i = int(x)
	case int16:
		i = int(x)
	case int32:
		i = int(x)
	case int64:
		i = int(x)
	case uint:
		i = int(x)
	case uint8:
		i = int(x)
	case uint16:
		i = int(x)
	case uint32:
		i = int(x)
	case uint64:
		i = int(x)
	}
	return i
}

// Get uint value, return 0 if fails.
func (v *Value) Uint() uint {
	if v.v == nil {
		return 0
	}

	var i uint = 0
	switch x := v.v.(type) {
	case int:
		i = uint(x)
	case int8:
		i = uint(x)
	case int16:
		i = uint(x)
	case int32:
		i = uint(x)
	case int64:
		i = uint(x)
	case uint:
		i = uint(x)
	case uint8:
		i = uint(x)
	case uint16:
		i = uint(x)
	case uint32:
		i = uint(x)
	case uint64:
		i = uint(x)
	}
	return i
}

// Get int64 value, return 0 if failes.
func (v *Value) Int64() int64 {
	if v.v == nil {
		return 0
	}

	var i int64 = 0
	switch x := v.v.(type) {
	case int:
		i = int64(x)
	case int8:
		i = int64(x)
	case int16:
		i = int64(x)
	case int32:
		i = int64(x)
	case int64:
		i = int64(x)
	case uint:
		i = int64(x)
	case uint8:
		i = int64(x)
	case uint16:
		i = int64(x)
	case uint32:
		i = int64(x)
	case uint64:
		i = int64(x)
	}
	return i
}

// Get uint64 value, return 0 if fails.
func (v *Value) Uint64() uint64 {
	if v.v == nil {
		return 0
	}

	var i uint64 = 0
	switch x := v.v.(type) {
	case int:
		i = uint64(x)
	case int8:
		i = uint64(x)
	case int16:
		i = uint64(x)
	case int32:
		i = uint64(x)
	case int64:
		i = uint64(x)
	case uint:
		i = uint64(x)
	case uint8:
		i = uint64(x)
	case uint16:
		i = uint64(x)
	case uint32:
		i = uint64(x)
	case uint64:
		i = uint64(x)
	}
	return i
}

// Get float64 value, return 0.0 if fails.
func (v *Value) Float64() float64 {
	if v.v == nil {
		return 0.0
	}

	var i = 0.0
	switch x := v.v.(type) {
	case float32:
		i = float64(x)
	case float64:
		i = float64(x)
	}
	return i
}

// Get bool value, return false if fails.
func (v *Value) Bool() bool {
	if v.v == nil {
		return false
	}

	var i = false
	switch x := v.v.(type) {
	case bool:
		i = bool(x)
	}
	return i
}

// Unmarshal is implemented by using yaml utility,
// value is firstly marshalled to yaml bytes,
// then the yaml bytes is unmarshal to target out.
// Thus, if 'out' is a struct, you must use the yaml struct tag.
func (v *Value) Unmarshal(out any) error {
	data, err := yaml.Marshal(v.v)
	if err != nil {
		return errors.Wrap(err, "Value.Unmarshal fail")
	}
	return yaml.Unmarshal(data, out)
}
