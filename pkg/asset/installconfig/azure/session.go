package azure

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/go-autorest/autorest"
	azureenv "github.com/Azure/go-autorest/autorest/azure"
	"github.com/jongio/azidext/go/azidext"
	azurekiota "github.com/microsoft/kiota-authentication-azure-go"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/types/azure"
)

const azureAuthEnv = "AZURE_AUTH_LOCATION"

var (
	defaultAuthFilePath = filepath.Join(os.Getenv("HOME"), ".azure", "osServicePrincipal.json")
	onceLoggers         = map[string]*sync.Once{}
)

// AuthenticationType identifies the authentication method used.
type AuthenticationType int

// The authentication types supported by the installer.
const (
	ClientSecretAuth AuthenticationType = iota
	ClientCertificateAuth
	ManagedIdentityAuth
	// AzureCLIAuth is authentication via az login (azidentity.NewAzureCLICredential).
	AzureCLIAuth
)

// Session is an object representing session for subscription
type Session struct {
	Authorizer   autorest.Authorizer
	Credentials  Credentials
	Environment  azureenv.Environment
	AuthProvider *azurekiota.AzureIdentityAuthenticationProvider
	TokenCreds   azcore.TokenCredential
	CloudConfig  cloud.Configuration
	AuthType     AuthenticationType
}

// Credentials is the data type for credentials as understood by the azure sdk
type Credentials struct {
	SubscriptionID            string `json:"subscriptionId,omitempty"`
	ClientID                  string `json:"clientId,omitempty"`
	ClientSecret              string `json:"clientSecret,omitempty"`
	TenantID                  string `json:"tenantId,omitempty"`
	ClientCertificatePath     string `json:"clientCertificate,omitempty"`
	ClientCertificatePassword string `json:"clientCertificatePassword,omitempty"`
}

// GetSession returns an azure session by using credentials found in ~/.azure/osServicePrincipal.json.
// If no credentials file is found, it tries the Azure CLI profile (az login). If that is also
// unavailable, it asks for credentials and stores them on disk in a config file.
func GetSession(cloudName azure.CloudEnvironment, armEndpoint string) (*Session, error) {
	return GetSessionWithCredentials(cloudName, armEndpoint, nil)
}

// GetSessionWithCredentials returns an Azure session by using prepopulated credentials.
// If there are no prepopulated credentials it falls back to reading credentials from the
// credentials file, the Azure CLI profile (az login), or asking the user and storing them on disk.
func GetSessionWithCredentials(cloudName azure.CloudEnvironment, armEndpoint string, credentials *Credentials) (*Session, error) {
	var cloudEnv azureenv.Environment
	var err error
	switch cloudName {
	case azure.StackCloud:
		cloudEnv, err = azureenv.EnvironmentFromURL(armEndpoint)
	default:
		cloudEnv, err = azureenv.EnvironmentFromName(string(cloudName))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get Azure environment for the %q cloud: %w", cloudName, err)
	}

	cloudConfig, err := GetCloudConfiguration(cloudName, armEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get cloud configuration for the %q cloud: %w", cloudName, err)
	}

	var cred azcore.TokenCredential
	var authType AuthenticationType
	if credentials == nil {
		credentials, authType, err = credentialsFromFileOrUser()
		if err != nil {
			return nil, err
		}
	} else {
		authType = authTypeFromCredentials(credentials)
	}
	switch authType {
	case ClientCertificateAuth:
		logrus.Warnf("Using client certs to authenticate. Please be warned cluster does not support certs and only the installer does.")
		cred, err = newTokenCredentialFromCertificates(credentials, *cloudConfig)
	case ClientSecretAuth:
		cred, err = newTokenCredentialFromCredentials(credentials, *cloudConfig)
	case AzureCLIAuth:
		logrus.Infof("Using Azure CLI credentials from az login")
		cred, err = newTokenCredentialFromAzureCLI(credentials)
	default:
		cred, err = newTokenCredentialFromMSI(credentials, *cloudConfig)
	}
	if err != nil {
		return nil, err
	}
	session, err := newSessionFromCredentials(cloudEnv, credentials, cred)
	if err != nil {
		return nil, err
	}
	session.CloudConfig = *cloudConfig
	session.AuthType = authType
	return session, nil
}

// GetCloudConfiguration gets a cloud configuration from the cloud name and endpoint.
func GetCloudConfiguration(cloudName azure.CloudEnvironment, armEndpoint string) (*cloud.Configuration, error) {
	var cloudEnv azureenv.Environment
	var err error
	switch cloudName {
	case azure.StackCloud:
		cloudEnv, err = azureenv.EnvironmentFromURL(armEndpoint)
	default:
		cloudEnv, err = azureenv.EnvironmentFromName(string(cloudName))
	}
	if err != nil {
		return nil, err
	}

	var cloudConfig cloud.Configuration
	switch cloudName {
	case azure.StackCloud:
		cloudConfig = cloud.Configuration{
			ActiveDirectoryAuthorityHost: cloudEnv.ActiveDirectoryEndpoint,
			Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
				cloud.ResourceManager: {
					Audience: cloudEnv.TokenAudience,
					Endpoint: cloudEnv.ResourceManagerEndpoint,
				},
			},
		}
	case azure.USGovernmentCloud:
		cloudConfig = cloud.AzureGovernment
	case azure.ChinaCloud:
		cloudConfig = cloud.AzureChina
	default:
		cloudConfig = cloud.AzurePublic
	}

	return &cloudConfig, nil
}

// authTypeFromCredentials returns the authentication method for a credentials
// file or prepopulated Credentials. Certificate and secret take priority;
// otherwise managed identity is used, including system-assigned MSI when ClientID is empty.
func authTypeFromCredentials(credentials *Credentials) AuthenticationType {
	switch {
	case credentials.ClientCertificatePath != "":
		return ClientCertificateAuth
	case credentials.ClientSecret != "":
		return ClientSecretAuth
	default:
		return ManagedIdentityAuth
	}
}

// credentialsFromFileOrUser returns credentials from ~/.azure/osServicePrincipal.json
// or AZURE_AUTH_LOCATION. If no credentials file is found, it tries the Azure CLI
// profile and otherwise prompts the user and stores a service principal on disk.
func credentialsFromFileOrUser() (*Credentials, AuthenticationType, error) {
	authFilePath := defaultAuthFilePath
	if f := os.Getenv(azureAuthEnv); len(f) > 0 {
		authFilePath = f
	}

	var authFile Credentials

	contents, err := os.ReadFile(authFilePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			cliCreds, cliErr := credentialsFromAzureCLIProfile()
			switch {
			case cliErr == nil:
				if cerr := checkCredentials(*cliCreds); cerr != nil {
					logrus.Warnf("Azure CLI profile found but incomplete: %v", cerr)
				} else {
					return cliCreds, AzureCLIAuth, nil
				}
			case errors.Is(cliErr, fs.ErrNotExist):
				logrus.Debugf("Azure CLI profile not found: %v", cliErr)
			default:
				logrus.Warnf("Azure CLI profile unusable: %v", cliErr)
			}
			// Fall back to asking the user interactively
			logrus.Infof("Asking user to provide authentication info")
			credentials, cerr := askForCredentials()
			if cerr != nil {
				return nil, 0, fmt.Errorf("failed to retrieve credentials from user: %w", cerr)
			}
			logrus.Infof("Saving user credentials to %q", authFilePath)
			if cerr = saveCredentials(*credentials, authFilePath); cerr != nil {
				return nil, 0, fmt.Errorf("failed to save credentials: %w", cerr)
			}
			authFile = *credentials
		} else {
			// File was found but we failed to read it, just error out and let the user handle it
			return nil, 0, err
		}
	} else {
		err = json.Unmarshal(contents, &authFile)
		if err != nil {
			return nil, 0, err
		}
	}

	if err := checkCredentials(authFile); err != nil {
		return nil, 0, err
	}

	if _, has := onceLoggers[authFilePath]; !has {
		onceLoggers[authFilePath] = new(sync.Once)
	}
	onceLoggers[authFilePath].Do(func() {
		logrus.Infof("Credentials loaded from file %q", authFilePath)
	})

	return &authFile, authTypeFromCredentials(&authFile), nil
}

func checkCredentials(creds Credentials) error {
	if creds.SubscriptionID == "" {
		return errors.New("could not retrieve subscriptionId from auth file")
	}
	if creds.TenantID == "" {
		return errors.New("could not retrieve tenantId from auth file")
	}
	if (creds.ClientSecret != "" || creds.ClientCertificatePath != "") && creds.ClientID == "" {
		return errors.New("could not retrieve clientId from auth file")
	}
	return nil
}

func askForCredentials() (*Credentials, error) {
	var subscriptionID, tenantID, clientID, clientSecret string

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "azure subscription id",
				Help:    "The azure subscription id to use for installation",
			},
		},
	}, &subscriptionID)
	if err != nil {
		return nil, err
	}

	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "azure tenant id",
				Help:    "The azure tenant id to use for installation",
			},
		},
	}, &tenantID)
	if err != nil {
		return nil, err
	}

	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "azure service principal client id",
				Help:    "The azure client id to use for installation (this is not your username)",
			},
		},
	}, &clientID)
	if err != nil {
		return nil, err
	}

	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Password{
				Message: "azure service principal client secret",
				Help:    "The azure secret access key corresponding to your client secret (this is not your password).",
			},
		},
	}, &clientSecret)
	if err != nil {
		return nil, err
	}

	return &Credentials{
		SubscriptionID: subscriptionID,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		TenantID:       tenantID,
	}, nil
}

func saveCredentials(credentials Credentials, filePath string) error {
	jsonCreds, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(filePath), 0700)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, jsonCreds, 0o600)
}

func newTokenCredentialFromCredentials(credentials *Credentials, cloudConfig cloud.Configuration) (azcore.TokenCredential, error) {
	options := azidentity.ClientSecretCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloudConfig,
		},
	}

	cred, err := azidentity.NewClientSecretCredential(credentials.TenantID, credentials.ClientID, credentials.ClientSecret, &options)
	if err != nil {
		return nil, fmt.Errorf("failed to get client credentials from secret: %w", err)
	}
	return cred, nil
}

func newTokenCredentialFromCertificates(credentials *Credentials, cloudConfig cloud.Configuration) (azcore.TokenCredential, error) {
	options := azidentity.ClientCertificateCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloudConfig,
		},
		SendCertificateChain: true,
	}

	data, err := os.ReadFile(credentials.ClientCertificatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read client certificate file: %w", err)
	}

	// NewClientCertificateCredential requires at least one *x509.Certificate,
	// and a crypto.PrivateKey. ParseCertificates returns these given
	// certificate data in PEM or PKCS12 format. It handles common scenarios
	// but has limitations, for example it doesn't load PEM encrypted private
	// keys.
	certs, key, err := azidentity.ParseCertificates(data, []byte(credentials.ClientCertificatePassword))
	if err != nil {
		return nil, fmt.Errorf("failed to parse client certificate: %w", err)
	}

	cred, err := azidentity.NewClientCertificateCredential(credentials.TenantID, credentials.ClientID, certs, key, &options)
	if err != nil {
		return nil, fmt.Errorf("failed to get client credentials from certificate: %w", err)
	}
	return cred, nil
}

func newTokenCredentialFromMSI(credentials *Credentials, cloudConfig cloud.Configuration) (azcore.TokenCredential, error) {
	options := azidentity.ManagedIdentityCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloudConfig,
		},
	}
	// User-assigned identity
	if credentials.ClientID != "" {
		options.ID = azidentity.ClientID(credentials.ClientID)
	}

	cred, err := azidentity.NewManagedIdentityCredential(&options)
	if err != nil {
		return nil, fmt.Errorf("failed to get client credentials from MSI: %w", err)
	}
	return cred, nil
}

// newTokenCredentialFromAzureCLI returns a TokenCredential for Azure CLI (az login).
func newTokenCredentialFromAzureCLI(credentials *Credentials) (azcore.TokenCredential, error) {
	cred, err := azidentity.NewAzureCLICredential(azureCLICredentialOptions(credentials))
	if err != nil {
		return nil, fmt.Errorf("failed to get Azure CLI credentials (have you run 'az login'?): %w", err)
	}
	return cred, nil
}

// azureCLICredentialOptions returns options for NewAzureCLICredential.
// Subscription is preferred when set; otherwise TenantID is used; otherwise
// nil selects the Azure CLI current account. Both cannot be set because
// Azure CLI rejects get-access-token when --subscription and --tenant are
// passed together.
func azureCLICredentialOptions(credentials *Credentials) *azidentity.AzureCLICredentialOptions {
	switch {
	case credentials.SubscriptionID != "":
		return &azidentity.AzureCLICredentialOptions{Subscription: credentials.SubscriptionID}
	case credentials.TenantID != "":
		return &azidentity.AzureCLICredentialOptions{TenantID: credentials.TenantID}
	}
	return nil
}

// credentialsFromAzureCLIProfile reads the default subscription and tenant IDs from
// ~/.azure/azureProfile.json so the installer can authenticate with AzureCLICredential
// without a service principal.
func credentialsFromAzureCLIProfile() (*Credentials, error) {
	// HOME/.azure is the same trusted CLI config dir as osServicePrincipal.json.
	data, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".azure", "azureProfile.json")) // #nosec G703
	if err != nil {
		return nil, err
	}
	// Strip UTF-8 BOM if present (Azure CLI writes one on Windows)
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))

	var profile struct {
		Subscriptions []struct {
			ID        string `json:"id"`
			TenantID  string `json:"tenantId"`
			IsDefault bool   `json:"isDefault"`
		} `json:"subscriptions"`
	}
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, err
	}
	for _, sub := range profile.Subscriptions {
		if sub.IsDefault {
			return &Credentials{
				SubscriptionID: sub.ID,
				TenantID:       sub.TenantID,
			}, nil
		}
	}
	return nil, fmt.Errorf("no default subscription found in Azure CLI profile")
}

func newSessionFromCredentials(cloudEnv azureenv.Environment, credentials *Credentials, cred azcore.TokenCredential) (*Session, error) {
	var scope []string
	// This can be empty for StackCloud
	if cloudEnv.MicrosoftGraphEndpoint != "" {
		// GovClouds need a properly set scope in the authenticator, otherwise we
		// get an 'Invalid audience' error when doing MSGraph API calls
		// https://learn.microsoft.com/en-us/graph/sdks/national-clouds?tabs=go
		scope = []string{endpointToScope(cloudEnv.MicrosoftGraphEndpoint)}
	}
	authProvider, err := azurekiota.NewAzureIdentityAuthenticationProviderWithScopes(cred, scope)
	if err != nil {
		return nil, fmt.Errorf("failed to get Azidentity authentication provider: %w", err)
	}

	// Use an adapter so azidentity in the Azure SDK can be used as
	// Authorizer when calling the Azure Management Packages, which we
	// currently use. Once the Azure SDK clients (found in /sdk) move to
	// stable, we can update our clients and they will be able to use the
	// creds directly without the authorizer. The schedule is here:
	// https://azure.github.io/azure-sdk/releases/latest/index.html#go
	authorizer := azidext.NewTokenCredentialAdapter(cred, []string{endpointToScope(cloudEnv.TokenAudience)})

	return &Session{
		Authorizer:   authorizer,
		Credentials:  *credentials,
		Environment:  cloudEnv,
		AuthProvider: authProvider,
		TokenCreds:   cred,
	}, nil
}

func endpointToScope(endpoint string) string {
	if !strings.HasSuffix(endpoint, "/.default") {
		endpoint += "/.default"
	}
	return endpoint
}
