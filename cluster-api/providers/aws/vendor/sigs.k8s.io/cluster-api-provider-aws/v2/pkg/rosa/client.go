// Package rosa provides a way to interact with the Red Hat OpenShift Service on AWS (ROSA) API.
package rosa

import (
	"context"
	"fmt"
	"os"

	sdk "github.com/openshift-online/ocm-sdk-go"
	ocmcfg "github.com/openshift/rosa/pkg/config"
	"github.com/openshift/rosa/pkg/ocm"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

const (
	ocmTokenKey  = "ocmToken"
	ocmAPIURLKey = "ocmApiUrl"
)

// NewOCMClient creates a new OCM client.
func NewOCMClient(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (*ocm.Client, error) {
	token, url, err := ocmCredentials(ctx, rosaScope)
	if err != nil {
		return nil, err
	}
	return ocm.NewClient().Logger(logrus.New()).Config(&ocmcfg.Config{
		AccessToken: token,
		URL:         url,
	}).Build()
}

func newOCMRawConnection(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (*sdk.Connection, error) {
	logger, err := sdk.NewGoLoggerBuilder().
		Debug(false).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}
	token, url, err := ocmCredentials(ctx, rosaScope)
	if err != nil {
		return nil, err
	}

	connection, err := sdk.NewConnectionBuilder().
		Logger(logger).
		Tokens(token).
		URL(url).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create ocm connection: %w", err)
	}

	return connection, nil
}

func ocmCredentials(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (string, string, error) {
	var token string
	var ocmAPIUrl string

	secret := rosaScope.CredentialsSecret()
	if secret != nil {
		if err := rosaScope.Client.Get(ctx, client.ObjectKeyFromObject(secret), secret); err != nil {
			return "", "", fmt.Errorf("failed to get credentials secret: %w", err)
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
		return "", "", fmt.Errorf("token is not provided, be sure to set OCM_TOKEN env variable or reference a credentials secret with key %s", ocmTokenKey)
	}
	return token, ocmAPIUrl, nil
}
