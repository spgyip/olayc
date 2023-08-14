package olayc

import (
	"strings"

	"github.com/pkg/errors"
)

// flagParser parses from args to kvs.
type flagParser struct {
	args []string
	kvs  []KV
}

// Parse arguments. Invalid arguments will be ignored.
// There are several forms:
// -key=value | --key=value
// -key value | --key value
// -key       | --key         # Same as -key=true, -key true
//
// Value interpretation should refer to `interpreted(string)`.
func (psr *flagParser) parse(args []string) int {
	psr.args = args
	end := false
	for !end {
		end, _ = psr.parseOne()
	}
	return len(psr.kvs)
}

// Return true if encounter end of the last argument.
// Return error if argument invalid.
func (psr *flagParser) parseOne() (bool, error) {
	if len(psr.args) == 0 {
		return true, nil
	}
	key := psr.args[0]
	psr.args = psr.args[1:]

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
		if len(psr.args) == 0 || (psr.args[0][0] == '-' && !isAllDigit(psr.args[0][1:])) {
			// E.g. '-onoff'
			strValue = "true"
		} else if len(psr.args) > 0 {
			// E.g. '-name foo'
			strValue = psr.args[0]
			psr.args = psr.args[1:]
		}
	} else {
		// E.g. '-name=foo'
		strValue = key[pos+1:]
		key = key[:pos]
	}

	// Interpret string to concrete type value
	value = interpret(strValue)
	psr.kvs = append(psr.kvs, KV{key, value})

	return false, nil
}
