package nutanix

// Metadata contains Nutanix metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	// PrismCentral is the domain name or IP address of the Prism Central.
	PrismCentral string `json:"prismCentral"`
	// Username is the name of the user to use to connect to the Prism Central.
	Username string `json:"username"`
	// Password is the password for the user to use to connect to the Prism Central.
	Password string `json:"password"`
	// Port is the port used to connect to the Prism Central.
	Port string `json:"port"`
}
