// Package aws collects AWS-specific configuration.
package aws

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/aws/validation"
)

var (
	defaultVPCCIDR = ipnet.MustParseCIDR("10.0.0.0/16")
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
		VPCCIDRBlock: defaultVPCCIDR,
		Region:       region,
	}, nil
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

	tmpl, err := template.New("aws-credentials").Parse(`# Created by openshift-install
# https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html
[default]
aws_access_key_id={{.KeyID}}
aws_secret_access_key={{.SecretKey}}
`)
	if err != nil {
		return err
	}

	path := defaults.SharedCredentialsFilename()
	logrus.Infof("Writing AWS credentials to %q (https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html)", path)
	err = os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, map[string]string{
		"KeyID":     keyID,
		"SecretKey": secretKey,
	})
}
