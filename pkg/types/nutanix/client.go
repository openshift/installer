package nutanix

import (
	"context"
	"fmt"
	"strconv"
	"time"

	nutanixclient "github.com/nutanix-cloud-native/prism-go-client"
	nutanixclientv3 "github.com/nutanix-cloud-native/prism-go-client/v3"
)

// CreateNutanixClient creates a Nutanix V3 Client.
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

// CreateNutanixClientFromPlatform creates a Nutanix V3 clinet based on the platform configuration.
func CreateNutanixClientFromPlatform(platform *Platform) (*nutanixclientv3.Client, error) {
	return CreateNutanixClient(context.TODO(),
		platform.PrismCentral.Endpoint.Address,
		strconv.Itoa(int(platform.PrismCentral.Endpoint.Port)),
		platform.PrismCentral.Username,
		platform.PrismCentral.Password)
}
