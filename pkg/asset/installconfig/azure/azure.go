package azure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/azure/validation"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// Platform collects azure-specific configuration.
func Platform() (*azure.Platform, error) {
	longRegions := make([]string, 0, len(validation.Regions))
	shortRegions := make([]string, 0, len(validation.Regions))
	for id, location := range validation.Regions {
		longRegions = append(longRegions, fmt.Sprintf("%s (%s)", id, location))
		shortRegions = append(shortRegions, id)
	}
	regionTransform := survey.TransformString(func(s string) string {
		return strings.SplitN(s, " ", 2)[0]
	})

	defaultRegion := "eastus"
	_, ok := validation.Regions[defaultRegion]
	if !ok {
		panic(fmt.Sprintf("installer bug: invalid default azure region %q", defaultRegion))
	}

	err := GetSession()
	if err != nil {
		return nil, err
	}

	sort.Strings(longRegions)
	sort.Strings(shortRegions)

	var region string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Region",
				Help:    "The azure region to be used for installation.",
				Default: fmt.Sprintf("%s (%s)", defaultRegion, validation.Regions[defaultRegion]),
				Options: longRegions,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				choice := regionTransform(ans).(string)
				i := sort.SearchStrings(shortRegions, choice)
				if i == len(shortRegions) || shortRegions[i] != choice {
					return errors.Errorf("invalid region %q", choice)
				}
				return nil
			}),
			Transform: regionTransform,
		},
	}, &region)
	if err != nil {
		return nil, err
	}

	return &azure.Platform{
		Region: region,
	}, nil
}

// GetSession returns an azure session by checking credentials
// and, if no creds are found, asks for them and stores them on disk in a config file
func GetSession() error {
	return getCredentials()
}

type AzureCredentials struct {
	SubscriptionID string `json:"subscriptionID,omitempty"`
	ClientID       string `json:"clientID,omitempty"`
	ClientSecret   string `json:"clientSecret,omitempty"`
	TenantID       string `json:"tenantID,omitempty"`
}

func getCredentials() error {
	path := os.Getenv("HOME") + "/.azure/osServicePrincipal.json"
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	}

	var subscriptionID, clientID, clientSecret string

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

	jsonCreds, err := json.Marshal(AzureCredentials{
		SubscriptionID: subscriptionID,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		TenantID:       "72f988bf-86f1-41af-91ab-2d7cd011db47",
	})

	logrus.Infof("Writing azure credentials to %q", path)
	err = os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, jsonCreds, 0600)
	if err != nil {
		return err
	}
	return nil
}
