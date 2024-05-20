package foundation

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

const (
	absolutePath = "foundation"
	userAgent    = "foundation"
	clientName   = "foundation"
)

//Foundation client with its services
type Client struct {

	//base client
	client *client.Client

	//Service for Imaging Nodes and Cluster Creation
	NodeImaging NodeImagingService

	//Service for File Management in foundation VM
	FileManagement FileManagementService

	//Service for Networking apis in foundation VM
	Networking NetworkingService
}

//This routine returns new Foundation API Client
func NewFoundationAPIClient(credentials client.Credentials) (*Client, error) {
	var baseClient *client.Client
	if credentials.FoundationEndpoint != "" {
		// for foundation client, url should be based on foundation's endpoint and port
		credentials.URL = fmt.Sprintf("%s:%s", credentials.FoundationEndpoint, credentials.FoundationPort)
		c, err := client.NewBaseClient(&credentials, absolutePath, true)
		if err != nil {
			return nil, err
		}
		baseClient = c
	} else {
		errorMsg := fmt.Sprintf("Foundation Client is missing. "+
			"Please provide required detail - %s in provider configuration.", strings.Join(credentials.RequiredFields[clientName], ", "))
		// create empty client if required fields are not provided
		baseClient = &client.Client{ErrorMsg: errorMsg}
	}

	//Fill user agent details
	baseClient.UserAgent = userAgent

	foundationClient := &Client{
		client: baseClient,
		NodeImaging: NodeImagingOperations{
			client: baseClient,
		},
		FileManagement: FileManagementOperations{
			client: baseClient,
		},
		Networking: NetworkingOperations{
			client: baseClient,
		},
	}
	return foundationClient, nil
}
