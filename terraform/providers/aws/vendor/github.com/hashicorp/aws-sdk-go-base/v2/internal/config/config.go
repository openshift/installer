// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
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
	HTTPClient                     *http.Client
	HTTPProxy                      string
	IamEndpoint                    string
	Insecure                       bool
	MaxRetries                     int
	Profile                        string
	Region                         string
	RetryMode                      aws.RetryMode
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

// HTTPTransportOptions returns functional options that configures an http.Transport.
// The returned options function is called on both AWS SDKv1 and v2 default HTTP clients.
func (c Config) HTTPTransportOptions() (func(*http.Transport), error) {
	var err error
	var proxyUrl *url.URL
	if c.HTTPProxy != "" {
		proxyUrl, err = url.Parse(c.HTTPProxy)
		if err != nil {
			return nil, fmt.Errorf("error parsing HTTP proxy URL: %w", err)
		}
	}

	opts := func(tr *http.Transport) {
		tr.MaxIdleConnsPerHost = awshttp.DefaultHTTPTransportMaxIdleConnsPerHost

		tlsConfig := tr.TLSClientConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
			}
			tr.TLSClientConfig = tlsConfig
		}

		if c.Insecure {
			tr.TLSClientConfig.InsecureSkipVerify = true
		}

		if proxyUrl != nil {
			tr.Proxy = http.ProxyURL(proxyUrl)
		}
	}

	return opts, nil
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

type AssumeRoleWithWebIdentity struct {
	RoleARN              string
	Duration             time.Duration
	Policy               string
	PolicyARNs           []string
	SessionName          string
	WebIdentityToken     string
	WebIdentityTokenFile string
}

func (c AssumeRoleWithWebIdentity) resolveWebIdentityTokenFile() (string, error) {
	v, err := expand.FilePath(c.WebIdentityTokenFile)
	if err != nil {
		return "", fmt.Errorf("expanding web identity token file: %w", err)
	}
	return v, nil
}

func (c AssumeRoleWithWebIdentity) HasValidTokenSource() bool {
	return c.WebIdentityToken != "" || c.WebIdentityTokenFile != ""
}

// Implements `stscreds.IdentityTokenRetriever`
func (c AssumeRoleWithWebIdentity) GetIdentityToken() ([]byte, error) {
	if c.WebIdentityToken != "" {
		return []byte(c.WebIdentityToken), nil
	}
	webIdentityTokenFile, err := c.resolveWebIdentityTokenFile()
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(webIdentityTokenFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read file at %s: %w", webIdentityTokenFile, err)
	}

	return b, nil
}
