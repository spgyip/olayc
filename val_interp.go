package olayc

import (
	"strconv"
	"strings"
)

// Interpret string value to concrete type, return nil if fail.
// foo, "123" => string
// 123, -123 => int
// 123.0 => float64
// true, false => bool
// 1,2,3 => array int
func interpret(s string) any {
	if isInt(s) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil
		}
		return i
	} else if isFloat(s) {
		i, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil
		}
		return i

	} else if isBool(s) {
		sl := strings.ToLower(s)
		if sl == "false" {
			return false
		}
		return true
	}
	return s
}

func isInt(s string) bool {
	for _, c := range s {
		if !(c >= '0' && c <= '9') {
			return false
		}
	}
	return true
}

func isFloat(s string) bool {
	countDot := 0
	for _, c := range s {
		if c == '.' {
			countDot++
			continue
		}
		if !(c >= '0' && c <= '9') {
			return false
		}
	}
	if countDot > 1 {
		return false
	}
	return true
}

func isBool(s string) bool {
	sl := strings.ToLower(s)
	if sl == "true" || sl == "false" {
		return true
	}
	return false
}
