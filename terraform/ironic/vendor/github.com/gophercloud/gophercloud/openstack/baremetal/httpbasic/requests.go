package httpbasic

import (
	"encoding/base64"
	"fmt"

	"github.com/gophercloud/gophercloud"
)

// EndpointOpts specifies a "http_basic" Ironic Endpoint
type EndpointOpts struct {
	IronicEndpoint     string
	IronicUser         string
	IronicUserPassword string
}

func initClientOpts(client *gophercloud.ProviderClient, eo EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc := new(gophercloud.ServiceClient)
	if eo.IronicEndpoint == "" {
		return nil, fmt.Errorf("IronicEndpoint is required")
	}
	if eo.IronicUser == "" || eo.IronicUserPassword == "" {
		return nil, fmt.Errorf("User and Password are required")
	}

	token := []byte(eo.IronicUser + ":" + eo.IronicUserPassword)
	encodedToken := base64.StdEncoding.EncodeToString(token)
	sc.MoreHeaders = map[string]string{"Authorization": "Basic " + encodedToken}
	sc.Endpoint = gophercloud.NormalizeURL(eo.IronicEndpoint)
	sc.ProviderClient = client
	return sc, nil
}

// NewBareMetalHTTPBasic creates a ServiceClient that may be used to access a
// "http_basic" bare metal service.
func NewBareMetalHTTPBasic(eo EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(&gophercloud.ProviderClient{}, eo)
	if err != nil {
		return nil, err
	}

	sc.Type = "baremetal"

	return sc, nil
}
