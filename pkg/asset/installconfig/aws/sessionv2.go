package aws

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/smithy-go/middleware"
	"github.com/sirupsen/logrus"
	ini "gopkg.in/ini.v1"

	"github.com/openshift/installer/pkg/version"
)

const (
	// OpenShiftInstallerUserAgent is the User Agent key to add to the AWS API request header.
	OpenShiftInstallerUserAgent = "OpenShift/4.x Installer"

	// OpenShiftInstallerGatherUserAgent is the User Agent key to add to the AWS API request header
	// when gather command is invoked.
	OpenShiftInstallerGatherUserAgent = "OpenShift/4.x Gather"

	// RetryMaxAttempts is the total number of times an API request is retried.
	RetryMaxAttempts = 25
)

var (
	credentialsFromConfigLogger = new(sync.Once)
)

// ConfigOptions is a set of functions that modify the provided config.LoadOptions.
type ConfigOptions []func(*config.LoadOptions) error

// getDefaultConfigOptions returns the default settings for config.LoadOptions.
func getDefaultConfigOptions() ConfigOptions {
	return ConfigOptions{
		config.WithRetryMaxAttempts(RetryMaxAttempts),
		config.WithAPIOptions([]func(*middleware.Stack) error{
			awsmiddleware.AddUserAgentKeyValue(OpenShiftInstallerUserAgent, version.Raw),
		}),
	}
}

// GetConfig returns an AWS config by checking credentials
// and, if no creds are found, asks for them and stores them on disk in a config file.
func GetConfig(ctx context.Context) (aws.Config, error) { return GetConfigWithOptions(ctx) }

// GetConfigWithOptions returns an AWS config by checking credentials
// and, if no creds are found, asks for them and stores them on disk in a config file.
func GetConfigWithOptions(ctx context.Context, options ...func(*config.LoadOptions) error) (aws.Config, error) {
	// Set the default options, which are overridden by user-defined ones if any.
	options = append(getDefaultConfigOptions(), options...)

	// Attempt to retrieve valid credentials.
	// If failed, ask the user for the credentials via the survey.
	_, err := getCredentialsV2(ctx, options)
	if err != nil {
		if err := getUserCredentialsV2(); err != nil {
			return aws.Config{}, err
		}
	}

	cfg, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to create AWS config: %w", err)
	}

	return cfg, err
}

// getCredentials returns the credentials by constructing an AWS Config
// and attempting to retrieve the credentials if possible.
// TODO: Remove suffix V2 when completing migration aws sdk v2 (i.e. removing session.go).
func getCredentialsV2(ctx context.Context, options ConfigOptions) (aws.Credentials, error) {
	cfg, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return aws.Credentials{}, fmt.Errorf("failed to create AWS config: %w", err)
	}

	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return aws.Credentials{}, err
	}

	credentialsFromConfigLogger.Do(func() {
		logrus.Infof("Credentials loaded from the AWS config using %q provider", creds.Source)
	})

	return creds, nil
}

// IsStaticCredentialsV2 returns whether the credentials value provider are
// static credentials safe for installer to transfer to cluster for use as-is.
// TODO: Remove suffix V2 when completing migration aws sdk v2 (i.e. removing session.go).
func IsStaticCredentialsV2(creds aws.Credentials) bool {
	if creds.Source == credentials.StaticCredentialsName {
		return creds.SessionToken == ""
	}
	return false
}

// getUserCredentialsV2 asks for aws access key id and secret in the survey
// and stores them on disk in a config file.
// TODO: Remove suffix V2 when completing migration aws sdk v2 (i.e. removing session.go).
func getUserCredentialsV2() error {
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

	path := config.DefaultSharedCredentialsFilename()
	if env := os.Getenv("AWS_SHARED_CREDENTIALS_FILE"); env != "" {
		path = env
	}
	logrus.Infof("Writing AWS credentials to %q (https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html)", path)
	err = os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		return err
	}

	creds, err := ini.Load(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to load credentials file %s: %w", path, err)
		}
		creds = ini.Empty()
		creds.Section("").Comment = "https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html"
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
			logrus.Error(fmt.Errorf("failed to remove partially-written credentials file: %w", err2))
		}
		return err
	}

	return os.Rename(tempPath, path)
}
