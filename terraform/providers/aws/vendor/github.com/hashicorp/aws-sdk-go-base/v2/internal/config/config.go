package config

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/hashicorp/aws-sdk-go-base/v2/internal/expand"
)

type Config struct {
	AccessKey                      string
	APNInfo                        *APNInfo
	AssumeRole                     *AssumeRole
	AssumeRoleWithWebIdentity      *AssumeRoleWithWebIdentity
	CallerDocumentationURL         string
	CallerName                     string
	CustomCABundle                 string
	EC2MetadataServiceEnableState  imds.ClientEnableState
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
	SkipRequestingAccountId        bool
	StsEndpoint                    string
	StsRegion                      string
	SuppressDebugLog               bool
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
	SourceIdentity    string
	Tags              map[string]string
	TransitiveTagKeys []string
}

type AssumeRoleWithWebIdentity struct {
	RoleARN              string
	Duration             time.Duration
	Policy               string
	PolicyARNs           []string
	SessionName          string
	WebIdentityToken     string
	WebIdentityTokenFile string
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

func (c AssumeRoleWithWebIdentity) ResolveWebIdentityTokenFile() (string, error) {
	v, err := expand.FilePath(c.WebIdentityTokenFile)
	if err != nil {
		return "", fmt.Errorf("expanding web identity token file: %w", err)
	}
	return v, nil
}

func (c AssumeRoleWithWebIdentity) GetIdentityToken() ([]byte, error) {
	if c.WebIdentityToken != "" {
		return []byte(c.WebIdentityToken), nil
	}
	webIdentityTokenFile, err := c.ResolveWebIdentityTokenFile()
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(webIdentityTokenFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read file at %s: %w", webIdentityTokenFile, err)
	}

	return b, nil
}
