package aws

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	survey "gopkg.in/AlecAivazis/survey.v1"
	ini "gopkg.in/ini.v1"

	typesaws "github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/version"
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

// SessionOptions is a function that modifies the provided session.Option.
type SessionOptions func(sess *session.Options)

// WithRegion configures the session.Option to set the AWS region.
func WithRegion(region string) SessionOptions {
	return func(sess *session.Options) {
		cfg := aws.NewConfig().WithRegion(region)
		sess.Config.MergeIn(cfg)
	}
}

// WithServiceEndpoints configures the session.Option to use provides services for AWS endpoints.
func WithServiceEndpoints(region string, services []typesaws.ServiceEndpoint) SessionOptions {
	return func(sess *session.Options) {
		resolver := newAWSResolver(region, services)
		cfg := aws.NewConfig().WithEndpointResolver(resolver)
		sess.Config.MergeIn(cfg)
	}
}

// GetSession returns an AWS session by checking credentials
// and, if no creds are found, asks for them and stores them on disk in a config file
func GetSession() (*session.Session, error) { return GetSessionWithOptions() }

// GetSessionWithOptions returns an AWS session by checking credentials
// and, if no creds are found, asks for them and stores them on disk in a config file
func GetSessionWithOptions(optFuncs ...SessionOptions) (*session.Session, error) {
	options := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}
	for _, optFunc := range optFuncs {
		optFunc(&options)
	}

	ssn := session.Must(session.NewSessionWithOptions(options))

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
	ssn = ssn.Copy(&aws.Config{MaxRetries: aws.Int(25)})
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

type awsResolver struct {
	region   string
	services map[string]typesaws.ServiceEndpoint

	// this is a list of known default endpoints for specific regions that would
	// otherwise require user to set the service overrides.
	// it's a map of region => service => resolved endpoint
	// this is only used when the user hasn't specified a override for the service in that region.
	defaultEndpoints map[string]map[string]endpoints.ResolvedEndpoint
}

func newAWSResolver(region string, services []typesaws.ServiceEndpoint) *awsResolver {
	resolver := &awsResolver{
		region:           region,
		services:         make(map[string]typesaws.ServiceEndpoint),
		defaultEndpoints: defaultEndpoints(),
	}
	for _, service := range services {
		service := service
		resolver.services[resolverKey(service.Name)] = service
	}
	return resolver
}

func (ar *awsResolver) EndpointFor(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
	if s, ok := ar.services[resolverKey(service)]; ok {
		logrus.Debugf("resolved AWS service %s (%s) to %q", service, region, s.URL)
		return endpoints.ResolvedEndpoint{
			URL:           s.URL,
			SigningRegion: ar.region,
		}, nil
	}
	if rv, ok := ar.defaultEndpoints[region]; ok {
		if v, ok := rv[service]; ok {
			return v, nil
		}
	}
	return endpoints.DefaultResolver().EndpointFor(service, region, optFns...)
}

func resolverKey(service string) string {
	return service
}

// this is a list of known default endpoints for specific regions that would
// otherwise require user to set the service overrides.
// it's a map of region => service => resolved endpoint
// this is only used when the user hasn't specified a override for the service in that region.
func defaultEndpoints() map[string]map[string]endpoints.ResolvedEndpoint {
	return map[string]map[string]endpoints.ResolvedEndpoint{
		endpoints.CnNorth1RegionID: {
			"route53": {
				URL:           "https://route53.amazonaws.com.cn",
				SigningRegion: endpoints.CnNorthwest1RegionID,
			},
		},
		endpoints.CnNorthwest1RegionID: {
			"route53": {
				URL:           "https://route53.amazonaws.com.cn",
				SigningRegion: endpoints.CnNorthwest1RegionID,
			},
		},
	}
}
