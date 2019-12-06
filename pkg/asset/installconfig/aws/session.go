package aws

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/openshift/installer/pkg/version"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	survey "gopkg.in/AlecAivazis/survey.v1"
	ini "gopkg.in/ini.v1"
)

const (
	sharedCredentialsProviderName = "SharedCredentialsProvider"
	envProviderName               = "EnvProvider"
)

var (
	onceLoggers = map[string]*sync.Once{
		sharedCredentialsProviderName: new(sync.Once),
		envProviderName:               new(sync.Once),
	}
)

// GetSession returns an AWS session by checking credentials
// and, if no creds are found, asks for them and stores them on disk in a config file
func GetSession() (*session.Session, error) {
	ssn := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	sharedCredentialsProvider := &credentials.SharedCredentialsProvider{}
	ssn.Config.Credentials = credentials.NewChainCredentials([]credentials.Provider{
		&credentials.EnvProvider{},
		sharedCredentialsProvider,
	})

	creds, err := ssn.Config.Credentials.Get()
	if err == nil {
		switch creds.ProviderName {
		case sharedCredentialsProviderName:
			onceLoggers[sharedCredentialsProviderName].Do(func() {
				logrus.Infof("Credentials loaded from the %q profile in file %q", sharedCredentialsProvider.Profile, sharedCredentialsProvider.Filename)
			})
		case envProviderName:
			onceLoggers[envProviderName].Do(func() {
				logrus.Info("Credentials loaded from default AWS environment variables")
			})
		}
	}
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
