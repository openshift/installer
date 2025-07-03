/*
Copyright 2025 The Kubernetes Authors.

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

// Package identityv2 provides the AWSPrincipalTypeProvider interface and its implementations.
package identityv2

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/gob"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	corev1 "k8s.io/api/core/v1"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	awsmetricsv2 "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/metricsv2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
)

// AWSPrincipalTypeProvider defines the interface for AWS Principal Type Provider.
type AWSPrincipalTypeProvider interface {
	aws.CredentialsProvider
	// Hash returns a unique hash of the data forming the V2 credentials
	// for this Principal
	Hash() (string, error)
	Name() string
}

// NewAWSStaticPrincipalTypeProvider will create a new AWSStaticPrincipalTypeProvider from a given AWSClusterStaticIdentity.
func NewAWSStaticPrincipalTypeProvider(identity *infrav1.AWSClusterStaticIdentity, secret *corev1.Secret) *AWSStaticPrincipalTypeProvider {
	accessKeyID := string(secret.Data["AccessKeyID"])
	secretAccessKey := string(secret.Data["SecretAccessKey"])
	sessionToken := string(secret.Data["SessionToken"])

	credProvider := credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, sessionToken)
	credCache := aws.NewCredentialsCache(credProvider)
	return &AWSStaticPrincipalTypeProvider{
		Principal:       identity,
		credentials:     credCache,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		SessionToken:    sessionToken,
	}
}

// GetAssumeRoleCredentialsCache will return the CredentialsCache of a given AWSRolePrincipalTypeProvider.
func GetAssumeRoleCredentialsCache(ctx context.Context, roleIdentityProvider *AWSRolePrincipalTypeProvider, optFns []func(*config.LoadOptions) error) (*aws.CredentialsCache, error) {
	cfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		return nil, err
	}

	stsOpts := sts.WithAPIOptions(
		awsmetricsv2.WithMiddlewares("identity provider", roleIdentityProvider.Principal),
		awsmetricsv2.WithCAPAUserAgentMiddleware())
	stsClient := sts.NewFromConfig(cfg, stsOpts)
	credsProvider := stscreds.NewAssumeRoleProvider(stsClient, roleIdentityProvider.Principal.Spec.RoleArn, func(o *stscreds.AssumeRoleOptions) {
		if roleIdentityProvider.Principal.Spec.ExternalID != "" {
			o.ExternalID = aws.String(roleIdentityProvider.Principal.Spec.ExternalID)
		}
		o.RoleSessionName = roleIdentityProvider.Principal.Spec.SessionName
		if roleIdentityProvider.Principal.Spec.InlinePolicy != "" {
			o.Policy = aws.String(roleIdentityProvider.Principal.Spec.InlinePolicy)
		}
		o.Duration = time.Duration(roleIdentityProvider.Principal.Spec.DurationSeconds) * time.Second
		// For testing
		if roleIdentityProvider.stsClient != nil {
			o.Client = roleIdentityProvider.stsClient
		}
	})

	return aws.NewCredentialsCache(credsProvider), nil
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
	credentials *aws.CredentialsCache
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
func (p *AWSStaticPrincipalTypeProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return p.credentials.Retrieve(ctx)
}

// Name returns the name of the AWSStaticPrincipalTypeProvider.
func (p *AWSStaticPrincipalTypeProvider) Name() string {
	return p.Principal.Name
}

// AWSRolePrincipalTypeProvider defines the specs for a AWSPrincipalTypeProvider with a role.
type AWSRolePrincipalTypeProvider struct {
	Principal      *infrav1.AWSClusterRoleIdentity
	credentials    *aws.CredentialsCache
	region         string
	sourceProvider AWSPrincipalTypeProvider
	log            logger.Wrapper
	stsClient      stscreds.AssumeRoleAPIClient
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
func (p *AWSRolePrincipalTypeProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	if p.credentials == nil {
		optFns := []func(*config.LoadOptions) error{config.WithRegion(p.region)}
		if p.sourceProvider != nil {
			sourceCreds, err := p.sourceProvider.Retrieve(ctx)
			if err != nil {
				return aws.Credentials{}, err
			}
			optFns = append(optFns, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(sourceCreds.AccessKeyID, sourceCreds.SecretAccessKey, sourceCreds.SessionToken)))
		}

		creds, err := GetAssumeRoleCredentialsCache(ctx, p, optFns)
		if err != nil {
			return aws.Credentials{}, err
		}
		// Update credentials
		p.credentials = creds
	}
	return p.credentials.Retrieve(ctx)
}
