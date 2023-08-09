package olayc

import (
	"strings"

	"github.com/pkg/errors"
)

// KV is composition of key and value.
type KV struct {
	key   string
	value any
}

// flagParser parses from args to kvs.
type flagParser struct {
	args []string
	kvs  []KV
}

// Parse arguments. Invalid arguments will be ignored.
// -name=foo | --name=foo            // parsed as <name, foo>
// -name foo | --name foo            // parsed as <name, foo>
// -on       | --on                  // parsed as <on, true>
// -on false | --on false            // parsed as <on, false>
// -on=false | --on=false            // parsed as <on, false>
//
// Value interpretation:
// "true" => true   # Not case sensitive
// "false" => false # Not case sensitive
func (f *flagParser) parse(args []string) int {
	f.args = args
	end := false
	for !end {
		end, _ = f.parseOne()
	}
	return len(f.kvs)
}

// Return true if encounter end of the last argument.
// Return error if argument invalid.
func (f *flagParser) parseOne() (bool, error) {
	if len(f.args) == 0 {
		return true, nil
	}
	key := f.args[0]
	f.args = f.args[1:]

	var value any
	var strValue string

	if len(key) >= 2 && key[0:2] == "--" {
		key = key[2:]
	} else if len(key) >= 1 && key[0:1] == "-" {
		key = key[1:]
	} else {
		return false, errors.Errorf("Invalid argument: %v", key)
	}

	// Find '=' in key, there are cases:
	// "-name=foo" => <name, foo>
	// "-name foo" => <name, foo>
	// "-onoff"     => <onoff, true>
	// "-onof -name=foo" => <onoff, true>, <name, foo>
	pos := strings.IndexByte(key, '=')
	if pos < 0 {
		if len(f.args) == 0 || (f.args[0][0] == '-' && !isAllDigit(f.args[0][1:])) {
			// E.g. '-onoff'
			strValue = "true"
		} else if len(f.args) > 0 {
			// E.g. '-name foo'
			strValue = f.args[0]
			f.args = f.args[1:]
		}
	} else {
		// E.g. '-name=foo'
		strValue = key[pos+1:]
		key = key[:pos]
	}

	// Interpret string to concrete type value
	value = interpret(strValue)
	f.kvs = append(f.kvs, KV{key, value})

	return false, nil
}
