package azure

import (
	"encoding/json"
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
	"github.com/pkg/errors"
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

// GetSession returns an azure session by using credentials found in ~/.azure/osServicePrincipal.json
// and, if no creds are found, asks for them and stores them on disk in a config file
func GetSession(cloudName azure.CloudEnvironment, armEndpoint string) (*Session, error) {
	return GetSessionWithCredentials(cloudName, armEndpoint, nil)
}

// GetSessionWithCredentials returns an Azure session by using prepopulated credentials.
// If there are no prepopulated credentials it falls back to reading credentials from file system
// or from user input.
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
		return nil, errors.Wrapf(err, "failed to get Azure environment for the %q cloud", cloudName)
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

	if credentials == nil {
		credentials, err = credentialsFromFileOrUser()
		if err != nil {
			return nil, err
		}
	}
	var cred azcore.TokenCredential
	var authType AuthenticationType
	switch {
	case credentials.ClientCertificatePath != "":
		logrus.Warnf("Using client certs to authenticate. Please be warned cluster does not support certs and only the installer does.")
		cred, err = newTokenCredentialFromCertificates(credentials, cloudConfig)
		authType = ClientCertificateAuth
	case credentials.ClientSecret != "":
		cred, err = newTokenCredentialFromCredentials(credentials, cloudConfig)
		authType = ClientSecretAuth
	default:
		cred, err = newTokenCredentialFromMSI(credentials, cloudConfig)
		authType = ManagedIdentityAuth
	}
	if err != nil {
		return nil, err
	}
	session, err := newSessionFromCredentials(cloudEnv, credentials, cred)
	if err != nil {
		return nil, err
	}
	session.CloudConfig = cloudConfig
	session.AuthType = authType
	return session, nil
}

// credentialsFromFileOrUser returns credentials found
// in ~/.azure/osServicePrincipal.json and, if no creds are found,
// asks for them and stores them on disk in a config file
func credentialsFromFileOrUser() (*Credentials, error) {
	authFilePath := defaultAuthFilePath
	if f := os.Getenv(azureAuthEnv); len(f) > 0 {
		authFilePath = f
	}

	var authFile Credentials

	contents, err := os.ReadFile(authFilePath)
	if err != nil {
		// If the file with creds was not found, ask user for auth info
		if errors.Is(err, fs.ErrNotExist) {
			logrus.Infof("Asking user to provide authentication info")
			credentials, cerr := askForCredentials()
			if cerr != nil {
				return nil, errors.Wrap(cerr, "failed to retrieve credentials from user")
			}
			logrus.Infof("Saving user credentials to %q", authFilePath)
			if cerr = saveCredentials(*credentials, authFilePath); cerr != nil {
				return nil, errors.Wrap(cerr, "failed to save credentials")
			}
			authFile = *credentials
		} else {
			// File was found but we failed to read it, just error out and let the user handle it
			return nil, err
		}
	} else {
		err = json.Unmarshal(contents, &authFile)
		if err != nil {
			return nil, err
		}
	}

	if err := checkCredentials(authFile); err != nil {
		return nil, err
	}

	if _, has := onceLoggers[authFilePath]; !has {
		onceLoggers[authFilePath] = new(sync.Once)
	}
	onceLoggers[authFilePath].Do(func() {
		logrus.Infof("Credentials loaded from file %q", authFilePath)
	})

	return &authFile, nil
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
	// If neither client secret nor client certificate are present, we default to Managed Identity
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
		return nil, errors.Wrap(err, "failed to get client credentials from secret")
	}
	return cred, nil
}

func newTokenCredentialFromCertificates(credentials *Credentials, cloudConfig cloud.Configuration) (azcore.TokenCredential, error) {
	options := azidentity.ClientCertificateCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloudConfig,
		},
	}

	data, err := os.ReadFile(credentials.ClientCertificatePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read client certificate file")
	}

	// NewClientCertificateCredential requires at least one *x509.Certificate,
	// and a crypto.PrivateKey. ParseCertificates returns these given
	// certificate data in PEM or PKCS12 format. It handles common scenarios
	// but has limitations, for example it doesn't load PEM encrypted private
	// keys.
	certs, key, err := azidentity.ParseCertificates(data, []byte(credentials.ClientCertificatePassword))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse client certificate")
	}

	cred, err := azidentity.NewClientCertificateCredential(credentials.TenantID, credentials.ClientID, certs, key, &options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get client credentials from certificate")
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
		return nil, errors.Wrap(err, "failed to get client credentials from MSI")
	}
	return cred, nil
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
		return nil, errors.Wrap(err, "failed to get Azidentity authentication provider")
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
