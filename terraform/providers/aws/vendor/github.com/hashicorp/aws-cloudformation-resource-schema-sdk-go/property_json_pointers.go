package cfschema

// PropertyJsonPointers is a list of PropertyJsonPointer.
type PropertyJsonPointers []PropertyJsonPointer

// ContainsPath returns true if an element matches the path.
func (ptrs PropertyJsonPointers) ContainsPath(path []string) bool {
	for _, ptr := range ptrs {
		if ptr.EqualsPath(path) {
			return true
		}
	}

	return false
}
