package networkextensions

// Error is the error type for this package, ready to be unwrapped with
// `errors.As`.
type Error string

// Error implements the error interface
func (err Error) Error() string {
	return string(err)
}
