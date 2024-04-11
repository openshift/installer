/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package identity provides the AWSPrincipalTypeProvider interface and its implementations.
package identity

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	corev1 "k8s.io/api/core/v1"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
)

// AWSPrincipalTypeProvider defines the interface for AWS Principal Type Provider.
type AWSPrincipalTypeProvider interface {
	credentials.Provider
	// Hash returns a unique hash of the data forming the credentials
	// for this Principal
	Hash() (string, error)
	Name() string
}

// NewAWSStaticPrincipalTypeProvider will create a new AWSStaticPrincipalTypeProvider from a given AWSClusterStaticIdentity.
func NewAWSStaticPrincipalTypeProvider(identity *infrav1.AWSClusterStaticIdentity, secret *corev1.Secret) *AWSStaticPrincipalTypeProvider {
	accessKeyID := string(secret.Data["AccessKeyID"])
	secretAccessKey := string(secret.Data["SecretAccessKey"])
	sessionToken := string(secret.Data["SessionToken"])

	return &AWSStaticPrincipalTypeProvider{
		Principal:       identity,
		credentials:     credentials.NewStaticCredentials(accessKeyID, secretAccessKey, sessionToken),
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		SessionToken:    sessionToken,
	}
}

// GetAssumeRoleCredentials will return the Credentials of a given AWSRolePrincipalTypeProvider.
func GetAssumeRoleCredentials(roleIdentityProvider *AWSRolePrincipalTypeProvider, awsConfig *aws.Config) *credentials.Credentials {
	sess := session.Must(session.NewSession(awsConfig))

	creds := stscreds.NewCredentials(sess, roleIdentityProvider.Principal.Spec.RoleArn, func(p *stscreds.AssumeRoleProvider) {
		if roleIdentityProvider.Principal.Spec.ExternalID != "" {
			p.ExternalID = aws.String(roleIdentityProvider.Principal.Spec.ExternalID)
		}
		p.RoleSessionName = roleIdentityProvider.Principal.Spec.SessionName
		if roleIdentityProvider.Principal.Spec.InlinePolicy != "" {
			p.Policy = aws.String(roleIdentityProvider.Principal.Spec.InlinePolicy)
		}
		p.Duration = time.Duration(roleIdentityProvider.Principal.Spec.DurationSeconds) * time.Second
		// For testing
		if roleIdentityProvider.stsClient != nil {
			p.Client = roleIdentityProvider.stsClient
		}
	})
	return creds
}

// NewAWSRolePrincipalTypeProvider will create a new AWSRolePrincipalTypeProvider from an AWSClusterRoleIdentity.
func NewAWSRolePrincipalTypeProvider(identity *infrav1.AWSClusterRoleIdentity, sourceProvider AWSPrincipalTypeProvider, region string, log logger.Wrapper) *AWSRolePrincipalTypeProvider {
	return &AWSRolePrincipalTypeProvider{
		credentials:    nil,
		stsClient:      nil,
		region:         region,
		Principal:      identity,
		sourceProvider: sourceProvider,
		log:            log.WithName("AWSRolePrincipalTypeProvider"),
	}
}

// AWSStaticPrincipalTypeProvider defines the specs for a static AWSPrincipalTypeProvider.
type AWSStaticPrincipalTypeProvider struct {
	Principal   *infrav1.AWSClusterStaticIdentity
	credentials *credentials.Credentials
	// these are for tests :/
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

// Hash returns the byte encoded AWSStaticPrincipalTypeProvider.
func (p *AWSStaticPrincipalTypeProvider) Hash() (string, error) {
	var roleIdentityValue bytes.Buffer
	err := gob.NewEncoder(&roleIdentityValue).Encode(p)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	return string(hash.Sum(roleIdentityValue.Bytes())), nil
}

// Retrieve returns the credential values for the AWSStaticPrincipalTypeProvider.
func (p *AWSStaticPrincipalTypeProvider) Retrieve() (credentials.Value, error) {
	return p.credentials.Get()
}

// Name returns the name of the AWSStaticPrincipalTypeProvider.
func (p *AWSStaticPrincipalTypeProvider) Name() string {
	return p.Principal.Name
}

// IsExpired checks the expiration state of the AWSStaticPrincipalTypeProvider.
func (p *AWSStaticPrincipalTypeProvider) IsExpired() bool {
	return p.credentials.IsExpired()
}

// AWSRolePrincipalTypeProvider defines the specs for a AWSPrincipalTypeProvider with a role.
type AWSRolePrincipalTypeProvider struct {
	Principal      *infrav1.AWSClusterRoleIdentity
	credentials    *credentials.Credentials
	region         string
	sourceProvider AWSPrincipalTypeProvider
	log            logger.Wrapper
	stsClient      stsiface.STSAPI
}

// Hash returns the byte encoded AWSRolePrincipalTypeProvider.
func (p *AWSRolePrincipalTypeProvider) Hash() (string, error) {
	var roleIdentityValue bytes.Buffer
	err := gob.NewEncoder(&roleIdentityValue).Encode(p)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	return string(hash.Sum(roleIdentityValue.Bytes())), nil
}

// Name returns the name of the AWSRolePrincipalTypeProvider.
func (p *AWSRolePrincipalTypeProvider) Name() string {
	return p.Principal.Name
}

// Retrieve returns the credential values for the AWSRolePrincipalTypeProvider.
func (p *AWSRolePrincipalTypeProvider) Retrieve() (credentials.Value, error) {
	if p.credentials == nil || p.IsExpired() {
		awsConfig := aws.NewConfig().WithRegion(p.region)
		if p.sourceProvider != nil {
			sourceCreds, err := p.sourceProvider.Retrieve()
			if err != nil {
				return credentials.Value{}, err
			}
			awsConfig = awsConfig.WithCredentials(credentials.NewStaticCredentialsFromCreds(sourceCreds))
		}

		creds := GetAssumeRoleCredentials(p, awsConfig)
		// Update credentials
		p.credentials = creds
	}
	return p.credentials.Get()
}

// IsExpired checks the expiration state of the AWSRolePrincipalTypeProvider.
func (p *AWSRolePrincipalTypeProvider) IsExpired() bool {
	return p.credentials.IsExpired()
}
