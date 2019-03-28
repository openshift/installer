package azure

//DNSConfig is an interface that provides means to fetch the DNS settings
type DNSConfig struct {
}

//GetBaseDomain returns the base domain to use
func (config DNSConfig) GetBaseDomain() (string, error) {
	return "", nil
}

//GetPublicZone returns the public zone id to create subdomain during deployment
func (config DNSConfig) GetPublicZone(name string) string { //returns ID
	return ""
}

//NewDNSConfig returns a new DNSConfig struct that helps configuring the DNS
//by querying your subscription and letting you choose
//which domain you wish to use for the cluster
func NewDNSConfig() (*DNSConfig, error) {
	//get session here, set session on config object.
	//each DNSConfig method implements the querying.
	//allows for fake session injection
	return &DNSConfig{}, nil
}
