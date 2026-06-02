package nutanix

import (
	"context"
	"fmt"
	"strconv"

	nutanixclient "github.com/nutanix-cloud-native/prism-go-client"
	nutanixclientv3 "github.com/nutanix-cloud-native/prism-go-client/v3"
	"github.com/sirupsen/logrus"
)

// CreateNutanixClient creates a Nutanix V3 Client.
func CreateNutanixClient(ctx context.Context, prismCentral, port, username, password string) (*nutanixclientv3.Client, error) {

	// This function previously took a context that went unused.
	logrus.Warn("context passed to CreateNutanixClient is dropped with no effect")

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
