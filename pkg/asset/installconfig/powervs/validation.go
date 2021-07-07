package powervs

// Just see if we can create an IBMPISession
// @TODO: Expand this to use the install config creds
func ValidateForProvisioning() error {
	_, err := GetSession()
	return err
}
