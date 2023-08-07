package olayc

import (
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

// Interpret type of s.
func typeInterpret(s string) reflect.Kind {
	var isAllDigit = true
	var cntDot = 0
	for _, c := range s {
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

	// If it's int or float64
	if isAllDigit && cntDot == 1 {
		return reflect.Float64
	} else if isAllDigit {
		return reflect.Int
	}

	sl := strings.ToLower(s)
	if sl == "true" || sl == "false" {
		return reflect.Bool
	}
	return reflect.String
}

// Interpret string value to concrete type, return nil if fail.
// foo, "123", "true", "false" => string
// 123, -123                   => int
// 123.0                       => float64
// true, false                 => bool
func interpret(s string) any {
	var v any
	var err error

	kind := typeInterpret(s)
	switch kind {
	case reflect.Int:
		v, err = strconv.Atoi(s)
	case reflect.Float64:
		v, err = strconv.ParseFloat(s, 32)
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
