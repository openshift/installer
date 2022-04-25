package azure

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Azure/go-autorest/autorest"
	azureenv "github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/types/azure"
)

const azureAuthEnv = "AZURE_AUTH_LOCATION"

var (
	defaultAuthFilePath = filepath.Join(os.Getenv("HOME"), ".azure", "osServicePrincipal.json")
	onceLoggers         = map[string]*sync.Once{}
)

//Session is an object representing session for subscription
type Session struct {
	GraphAuthorizer autorest.Authorizer
	Authorizer      autorest.Authorizer
	Credentials     Credentials
	Environment     azureenv.Environment
}

//Credentials is the data type for credentials as understood by the azure sdk
type Credentials struct {
	SubscriptionID string `json:"subscriptionId,omitempty"`
	ClientID       string `json:"clientId,omitempty"`
	ClientSecret   string `json:"clientSecret,omitempty"`
	TenantID       string `json:"tenantId,omitempty"`
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

	if credentials == nil {
		credentials, err = credentialsFromFileOrUser(&cloudEnv)
		if err != nil {
			return nil, err
		}
	}

	return newSessionFromCredentials(cloudEnv, credentials)
}

// credentialsFromFileOrUser returns credentials found
// in ~/.azure/osServicePrincipal.json and, if no creds are found,
// asks for them and stores them on disk in a config file
func credentialsFromFileOrUser(cloudEnv *azureenv.Environment) (*Credentials, error) {
	authFilePath := defaultAuthFilePath
	if f := os.Getenv(azureAuthEnv); len(f) > 0 {
		authFilePath = f
	}
	// NewAuthorizerFromFileWithResource uses `auth.GetSettingsFromFile`, which uses the `azureAuthEnv` to fetch the auth credentials.
	// therefore setting the local env here to authFilePath allows NewAuthorizerFromFileWithResource to load credentials.
	os.Setenv(azureAuthEnv, authFilePath)
	_, err := auth.NewAuthorizerFromFileWithResource(cloudEnv.ResourceManagerEndpoint)
	if err != nil {
		logrus.Debug("Could not get an azure authorizer from file. Asking user to provide authentication info")
		credentials, err := askForCredentials()
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve credentials from user")
		}
		logrus.Infof("Saving user credentials to %q", authFilePath)
		if err = saveCredentials(*credentials, authFilePath); err != nil {
			return nil, errors.Wrap(err, "failed to save credentials")
		}
	}

	//If the authorizer worked right away, we need to read credentials details
	authSettings, err := auth.GetSettingsFromFile()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get settings from file")
	}

	credentials, err := getCredentials(authSettings)
	if err != nil {
		return nil, errors.Wrap(err, "failed to map authsettings to credentials")
	}

	if _, has := onceLoggers[authFilePath]; !has {
		onceLoggers[authFilePath] = new(sync.Once)
	}
	onceLoggers[authFilePath].Do(func() {
		logrus.Infof("Credentials loaded from file %q", authFilePath)
	})

	return credentials, nil
}

func getCredentials(fs auth.FileSettings) (*Credentials, error) {
	subscriptionID := fs.GetSubscriptionID()
	if subscriptionID == "" {
		return nil, errors.New("could not retrieve subscriptionId from auth file")
	}

	clientID := fs.Values[auth.ClientID]
	if clientID == "" {
		return nil, errors.New("could not retrieve clientId from auth file")
	}
	clientSecret := fs.Values[auth.ClientSecret]
	if clientSecret == "" {
		return nil, errors.New("could not retrieve clientSecret from auth file")
	}
	tenantID := fs.Values[auth.TenantID]
	if tenantID == "" {
		return nil, errors.New("could not retrieve tenantId from auth file")
	}
	return &Credentials{
		SubscriptionID: subscriptionID,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		TenantID:       tenantID,
	}, nil
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

	return ioutil.WriteFile(filePath, jsonCreds, 0600)
}

func newSessionFromCredentials(cloudEnv azureenv.Environment, credentials *Credentials) (*Session, error) {
	c := &auth.ClientCredentialsConfig{
		TenantID:     credentials.TenantID,
		ClientID:     credentials.ClientID,
		ClientSecret: credentials.ClientSecret,
		AADEndpoint:  cloudEnv.ActiveDirectoryEndpoint,
	}
	c.Resource = cloudEnv.TokenAudience
	authorizer, err := c.Authorizer()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get client credentials authorizer")
	}

	c.Resource = cloudEnv.GraphEndpoint
	graphAuthorizer, err := c.Authorizer()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get GraphEndpoint authorizer")
	}
	return &Session{
		GraphAuthorizer: graphAuthorizer,
		Authorizer:      authorizer,
		Credentials:     *credentials,
		Environment:     cloudEnv,
	}, nil
}
