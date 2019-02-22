// Package aws collects AWS-specific configuration.
package aws

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/aws/validation"
	"github.com/openshift/installer/pkg/version"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	survey "gopkg.in/AlecAivazis/survey.v1"
	ini "gopkg.in/ini.v1"
)

// Platform collects AWS-specific configuration.
func Platform() (*aws.Platform, error) {
	longRegions := make([]string, 0, len(validation.Regions))
	shortRegions := make([]string, 0, len(validation.Regions))
	for id, location := range validation.Regions {
		longRegions = append(longRegions, fmt.Sprintf("%s (%s)", id, location))
		shortRegions = append(shortRegions, id)
	}
	regionTransform := survey.TransformString(func(s string) string {
		return strings.SplitN(s, " ", 2)[0]
	})

	defaultRegion := "us-east-1"
	_, ok := validation.Regions[defaultRegion]
	if !ok {
		panic(fmt.Sprintf("installer bug: invalid default AWS region %q", defaultRegion))
	}

	ssn, err := GetSession()
	if err != nil {
		return nil, err
	}

	defaultRegionPointer := ssn.Config.Region
	if defaultRegionPointer != nil && *defaultRegionPointer != "" {
		_, ok := validation.Regions[*defaultRegionPointer]
		if ok {
			defaultRegion = *defaultRegionPointer
		} else {
			logrus.Warnf("Unrecognized AWS region %q, defaulting to %s", *defaultRegionPointer, defaultRegion)
		}
	}

	sort.Strings(longRegions)
	sort.Strings(shortRegions)

	var region string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Region",
				Help:    "The AWS region to be used for installation.",
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

	return &aws.Platform{
		Region: region,
	}, nil
}

// GetSession returns an AWS session by checking credentials
// and, if no creds are found, asks for them and stores them on disk in a config file
func GetSession() (*session.Session, error) {
	ssn := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	ssn.Config.Credentials = credentials.NewChainCredentials([]credentials.Provider{
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{},
	})
	_, err := ssn.Config.Credentials.Get()
	if err == credentials.ErrNoValidProvidersFoundInChain {
		err = getCredentials()
		if err != nil {
			return nil, err
		}
	}
	ssn.Handlers.Build.PushBackNamed(request.NamedHandler{
		Name: "openshiftInstaller.OpenshiftInstallerUserAgentHandler",
		Fn:   request.MakeAddToUserAgentHandler("OpenShift/4.x Installer", version.Raw),
	})
	return ssn, nil
}

func getCredentials() error {
	var keyID string
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "AWS Access Key ID",
				Help:    "The AWS access key ID to use for installation (this is not your username).\nhttps://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html",
			},
		},
	}, &keyID)
	if err != nil {
		return err
	}

	var secretKey string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Password{
				Message: "AWS Secret Access Key",
				Help:    "The AWS secret access key corresponding to your access key ID (this is not your password).",
			},
		},
	}, &secretKey)
	if err != nil {
		return err
	}

	path := defaults.SharedCredentialsFilename()
	logrus.Infof("Writing AWS credentials to %q (https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html)", path)
	err = os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		return err
	}

	creds, err := ini.Load(path)
	if err != nil {
		if os.IsNotExist(err) {
			creds = ini.Empty()
			creds.Section("").Comment = "https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html"
		} else {
			return err
		}
	}

	profile := os.Getenv("AWS_PROFILE")
	if profile == "" {
		profile = "default"
	}

	creds.Section(profile).Key("aws_access_key_id").SetValue(keyID)
	creds.Section(profile).Key("aws_secret_access_key").SetValue(secretKey)

	tempPath := path + ".tmp"
	file, err := os.OpenFile(tempPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = creds.WriteTo(file)
	if err != nil {
		err2 := os.Remove(tempPath)
		if err2 != nil {
			logrus.Error(errors.Wrap(err2, "failed to remove partially-written credentials file"))
		}
		return err
	}

	return os.Rename(tempPath, path)
}
