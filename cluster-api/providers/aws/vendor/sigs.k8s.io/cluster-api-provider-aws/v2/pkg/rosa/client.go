package rosa

import (
	"context"
	"fmt"
	"os"

	sdk "github.com/openshift-online/ocm-sdk-go"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

const (
	ocmTokenKey  = "ocmToken"
	ocmAPIURLKey = "ocmApiUrl"
)

type RosaClient struct {
	ocm       *sdk.Connection
	rosaScope *scope.ROSAControlPlaneScope
}

// NewRosaClientWithConnection creates a client with a preexisting connection for testing purposes.
func NewRosaClientWithConnection(connection *sdk.Connection, rosaScope *scope.ROSAControlPlaneScope) *RosaClient {
	return &RosaClient{
		ocm:       connection,
		rosaScope: rosaScope,
	}
}

func NewRosaClient(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (*RosaClient, error) {
	var token string
	var ocmAPIUrl string

	secret := rosaScope.CredentialsSecret()
	if secret != nil {
		if err := rosaScope.Client.Get(ctx, client.ObjectKeyFromObject(secret), secret); err != nil {
			return nil, fmt.Errorf("failed to get credentials secret: %w", err)
		}

		token = string(secret.Data[ocmTokenKey])
		ocmAPIUrl = string(secret.Data[ocmAPIURLKey])
	} else {
		// fallback to env variables if secrert is not set
		token = os.Getenv("OCM_TOKEN")
		if ocmAPIUrl = os.Getenv("OCM_API_URL"); ocmAPIUrl == "" {
			ocmAPIUrl = "https://api.openshift.com"
		}
	}

	if token == "" {
		return nil, fmt.Errorf("token is not provided, be sure to set OCM_TOKEN env variable or reference a credentials secret with key %s", ocmTokenKey)
	}

	// Create a logger that has the debug level enabled:
	logger, err := sdk.NewGoLoggerBuilder().
		Debug(true).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	connection, err := sdk.NewConnectionBuilder().
		Logger(logger).
		Tokens(token).
		URL(ocmAPIUrl).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create ocm connection: %w", err)
	}

	return &RosaClient{
		ocm:       connection,
		rosaScope: rosaScope,
	}, nil
}

func (c *RosaClient) Close() error {
	return c.ocm.Close()
}

func (c *RosaClient) GetConnectionURL() string {
	return c.ocm.URL()
}

func (c *RosaClient) GetConnectionTokens() (string, string, error) {
	return c.ocm.Tokens()
}
