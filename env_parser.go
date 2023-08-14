package olayc

import (
	"strings"
)

// envParser parses from ENVs to kvs.
type envParser struct {
	kvs []KV
}

// Parse environments to kvs. The env must be in the form "key=value".
// The key is converted to lower case and the seperator '_' is replaced by '.'.
// E.g. 'LC_CTYPE=UTF-8', is converted to 'lc.ctype=UTF-8'.
// The anterior '_' in key will be trimed, e.g. '_P9K_SSH_TTY' is converted to `p9k.ssh.tty`.
//
// Value interpretation should refer to `func interpreted(string)`.
func (psr *envParser) parse(envs []string) int {
	// Replace '_' to '.', trim the anterior '_'.
	replaceFunc := func(s string) string {
		sl := []byte(s)
		var i int
		for i = 0; i < len(sl); i++ {
			if sl[i] != '_' {
				break
			}
		}
		sl = sl[i:] // Trim anterior '_'
		for i = 0; i < len(sl); i++ {
			if sl[i] == '_' {
				sl[i] = '.'
			}
		}
		return string(sl)
	}

	for _, e := range envs {
		sps := strings.SplitN(e, "=", 2)
		if len(sps) != 2 {
			continue
		}

		var key = strings.ToLower(replaceFunc(sps[0]))
		var value any = interpret(sps[1])
		if len(key) > 0 {
			psr.kvs = append(psr.kvs, KV{key, value})
		}
	}
	return len(psr.kvs)
}
