package config

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/aws-sdk-go-base/v2/internal/expand"
)

type Config struct {
	AccessKey                      string
	APNInfo                        *APNInfo
	AssumeRole                     *AssumeRole
	CallerDocumentationURL         string
	CallerName                     string
	CustomCABundle                 string
	EC2MetadataServiceEndpoint     string
	EC2MetadataServiceEndpointMode string
	HTTPProxy                      string
	IamEndpoint                    string
	Insecure                       bool
	MaxRetries                     int
	Profile                        string
	Region                         string
	SecretKey                      string
	SharedCredentialsFiles         []string
	SharedConfigFiles              []string
	SkipCredsValidation            bool
	SkipEC2MetadataApiCheck        bool
	SkipRequestingAccountId        bool
	StsEndpoint                    string
	StsRegion                      string
	Token                          string
	UseDualStackEndpoint           bool
	UseFIPSEndpoint                bool
	UserAgent                      UserAgentProducts
}

type AssumeRole struct {
	RoleARN           string
	Duration          time.Duration
	ExternalID        string
	Policy            string
	PolicyARNs        []string
	SessionName       string
	Tags              map[string]string
	TransitiveTagKeys []string
}

func (c Config) CustomCABundleReader() (*bytes.Reader, error) {
	if c.CustomCABundle == "" {
		return nil, nil
	}
	bundleFile, err := expand.FilePath(c.CustomCABundle)
	if err != nil {
		return nil, fmt.Errorf("expanding custom CA bundle: %w", err)
	}
	bundle, err := os.ReadFile(bundleFile)
	if err != nil {
		return nil, fmt.Errorf("reading custom CA bundle: %w", err)
	}
	return bytes.NewReader(bundle), nil
}

func (c Config) ResolveSharedConfigFiles() ([]string, error) {
	v, err := expand.FilePaths(c.SharedConfigFiles)
	if err != nil {
		return []string{}, fmt.Errorf("expanding shared config files: %w", err)
	}
	return v, nil
}

func (c Config) ResolveSharedCredentialsFiles() ([]string, error) {
	v, err := expand.FilePaths(c.SharedCredentialsFiles)
	if err != nil {
		return []string{}, fmt.Errorf("expanding shared credentials files: %w", err)
	}
	return v, nil
}
