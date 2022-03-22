package karbon

import (
	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

const (
	absolutePath = "karbon"
	userAgent    = "nutanix"
)

// Client manages the V3 API
type Client struct {
	client          *client.Client
	Cluster         ClusterService
	PrivateRegistry PrivateRegistryService
	Meta            MetaService
}

// NewKarbonAPIClient return a client to operate Karbon resources
func NewKarbonAPIClient(credentials client.Credentials) (*Client, error) {
	c, err := client.NewClient(&credentials, userAgent, absolutePath)

	if err != nil {
		return nil, err
	}

	f := &Client{
		client: c,
		Cluster: ClusterOperations{
			client: c,
		},
		PrivateRegistry: PrivateRegistryOperations{
			client: c,
		},
		Meta: MetaOperations{
			client: c,
		},
	}

	return f, nil
}
