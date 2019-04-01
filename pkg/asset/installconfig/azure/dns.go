package azure

import (
	"github.com/pkg/errors"
)

//DNSConfig implements dns.ConfigProvider interface which provides methods to choose the DNS settings
type DNSConfig struct {
	Session *Session
}

//GetBaseDomain returns the base domain to use
func (config DNSConfig) GetBaseDomain() (string, error) {
	//call azure api using the session to retrieve available base domain
	return "eastus.cloudapp.azure.com", nil
}

//GetPublicZone returns the public zone id to create subdomain during deployment
func (config DNSConfig) GetPublicZone(name string) (string, error) { //returns ID
	//call azure api using the session to return reference to available public zone
	return "", nil
}

//NewDNSConfig returns a new DNSConfig struct that helps configuring the DNS
//by querying your subscription and letting you choose
//which domain you wish to use for the cluster
func NewDNSConfig() (*DNSConfig, error) {
	session, err := GetSession()
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve session information")
	}
	//get session here, set session on config object.
	//each DNSConfig method implements the querying.
	//allows for fake session injection and tests
	return &DNSConfig{Session: session}, nil
}
