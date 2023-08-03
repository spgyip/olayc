package olayc

import (
	"strings"

	"github.com/pkg/errors"
)

type kv struct {
	key   string
	value string
}

type flags struct {
	args []string
	kvs  []kv
}

// newFlags allocates and returns a flags.
func newFlags() *flags {
	return &flags{}
}

// Parse arguments.
func (f *flags) parse(args []string) {
	f.args = args
	end := false
	for !end {
		end, _ = f.parseOne()
	}
}

// Return true if end of last arguments.
// Return error if argument invalid.
func (f *flags) parseOne() (bool, error) {
	if len(f.args) == 0 {
		return true, nil
	}

	key, value := f.args[0], ""
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
		if len(f.args) == 0 {
			return false, errors.Errorf("No value specified for key: %v", key)
		}

		value = f.args[0]
		f.args = f.args[1:]
	} else {
		key = sps[0]
		value = sps[1]
	}
	f.kvs = append(f.kvs, kv{key, value})

	return false, nil
}
