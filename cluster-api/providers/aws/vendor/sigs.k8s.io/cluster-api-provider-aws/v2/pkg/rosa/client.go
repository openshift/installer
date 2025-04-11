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
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/util/system"
	"sigs.k8s.io/cluster-api-provider-aws/v2/version"
)

const (
	ocmTokenKey        = "ocmToken"
	ocmAPIURLKey       = "ocmApiUrl"
	ocmClientIDKey     = "ocmClientID"
	ocmClientSecretKey = "ocmClientSecret"
	capaAgentName      = "CAPA"
)

// NewOCMClient creates a new OCM client.
func NewOCMClient(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (*ocm.Client, error) {
	token, url, clientID, clientSecret, err := ocmCredentials(ctx, rosaScope)
	if err != nil {
		return nil, err
	}

	ocmConfig := ocmcfg.Config{
		URL:       url,
		UserAgent: capaAgentName,
		Version:   version.Get().GitVersion,
	}

	if clientID != "" && clientSecret != "" {
		ocmConfig.ClientID = clientID
		ocmConfig.ClientSecret = clientSecret
	} else if token != "" {
		ocmConfig.AccessToken = token
	}

	return ocm.NewClient().Logger(logrus.New()).Config(&ocmConfig).Build()
}

// NewWrappedOCMClient creates a new OCM client wrapped in ocmclient struct that implements OCMClient interface.
// This is needed to be able to mock OCM in tests. NewOCMClient is left unchanged so we don't change public interface.
func NewWrappedOCMClient(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (OCMClient, error) {
	ocmClient, err := NewOCMClient(ctx, rosaScope)
	c := ocmclient{
		ocmClient: ocmClient,
	}

	return &c, err
}

func newOCMRawConnection(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (*sdk.Connection, error) {
	ocmSdkLogger, err := sdk.NewGoLoggerBuilder().
		Debug(false).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	token, url, clientID, clientSecret, err := ocmCredentials(ctx, rosaScope)
	if err != nil {
		return nil, err
	}

	connBuilder := sdk.NewConnectionBuilder().
		Logger(ocmSdkLogger).
		URL(url).
		Agent(capaAgentName + "/" + version.Get().GitVersion + " " + sdk.DefaultAgent)

	if clientID != "" && clientSecret != "" {
		connBuilder.Client(clientID, clientSecret)
	} else if token != "" {
		connBuilder.Tokens(token)
	}

	connection, err := connBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create ocm connection: %w", err)
	}

	return connection, nil
}

func ocmCredentials(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (string, string, string, string, error) {
	var token string           // Offline SSO token
	var ocmClientID string     // Service account client id
	var ocmClientSecret string // Service account client secret
	var ocmAPIUrl string       // https://api.openshift.com by default
	var secret *corev1.Secret

	secret = rosaScope.CredentialsSecret() // We'll retrieve the OCM credentials ref from the ROSA control plane
	if secret != nil {
		if err := rosaScope.Client.Get(ctx, client.ObjectKeyFromObject(secret), secret); err != nil {
			return "", "", "", "", fmt.Errorf("failed to get credentials secret: %w", err)
		}
	} else { // If the reference to OCM secret wasn't specified in the ROSA control plane, we'll try to use a predefined secret name from the capa namespace
		secret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "rosa-creds-secret",
				Namespace: system.GetManagerNamespace(),
			},
		}

		err := rosaScope.Client.Get(ctx, client.ObjectKeyFromObject(secret), secret)
		// We'll ignore non-existent secret so that we can try the ENV variable fallback below
		// TODO: once the ENV variable fallback is gone, we can no longer ignore non-existent secret here
		if err != nil && !apierrors.IsNotFound(err) {
			return "", "", "", "", fmt.Errorf("failed to get credentials secret: %w", err)
		}
	}

	token = string(secret.Data[ocmTokenKey])
	ocmAPIUrl = string(secret.Data[ocmAPIURLKey])
	ocmClientID = string(secret.Data[ocmClientIDKey])
	ocmClientSecret = string(secret.Data[ocmClientSecretKey])

	// Deprecation warning in case SSO offline token was used
	if token != "" {
		rosaScope.Info("Using SSO offline token is deprecated, use service account credentials instead")
	}

	if token == "" && (ocmClientID == "" || ocmClientSecret == "") {
		// TODO: the ENV variables are to be removed with the next code release
		// Last fall-back is to use OCM_TOKEN & OCM_API_URL environment variables (soon to be deprecated)
		token = os.Getenv("OCM_TOKEN")
		ocmAPIUrl = os.Getenv("OCM_API_URL")

		if token != "" {
			rosaScope.Info("Defining OCM credentials in environment variable is deprecated, use secret with service account credentials instead")
		} else {
			return "", "", "", "",
				fmt.Errorf("OCM credentials have not been provided. Make sure to set the secret with service account credentials")
		}
	}

	if ocmAPIUrl == "" {
		ocmAPIUrl = "https://api.openshift.com" // Defaults to production URL
	}

	return token, ocmAPIUrl, ocmClientID, ocmClientSecret, nil
}
