package nutanix

import (
	"fmt"

	nutanixclient "github.com/nutanix-cloud-native/prism-go-client"
	nutanixclientv3 "github.com/nutanix-cloud-native/prism-go-client/v3"
)

// CreateNutanixClient creates a Nutanix V3 Client
func CreateNutanixClient(prismCentral, port, username, password string) (*nutanixclientv3.Client, error) {
	cred := nutanixclient.Credentials{
		URL:      fmt.Sprintf("%s:%s", prismCentral, port),
		Username: username,
		Password: password,
		Port:     port,
		Endpoint: prismCentral,
	}

	return nutanixclientv3.NewV3Client(cred)
}
