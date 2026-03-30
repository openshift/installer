/*
Copyright 2018 The OpenShift Authors.

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

package aws

import (
	"context"
	"errors"
	"fmt"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	smithymiddleware "github.com/aws/smithy-go/middleware"

	"github.com/openshift/cloud-credential-operator/pkg/version"
)

//go:generate mockgen -source=./client.go -destination=./mock/client_generated.go -package=mock

// Client is a wrapper object for actual AWS SDK clients to allow for easier testing.
type Client interface {
	//IAM
	CreateAccessKey(context.Context, *iam.CreateAccessKeyInput, ...func(*iam.Options)) (*iam.CreateAccessKeyOutput, error)
	CreateOpenIDConnectProvider(context.Context, *iam.CreateOpenIDConnectProviderInput, ...func(*iam.Options)) (*iam.CreateOpenIDConnectProviderOutput, error)
	CreateRole(context.Context, *iam.CreateRoleInput, ...func(*iam.Options)) (*iam.CreateRoleOutput, error)
	CreateUser(context.Context, *iam.CreateUserInput, ...func(*iam.Options)) (*iam.CreateUserOutput, error)
	DeleteAccessKey(context.Context, *iam.DeleteAccessKeyInput, ...func(*iam.Options)) (*iam.DeleteAccessKeyOutput, error)
	DeleteUser(context.Context, *iam.DeleteUserInput, ...func(*iam.Options)) (*iam.DeleteUserOutput, error)
	DeleteUserPolicy(context.Context, *iam.DeleteUserPolicyInput, ...func(*iam.Options)) (*iam.DeleteUserPolicyOutput, error)
	GetOpenIDConnectProvider(context.Context, *iam.GetOpenIDConnectProviderInput, ...func(*iam.Options)) (*iam.GetOpenIDConnectProviderOutput, error)
	GetRole(context.Context, *iam.GetRoleInput, ...func(*iam.Options)) (*iam.GetRoleOutput, error)
	ListRoles(context.Context, *iam.ListRolesInput, ...func(*iam.Options)) (*iam.ListRolesOutput, error)
	DeleteRole(context.Context, *iam.DeleteRoleInput, ...func(*iam.Options)) (*iam.DeleteRoleOutput, error)
	ListRolePolicies(context.Context, *iam.ListRolePoliciesInput, ...func(*iam.Options)) (*iam.ListRolePoliciesOutput, error)
	DeleteRolePolicy(context.Context, *iam.DeleteRolePolicyInput, ...func(*iam.Options)) (*iam.DeleteRolePolicyOutput, error)
	GetUser(context.Context, *iam.GetUserInput, ...func(*iam.Options)) (*iam.GetUserOutput, error)
	GetUserPolicy(context.Context, *iam.GetUserPolicyInput, ...func(*iam.Options)) (*iam.GetUserPolicyOutput, error)
	ListAccessKeys(context.Context, *iam.ListAccessKeysInput, ...func(*iam.Options)) (*iam.ListAccessKeysOutput, error)
	ListOpenIDConnectProviders(context.Context, *iam.ListOpenIDConnectProvidersInput, ...func(*iam.Options)) (*iam.ListOpenIDConnectProvidersOutput, error)
	DeleteOpenIDConnectProvider(context.Context, *iam.DeleteOpenIDConnectProviderInput, ...func(*iam.Options)) (*iam.DeleteOpenIDConnectProviderOutput, error)
	ListUserPolicies(context.Context, *iam.ListUserPoliciesInput, ...func(*iam.Options)) (*iam.ListUserPoliciesOutput, error)
	PutRolePolicy(context.Context, *iam.PutRolePolicyInput, ...func(*iam.Options)) (*iam.PutRolePolicyOutput, error)
	PutUserPolicy(context.Context, *iam.PutUserPolicyInput, ...func(*iam.Options)) (*iam.PutUserPolicyOutput, error)
	SimulatePrincipalPolicy(context.Context, *iam.SimulatePrincipalPolicyInput, ...func(*iam.Options)) (*iam.SimulatePrincipalPolicyOutput, error)
	TagOpenIDConnectProvider(context.Context, *iam.TagOpenIDConnectProviderInput, ...func(*iam.Options)) (*iam.TagOpenIDConnectProviderOutput, error)
	TagUser(context.Context, *iam.TagUserInput, ...func(*iam.Options)) (*iam.TagUserOutput, error)
	UpdateAssumeRolePolicy(context.Context, *iam.UpdateAssumeRolePolicyInput, ...func(*iam.Options)) (*iam.UpdateAssumeRolePolicyOutput, error)

	//S3
	CreateBucket(context.Context, *s3.CreateBucketInput, ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
	PutBucketTagging(context.Context, *s3.PutBucketTaggingInput, ...func(*s3.Options)) (*s3.PutBucketTaggingOutput, error)
	GetBucketTagging(context.Context, *s3.GetBucketTaggingInput, ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error)
	DeleteBucket(context.Context, *s3.DeleteBucketInput, ...func(*s3.Options)) (*s3.DeleteBucketOutput, error)
	PutObject(context.Context, *s3.PutObjectInput, ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	ListObjects(context.Context, *s3.ListObjectsInput, ...func(*s3.Options)) (*s3.ListObjectsOutput, error)
	GetObjectTagging(context.Context, *s3.GetObjectTaggingInput, ...func(*s3.Options)) (*s3.GetObjectTaggingOutput, error)
	DeleteObject(context.Context, *s3.DeleteObjectInput, ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	PutPublicAccessBlock(context.Context, *s3.PutPublicAccessBlockInput, ...func(*s3.Options)) (*s3.PutPublicAccessBlockOutput, error)
	PutBucketPolicy(context.Context, *s3.PutBucketPolicyInput, ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error)

	//CloudFront
	CreateCloudFrontOriginAccessIdentity(context.Context, *cloudfront.CreateCloudFrontOriginAccessIdentityInput, ...func(*cloudfront.Options)) (*cloudfront.CreateCloudFrontOriginAccessIdentityOutput, error)
	DeleteCloudFrontOriginAccessIdentity(context.Context, *cloudfront.DeleteCloudFrontOriginAccessIdentityInput, ...func(*cloudfront.Options)) (*cloudfront.DeleteCloudFrontOriginAccessIdentityOutput, error)
	GetCloudFrontOriginAccessIdentity(context.Context, *cloudfront.GetCloudFrontOriginAccessIdentityInput, ...func(*cloudfront.Options)) (*cloudfront.GetCloudFrontOriginAccessIdentityOutput, error)
	ListCloudFrontOriginAccessIdentities(context.Context, *cloudfront.ListCloudFrontOriginAccessIdentitiesInput, ...func(*cloudfront.Options)) (*cloudfront.ListCloudFrontOriginAccessIdentitiesOutput, error)
	CreateDistributionWithTags(context.Context, *cloudfront.CreateDistributionWithTagsInput, ...func(*cloudfront.Options)) (*cloudfront.CreateDistributionWithTagsOutput, error)
	DeleteDistribution(context.Context, *cloudfront.DeleteDistributionInput, ...func(*cloudfront.Options)) (*cloudfront.DeleteDistributionOutput, error)
	GetDistribution(context.Context, *cloudfront.GetDistributionInput, ...func(*cloudfront.Options)) (*cloudfront.GetDistributionOutput, error)
	UpdateDistribution(context.Context, *cloudfront.UpdateDistributionInput, ...func(*cloudfront.Options)) (*cloudfront.UpdateDistributionOutput, error)
	ListDistributions(context.Context, *cloudfront.ListDistributionsInput, ...func(*cloudfront.Options)) (*cloudfront.ListDistributionsOutput, error)
	ListTagsForResource(context.Context, *cloudfront.ListTagsForResourceInput, ...func(*cloudfront.Options)) (*cloudfront.ListTagsForResourceOutput, error)
}

// ClientParams holds the various optional tunables that can be used to modify the AWS
// client that will be used for API calls.
type ClientParams struct {
	InfraName string
	Region    string
	Endpoint  string
	CABundle  string
}

type awsClient struct {
	iamClient        *iam.Client
	s3Client         *s3.Client
	cloudFrontClient *cloudfront.Client
}

func (c *awsClient) CreateAccessKey(ctx context.Context, input *iam.CreateAccessKeyInput, opts ...func(*iam.Options)) (*iam.CreateAccessKeyOutput, error) {
	return c.iamClient.CreateAccessKey(ctx, input, opts...)
}

func (c *awsClient) CreateRole(ctx context.Context, input *iam.CreateRoleInput, opts ...func(*iam.Options)) (*iam.CreateRoleOutput, error) {
	return c.iamClient.CreateRole(ctx, input, opts...)
}

func (c *awsClient) CreateUser(ctx context.Context, input *iam.CreateUserInput, opts ...func(*iam.Options)) (*iam.CreateUserOutput, error) {
	return c.iamClient.CreateUser(ctx, input, opts...)
}

func (c *awsClient) DeleteAccessKey(ctx context.Context, input *iam.DeleteAccessKeyInput, opts ...func(*iam.Options)) (*iam.DeleteAccessKeyOutput, error) {
	return c.iamClient.DeleteAccessKey(ctx, input, opts...)
}

func (c *awsClient) DeleteUser(ctx context.Context, input *iam.DeleteUserInput, opts ...func(*iam.Options)) (*iam.DeleteUserOutput, error) {
	return c.iamClient.DeleteUser(ctx, input, opts...)
}

func (c *awsClient) DeleteUserPolicy(ctx context.Context, input *iam.DeleteUserPolicyInput, opts ...func(*iam.Options)) (*iam.DeleteUserPolicyOutput, error) {
	return c.iamClient.DeleteUserPolicy(ctx, input, opts...)
}
func (c *awsClient) GetUser(ctx context.Context, input *iam.GetUserInput, opts ...func(*iam.Options)) (*iam.GetUserOutput, error) {
	return c.iamClient.GetUser(ctx, input, opts...)
}

func (c *awsClient) ListAccessKeys(ctx context.Context, input *iam.ListAccessKeysInput, opts ...func(*iam.Options)) (*iam.ListAccessKeysOutput, error) {
	return c.iamClient.ListAccessKeys(ctx, input, opts...)
}

func (c *awsClient) ListUserPolicies(ctx context.Context, input *iam.ListUserPoliciesInput, opts ...func(*iam.Options)) (*iam.ListUserPoliciesOutput, error) {
	return c.iamClient.ListUserPolicies(ctx, input, opts...)
}

func (c *awsClient) PutRolePolicy(ctx context.Context, input *iam.PutRolePolicyInput, opts ...func(*iam.Options)) (*iam.PutRolePolicyOutput, error) {
	return c.iamClient.PutRolePolicy(ctx, input, opts...)
}

func (c *awsClient) PutUserPolicy(ctx context.Context, input *iam.PutUserPolicyInput, opts ...func(*iam.Options)) (*iam.PutUserPolicyOutput, error) {
	return c.iamClient.PutUserPolicy(ctx, input, opts...)
}

func (c *awsClient) GetRole(ctx context.Context, input *iam.GetRoleInput, opts ...func(*iam.Options)) (*iam.GetRoleOutput, error) {
	return c.iamClient.GetRole(ctx, input, opts...)
}

func (c *awsClient) GetUserPolicy(ctx context.Context, input *iam.GetUserPolicyInput, opts ...func(*iam.Options)) (*iam.GetUserPolicyOutput, error) {
	return c.iamClient.GetUserPolicy(ctx, input, opts...)
}

func (c *awsClient) SimulatePrincipalPolicy(ctx context.Context, input *iam.SimulatePrincipalPolicyInput, opts ...func(*iam.Options)) (*iam.SimulatePrincipalPolicyOutput, error) {
	return c.iamClient.SimulatePrincipalPolicy(ctx, input, opts...)
}

func (c *awsClient) TagUser(ctx context.Context, input *iam.TagUserInput, opts ...func(*iam.Options)) (*iam.TagUserOutput, error) {
	return c.iamClient.TagUser(ctx, input, opts...)
}

func (c *awsClient) ListOpenIDConnectProviders(ctx context.Context, input *iam.ListOpenIDConnectProvidersInput, opts ...func(*iam.Options)) (*iam.ListOpenIDConnectProvidersOutput, error) {
	return c.iamClient.ListOpenIDConnectProviders(ctx, input, opts...)
}

func (c *awsClient) CreateOpenIDConnectProvider(ctx context.Context, input *iam.CreateOpenIDConnectProviderInput, opts ...func(*iam.Options)) (*iam.CreateOpenIDConnectProviderOutput, error) {
	return c.iamClient.CreateOpenIDConnectProvider(ctx, input, opts...)
}

func (c *awsClient) TagOpenIDConnectProvider(ctx context.Context, input *iam.TagOpenIDConnectProviderInput, opts ...func(*iam.Options)) (*iam.TagOpenIDConnectProviderOutput, error) {
	return c.iamClient.TagOpenIDConnectProvider(ctx, input, opts...)
}

func (c *awsClient) UpdateAssumeRolePolicy(ctx context.Context, input *iam.UpdateAssumeRolePolicyInput, opts ...func(*iam.Options)) (*iam.UpdateAssumeRolePolicyOutput, error) {
	return c.iamClient.UpdateAssumeRolePolicy(ctx, input, opts...)
}

func (c *awsClient) GetOpenIDConnectProvider(ctx context.Context, input *iam.GetOpenIDConnectProviderInput, opts ...func(*iam.Options)) (*iam.GetOpenIDConnectProviderOutput, error) {
	return c.iamClient.GetOpenIDConnectProvider(ctx, input, opts...)
}

func (c *awsClient) DeleteOpenIDConnectProvider(ctx context.Context, input *iam.DeleteOpenIDConnectProviderInput, opts ...func(*iam.Options)) (*iam.DeleteOpenIDConnectProviderOutput, error) {
	return c.iamClient.DeleteOpenIDConnectProvider(ctx, input, opts...)
}

func (c *awsClient) ListRoles(ctx context.Context, input *iam.ListRolesInput, opts ...func(*iam.Options)) (*iam.ListRolesOutput, error) {
	return c.iamClient.ListRoles(ctx, input, opts...)
}

func (c *awsClient) DeleteRole(ctx context.Context, input *iam.DeleteRoleInput, opts ...func(*iam.Options)) (*iam.DeleteRoleOutput, error) {
	return c.iamClient.DeleteRole(ctx, input, opts...)
}

func (c *awsClient) ListRolePolicies(ctx context.Context, input *iam.ListRolePoliciesInput, opts ...func(*iam.Options)) (*iam.ListRolePoliciesOutput, error) {
	return c.iamClient.ListRolePolicies(ctx, input, opts...)
}

func (c *awsClient) DeleteRolePolicy(ctx context.Context, input *iam.DeleteRolePolicyInput, opts ...func(*iam.Options)) (*iam.DeleteRolePolicyOutput, error) {
	return c.iamClient.DeleteRolePolicy(ctx, input, opts...)
}

func (c *awsClient) CreateBucket(ctx context.Context, input *s3.CreateBucketInput, opts ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return c.s3Client.CreateBucket(ctx, input, opts...)
}

func (c *awsClient) PutBucketTagging(ctx context.Context, input *s3.PutBucketTaggingInput, opts ...func(*s3.Options)) (*s3.PutBucketTaggingOutput, error) {
	return c.s3Client.PutBucketTagging(ctx, input, opts...)
}

func (c *awsClient) GetBucketTagging(ctx context.Context, input *s3.GetBucketTaggingInput, opts ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
	return c.s3Client.GetBucketTagging(ctx, input, opts...)
}

func (c *awsClient) DeleteBucket(ctx context.Context, input *s3.DeleteBucketInput, opts ...func(*s3.Options)) (*s3.DeleteBucketOutput, error) {
	return c.s3Client.DeleteBucket(ctx, input, opts...)
}

func (c *awsClient) PutObject(ctx context.Context, input *s3.PutObjectInput, opts ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return c.s3Client.PutObject(ctx, input, opts...)
}

func (c *awsClient) ListObjects(ctx context.Context, input *s3.ListObjectsInput, opts ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
	return c.s3Client.ListObjects(ctx, input, opts...)
}

func (c *awsClient) GetObjectTagging(ctx context.Context, input *s3.GetObjectTaggingInput, opts ...func(*s3.Options)) (*s3.GetObjectTaggingOutput, error) {
	return c.s3Client.GetObjectTagging(ctx, input, opts...)
}

func (c *awsClient) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	return c.s3Client.DeleteObject(ctx, input, opts...)
}

func (c *awsClient) PutPublicAccessBlock(ctx context.Context, input *s3.PutPublicAccessBlockInput, opts ...func(*s3.Options)) (*s3.PutPublicAccessBlockOutput, error) {
	return c.s3Client.PutPublicAccessBlock(ctx, input, opts...)
}

func (c *awsClient) PutBucketPolicy(ctx context.Context, input *s3.PutBucketPolicyInput, opts ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error) {
	return c.s3Client.PutBucketPolicy(ctx, input, opts...)
}

func (c *awsClient) CreateCloudFrontOriginAccessIdentity(ctx context.Context, input *cloudfront.CreateCloudFrontOriginAccessIdentityInput, opts ...func(*cloudfront.Options)) (*cloudfront.CreateCloudFrontOriginAccessIdentityOutput, error) {
	return c.cloudFrontClient.CreateCloudFrontOriginAccessIdentity(ctx, input, opts...)
}

func (c *awsClient) DeleteCloudFrontOriginAccessIdentity(ctx context.Context, input *cloudfront.DeleteCloudFrontOriginAccessIdentityInput, opts ...func(*cloudfront.Options)) (*cloudfront.DeleteCloudFrontOriginAccessIdentityOutput, error) {
	return c.cloudFrontClient.DeleteCloudFrontOriginAccessIdentity(ctx, input, opts...)
}

func (c *awsClient) GetCloudFrontOriginAccessIdentity(ctx context.Context, input *cloudfront.GetCloudFrontOriginAccessIdentityInput, opts ...func(*cloudfront.Options)) (*cloudfront.GetCloudFrontOriginAccessIdentityOutput, error) {
	return c.cloudFrontClient.GetCloudFrontOriginAccessIdentity(ctx, input, opts...)
}

func (c *awsClient) ListCloudFrontOriginAccessIdentities(ctx context.Context, input *cloudfront.ListCloudFrontOriginAccessIdentitiesInput, opts ...func(*cloudfront.Options)) (*cloudfront.ListCloudFrontOriginAccessIdentitiesOutput, error) {
	return c.cloudFrontClient.ListCloudFrontOriginAccessIdentities(ctx, input, opts...)
}
func (c *awsClient) CreateDistributionWithTags(ctx context.Context, input *cloudfront.CreateDistributionWithTagsInput, opts ...func(*cloudfront.Options)) (*cloudfront.CreateDistributionWithTagsOutput, error) {
	return c.cloudFrontClient.CreateDistributionWithTags(ctx, input, opts...)
}

func (c *awsClient) DeleteDistribution(ctx context.Context, input *cloudfront.DeleteDistributionInput, opts ...func(*cloudfront.Options)) (*cloudfront.DeleteDistributionOutput, error) {
	return c.cloudFrontClient.DeleteDistribution(ctx, input, opts...)
}

func (c *awsClient) GetDistribution(ctx context.Context, input *cloudfront.GetDistributionInput, opts ...func(*cloudfront.Options)) (*cloudfront.GetDistributionOutput, error) {
	return c.cloudFrontClient.GetDistribution(ctx, input, opts...)
}

func (c *awsClient) UpdateDistribution(ctx context.Context, input *cloudfront.UpdateDistributionInput, opts ...func(*cloudfront.Options)) (*cloudfront.UpdateDistributionOutput, error) {
	return c.cloudFrontClient.UpdateDistribution(ctx, input, opts...)
}

func (c *awsClient) ListDistributions(ctx context.Context, input *cloudfront.ListDistributionsInput, opts ...func(*cloudfront.Options)) (*cloudfront.ListDistributionsOutput, error) {
	return c.cloudFrontClient.ListDistributions(ctx, input, opts...)
}

func (c *awsClient) ListTagsForResource(ctx context.Context, input *cloudfront.ListTagsForResourceInput, opts ...func(*cloudfront.Options)) (*cloudfront.ListTagsForResourceOutput, error) {
	return c.cloudFrontClient.ListTagsForResource(ctx, input, opts...)
}

// NewClient creates our client wrapper object for the actual AWS clients we use.
func NewClient(accessKeyID, secretAccessKey []byte, params *ClientParams) (Client, error) {
	awsOpts := []func(o *config.LoadOptions) error{}
	agentText := "defaultAgent"
	var endpoint string

	if params != nil {
		if params.Region != "" {
			awsOpts = append(awsOpts, config.WithRegion(params.Region))
		}

		if params.Endpoint != "" {
			endpoint = params.Endpoint
		}

		if params.CABundle != "" {
			awsOpts = append(awsOpts, config.WithCustomCABundle(strings.NewReader(params.CABundle)))
		}

		if params.InfraName != "" {
			agentText = params.InfraName
		}
	}

	awsOpts = append(awsOpts,
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				string(accessKeyID), string(secretAccessKey), ""),
		),
	)

	agentKey := fmt.Sprintf("%s/%s", "openshift.io cloud-credential-operator", version.Get().String())
	awsOpts = append(awsOpts,
		config.WithAPIOptions([]func(*smithymiddleware.Stack) error{
			middleware.AddUserAgentKeyValue(agentKey, agentText),
		}),
	)

	cfg, err := config.LoadDefaultConfig(context.TODO(), awsOpts...)
	if err != nil {
		return nil, err
	}

	client, err := NewClientFromConfig(cfg, endpoint)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// NewClientFromConfig will return a basic Client using only the provided aws.Config
func NewClientFromConfig(cfg awssdk.Config, endpoint string) (Client, error) {
	iamOpts := []func(o *iam.Options){}
	s3Opts := []func(o *s3.Options){}
	cfOpts := []func(o *cloudfront.Options){}

	if endpoint != "" {
		iamOpts = append(iamOpts, func(o *iam.Options) {
			o.BaseEndpoint = &endpoint
		})
	}

	return &awsClient{
		iamClient:        iam.NewFromConfig(cfg, iamOpts...),
		s3Client:         s3.NewFromConfig(cfg, s3Opts...),
		cloudFrontClient: cloudfront.NewFromConfig(cfg, cfOpts...),
	}, nil
}

func ErrCodeEquals(err error, code string) bool {
	var awsErr smithy.APIError
	return err != nil && errors.As(err, &awsErr) && awsErr.ErrorCode() == code
}

func NewAPIError(code, message string) error {
	return &smithy.GenericAPIError{Code: code, Message: message}
}
