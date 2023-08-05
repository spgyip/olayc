package olayc

import (
	"strings"

	"github.com/pkg/errors"
)

type kv struct {
	key   string
	value any
}

type flagParser struct {
	args []string
	kvs  []kv
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
	var value any
	f.args = f.args[1:]

	if len(key) >= 2 && key[0:2] == "--" {
		key = key[2:]
	} else if len(key) >= 1 && key[0:1] == "-" {
		key = key[1:]
	} else {
		return false, errors.Errorf("argument invalid: %v", key)
	}

	sps := strings.Split(key, "=")
	if len(sps) == 1 {
		if len(f.args) == 0 || f.args[0][0] == '-' {
			// This is the last argument,
			// or the next argument is not a value with prefix "-",
			// The value is parsed as true value of bool type.
			// E.g. '-on'
			value = "true"
		} else if len(f.args) > 0 {
			// E.g. '-name=foo'
			// Get value from the next argument.
			value = f.args[0]
			f.args = f.args[1:]
		}
	} else {
		// Format of '-name=foo'
		key = sps[0]
		value = sps[1]
	}

	// Interpret string value "true", "false" to bool type.
	// TODO: interpret type int/uint/float64/array
	if strVal, ok := value.(string); ok {
		if strings.ToLower(strVal) == "true" {
			value = true
		} else if strings.ToLower(strVal) == "false" {
			value = false
		}
	}
	f.kvs = append(f.kvs, kv{key, value})

	return false, nil
}
