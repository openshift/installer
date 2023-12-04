package dns

// UserProvisionedDNS indicates whether the DNS solution is provisioned by the Installer or the user.
type UserProvisionedDNS string

const (
	// UserProvisionedDNSEnabled indicates that the DNS solution is provisioned and provided by the user.
	UserProvisionedDNSEnabled UserProvisionedDNS = "Enabled"

	// UserProvisionedDNSDisabled indicates that the DNS solution is provisioned by the Installer.
	UserProvisionedDNSDisabled UserProvisionedDNS = "Disabled"
)

// ValidUserProvisionedDNSType verifies the userProvisionedDNS string is valid.
func ValidUserProvisionedDNSType(userProvisionedDNS UserProvisionedDNS) bool {
	return userProvisionedDNS == UserProvisionedDNSEnabled || userProvisionedDNS == UserProvisionedDNSDisabled
}
