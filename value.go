package olayc

import (
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// Value represents a configure value, it can be scalar node or sub-tree node.
type Value any

// Unmarshal value to out, value is marshalled to yaml then unmarshalled to out.
func UnmarshalValue(v Value, out any) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return errors.Wrap(err, "Value.Unmarshal fail")
	}
	return yaml.Unmarshal(data, out)
}
