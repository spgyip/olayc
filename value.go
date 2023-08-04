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
