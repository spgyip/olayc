package olayc

// Value represents a configure value, it can be scalar node or sub-tree node.
type Value any

// Get string from value, if value is nil or not type of string, return defaultValue.
/*func (v *Value) String(defaultValue string) string {
	s, ok := v.(string)
	if !ok {
		return defaultValue
	}
	return s
}*/

// ConfigSource is the interface defines configure source behaviors.
type ConfigSource interface {
	// Get value with the given key, return nil if doesn't exist.
	// The key is splitted by seperator '.', e.g. 'foo.name'.
	// The key is case sensitive, thus, 'foo.Name' is different from 'foo.name'.
	Get(string) Value
}
