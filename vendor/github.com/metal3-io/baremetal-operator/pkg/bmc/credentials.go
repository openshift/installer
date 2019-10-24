package bmc

// Credentials holds the information for authenticating with the BMC.
type Credentials struct {
	Username string
	Password string
}

// Validate returns an error if the credentials are invalid
func (creds Credentials) Validate() error {
	if creds.Username == "" {
		return &CredentialsValidationError{message: "Missing BMC connection detail 'username' in credentials"}
	}
	if creds.Password == "" {
		return &CredentialsValidationError{message: "Missing BMC connection details 'password' in credentials"}
	}
	return nil
}
