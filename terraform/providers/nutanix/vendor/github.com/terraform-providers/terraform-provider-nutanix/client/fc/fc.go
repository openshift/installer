package foundationcentral

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

const (
	libraryVersion = "v1"
	absolutePath   = "api/fc/" + libraryVersion
	userAgent      = "nutanix/" + libraryVersion
	clientName     = "foundation_central"
)

// Client manages the foundation central API
type Client struct {
	client  *client.Client
	Service Service
}

// NewFoundationCentralClient return a client to operate foundation central resources
func NewFoundationCentralClient(credentials client.Credentials) (*Client, error) {
	var baseClient *client.Client

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		c, err := client.NewClient(&credentials, userAgent, absolutePath, false)
		if err != nil {
			return nil, err
		}
		baseClient = c
	} else {
		errorMsg := fmt.Sprintf("Foundation Central Client is missing. "+
			"Please provide required details - %s in provider configuration.", strings.Join(credentials.RequiredFields[clientName], ", "))

		baseClient = &client.Client{UserAgent: userAgent, ErrorMsg: errorMsg}
	}

	fc := &Client{
		client: baseClient,
		Service: Operations{
			client: baseClient,
		},
	}
	return fc, nil
}
