package olayc

import (
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// Value represents a configure value, it can be scalar node or sub-tree node.
type Value any

// Unmarshal is implemented by using yaml utility,
// value is firstly marshalled to yaml bytes,
// then the yaml bytes is unmarshal to target out.
// Thus, if 'out' is a struct, you must use the yaml struct tag.
func Unmarshal(v Value, out any) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return errors.Wrap(err, "Value.Unmarshal fail")
	}
	return yaml.Unmarshal(data, out)
}
