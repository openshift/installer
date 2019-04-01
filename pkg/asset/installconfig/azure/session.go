package azure

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Azure/go-autorest/autorest"
	azureenv "github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/AlecAivazis/survey.v1"
)

var authFileLocation = os.Getenv("HOME") + "/.azure/osServicePrincipal.json"

//Session is an object representing session for subscription
type Session struct {
	SubscriptionID string
	Authorizer     autorest.Authorizer
}

// GetSession returns an azure session by checking credentials
// and, if no creds are found, asks for them and stores them on disk in a config file
func GetSession() (*Session, error) {
	err := getCredentials()
	if err != nil {
		return nil, err
	}
	return newSessionFromFile()
}

func newSessionFromFile() (*Session, error) {
	os.Setenv("AZURE_AUTH_LOCATION", authFileLocation)
	authorizer, err := auth.NewAuthorizerFromFileWithResource(azureenv.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "Can't initialize authorizer")
	}

	authInfo, err := readCredentialsFromFile()
	if err != nil {
		return nil, err
	}

	session := &Session{
		SubscriptionID: authInfo.SubscriptionID,
		Authorizer:     authorizer,
	}
	return session, nil
}

func readCredentialsFromFile() (*azureCredentials, error) {
	authInfo := &azureCredentials{}
	authBytes, err := ioutil.ReadFile(authFileLocation)
	if err != nil {
		return nil, errors.Wrapf(err, "Can't read azure authorization file : %s", authFileLocation)
	}
	err = json.Unmarshal(authBytes, authInfo)
	if err != nil {
		return nil, errors.Wrap(err, "Can't get authinfo")
	}
	return authInfo, nil
}

type azureCredentials struct {
	SubscriptionID string `json:"subscriptionId,omitempty"`
	ClientID       string `json:"clientId,omitempty"`
	ClientSecret   string `json:"clientSecret,omitempty"`
	TenantID       string `json:"tenantId,omitempty"`
}

func getCredentials() error {
	if _, err := os.Stat(authFileLocation); err == nil {
		return nil
	}

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
		return err
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
		return err
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
		return err
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
		return err
	}

	jsonCreds, err := json.Marshal(azureCredentials{
		SubscriptionID: subscriptionID,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		TenantID:       tenantID,
	})

	logrus.Infof("Writing azure credentials to %q", authFileLocation)
	err = os.MkdirAll(filepath.Dir(authFileLocation), 0700)
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = ioutil.WriteFile(authFileLocation, jsonCreds, 0600)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
