package nutanix

import (
	"context"
	"fmt"
	"time"

	nutanixclient "github.com/nutanix-cloud-native/prism-go-client/pkg/nutanix"
	nutanixclientv3 "github.com/nutanix-cloud-native/prism-go-client/pkg/nutanix/v3"
)

// CreateNutanixClient creates a Nutanix V3 Client
func CreateNutanixClient(ctx context.Context, prismCentral, port, username, password string) (*nutanixclientv3.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	cred := nutanixclient.Credentials{
		URL:      fmt.Sprintf("%s:%s", prismCentral, port),
		Username: username,
		Password: password,
		Port:     port,
		Endpoint: prismCentral,
	}

	return nutanixclientv3.NewV3Client(cred)
}
