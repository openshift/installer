package op

// Convenience functions for returning pointers to values

// Int returns a pointer to the int value provided
func Int(v int) *int {
	return &v
}

// String returns a pointer to the string value provided
func String(v string) *string {
	return &v
}

// Bool returns a pointer to the bool value provided
func Bool(v bool) *bool {
	return &v
}
