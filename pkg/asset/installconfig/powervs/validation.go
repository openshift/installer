package powervs

// ValidateForProvisioning only validates credentials
// @TODO: Expand this to use the install config creds
func ValidateForProvisioning() error {
	_, err := GetSession()
	return err
}
