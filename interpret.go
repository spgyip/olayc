package olayc

import (
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

// Interpret type of string s.
func typeInterpret(s string) reflect.Kind {
	var isAllDigit = true
	var cntDot = 0
	var isSigned = false
	for i, c := range s {
		if i == 0 && c == '-' {
			isSigned = true
			continue
		}
		if !unicode.IsDigit(c) && c != '.' {
			isAllDigit = false
			break
		}
		if c == '.' {
			cntDot++
		}
		if cntDot > 1 {
			break
		}
	}

	// Default type is string
	var kind reflect.Kind = reflect.String
	if isAllDigit && cntDot <= 1 {
		if cntDot == 0 {
			if isSigned {
				kind = reflect.Int64
			} else {
				kind = reflect.Uint64
			}
		} else {
			kind = reflect.Float64
		}
	} else {
		// Bool
		sl := strings.ToLower(s)
		if sl == "true" || sl == "false" {
			kind = reflect.Bool
		}
	}
	return kind
}

// Interpret string to type value, return nil if fails.
// Int/Uint type is interpreted as 64-bits(int64/uint64), allows downcast to 32-bits type(int32/uint32) as necessary.
// foo, "123", "true", "false" => string
// 123                         => uint64
// -123                        => int64
// 123.0                       => float64
// true, false                 => bool
func interpret(s string) any {
	var v any
	var err error

	kind := typeInterpret(s)
	switch kind {
	case reflect.Int64:
		v, err = strconv.ParseInt(s, 10, 64)
	case reflect.Uint64:
		v, err = strconv.ParseUint(s, 10, 64)
	case reflect.Float64:
		v, err = strconv.ParseFloat(s, 64)
	case reflect.String:
		v, err = s, nil
	case reflect.Bool:
		sl := strings.ToLower(s)
		if sl == "true" {
			v = true
		} else {
			v = false
		}
	default:
		v, err = nil, errors.Errorf("unknown kind: %v", kind)
	}

	if err != nil {
		return nil
	}
	return v
}
