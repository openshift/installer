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
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/openshift/cloud-credential-operator/pkg/version"
)

//go:generate mockgen -source=./client.go -destination=./mock/client_generated.go -package=mock

// Client is a wrapper object for actual AWS SDK clients to allow for easier testing.
type Client interface {
	//IAM
	CreateAccessKey(*iam.CreateAccessKeyInput) (*iam.CreateAccessKeyOutput, error)
	CreateOpenIDConnectProvider(*iam.CreateOpenIDConnectProviderInput) (*iam.CreateOpenIDConnectProviderOutput, error)
	CreateRole(*iam.CreateRoleInput) (*iam.CreateRoleOutput, error)
	CreateUser(*iam.CreateUserInput) (*iam.CreateUserOutput, error)
	DeleteAccessKey(*iam.DeleteAccessKeyInput) (*iam.DeleteAccessKeyOutput, error)
	DeleteUser(*iam.DeleteUserInput) (*iam.DeleteUserOutput, error)
	DeleteUserPolicy(*iam.DeleteUserPolicyInput) (*iam.DeleteUserPolicyOutput, error)
	GetOpenIDConnectProvider(input *iam.GetOpenIDConnectProviderInput) (*iam.GetOpenIDConnectProviderOutput, error)
	GetRole(input *iam.GetRoleInput) (*iam.GetRoleOutput, error)
	ListRoles(input *iam.ListRolesInput) (*iam.ListRolesOutput, error)
	DeleteRole(input *iam.DeleteRoleInput) (*iam.DeleteRoleOutput, error)
	ListRolePolicies(input *iam.ListRolePoliciesInput) (*iam.ListRolePoliciesOutput, error)
	DeleteRolePolicy(input *iam.DeleteRolePolicyInput) (*iam.DeleteRolePolicyOutput, error)
	GetUser(*iam.GetUserInput) (*iam.GetUserOutput, error)
	GetUserPolicy(*iam.GetUserPolicyInput) (*iam.GetUserPolicyOutput, error)
	ListAccessKeys(*iam.ListAccessKeysInput) (*iam.ListAccessKeysOutput, error)
	ListOpenIDConnectProviders(*iam.ListOpenIDConnectProvidersInput) (*iam.ListOpenIDConnectProvidersOutput, error)
	DeleteOpenIDConnectProvider(input *iam.DeleteOpenIDConnectProviderInput) (*iam.DeleteOpenIDConnectProviderOutput, error)
	ListUserPolicies(*iam.ListUserPoliciesInput) (*iam.ListUserPoliciesOutput, error)
	PutRolePolicy(*iam.PutRolePolicyInput) (*iam.PutRolePolicyOutput, error)
	PutUserPolicy(*iam.PutUserPolicyInput) (*iam.PutUserPolicyOutput, error)
	SimulatePrincipalPolicy(*iam.SimulatePrincipalPolicyInput) (*iam.SimulatePolicyResponse, error)
	SimulatePrincipalPolicyPages(*iam.SimulatePrincipalPolicyInput, func(*iam.SimulatePolicyResponse, bool) bool) error
	TagOpenIDConnectProvider(*iam.TagOpenIDConnectProviderInput) (*iam.TagOpenIDConnectProviderOutput, error)
	TagUser(*iam.TagUserInput) (*iam.TagUserOutput, error)
	UpdateAssumeRolePolicy(*iam.UpdateAssumeRolePolicyInput) (*iam.UpdateAssumeRolePolicyOutput, error)

	//S3
	CreateBucket(*s3.CreateBucketInput) (*s3.CreateBucketOutput, error)
	PutBucketTagging(*s3.PutBucketTaggingInput) (*s3.PutBucketTaggingOutput, error)
	GetBucketTagging(input *s3.GetBucketTaggingInput) (*s3.GetBucketTaggingOutput, error)
	DeleteBucket(input *s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error)
	PutObject(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
	ListObjects(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error)
	GetObjectTagging(input *s3.GetObjectTaggingInput) (*s3.GetObjectTaggingOutput, error)
	DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
	PutPublicAccessBlock(input *s3.PutPublicAccessBlockInput) (*s3.PutPublicAccessBlockOutput, error)
	PutBucketPolicy(input *s3.PutBucketPolicyInput) (*s3.PutBucketPolicyOutput, error)

	//CloudFront
	CreateCloudFrontOriginAccessIdentity(input *cloudfront.CreateCloudFrontOriginAccessIdentityInput) (*cloudfront.CreateCloudFrontOriginAccessIdentityOutput, error)
	DeleteCloudFrontOriginAccessIdentity(input *cloudfront.DeleteCloudFrontOriginAccessIdentityInput) (*cloudfront.DeleteCloudFrontOriginAccessIdentityOutput, error)
	GetCloudFrontOriginAccessIdentity(input *cloudfront.GetCloudFrontOriginAccessIdentityInput) (*cloudfront.GetCloudFrontOriginAccessIdentityOutput, error)
	ListCloudFrontOriginAccessIdentities(input *cloudfront.ListCloudFrontOriginAccessIdentitiesInput) (*cloudfront.ListCloudFrontOriginAccessIdentitiesOutput, error)
	CreateCloudFrontDistributionWithTags(input *cloudfront.CreateDistributionWithTagsInput) (*cloudfront.CreateDistributionWithTagsOutput, error)
	DeleteCloudFrontDistribution(input *cloudfront.DeleteDistributionInput) (*cloudfront.DeleteDistributionOutput, error)
	GetCloudFrontDistribution(input *cloudfront.GetDistributionInput) (*cloudfront.GetDistributionOutput, error)
	UpdateCloudFrontDistribution(input *cloudfront.UpdateDistributionInput) (*cloudfront.UpdateDistributionOutput, error)
	ListCloudFrontDistributions(input *cloudfront.ListDistributionsInput) (*cloudfront.ListDistributionsOutput, error)
	ListTagsForCloudFrontResource(input *cloudfront.ListTagsForResourceInput) (*cloudfront.ListTagsForResourceOutput, error)
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
	iamClient        iamiface.IAMAPI
	s3Client         s3iface.S3API
	cloudFrontClient cloudfrontiface.CloudFrontAPI
}

func (c *awsClient) CreateAccessKey(input *iam.CreateAccessKeyInput) (*iam.CreateAccessKeyOutput, error) {
	return c.iamClient.CreateAccessKey(input)
}

func (c *awsClient) CreateRole(input *iam.CreateRoleInput) (*iam.CreateRoleOutput, error) {
	return c.iamClient.CreateRole(input)
}

func (c *awsClient) CreateUser(input *iam.CreateUserInput) (*iam.CreateUserOutput, error) {
	return c.iamClient.CreateUser(input)
}

func (c *awsClient) DeleteAccessKey(input *iam.DeleteAccessKeyInput) (*iam.DeleteAccessKeyOutput, error) {
	return c.iamClient.DeleteAccessKey(input)
}

func (c *awsClient) DeleteUser(input *iam.DeleteUserInput) (*iam.DeleteUserOutput, error) {
	return c.iamClient.DeleteUser(input)
}

func (c *awsClient) DeleteUserPolicy(input *iam.DeleteUserPolicyInput) (*iam.DeleteUserPolicyOutput, error) {
	return c.iamClient.DeleteUserPolicy(input)
}
func (c *awsClient) GetUser(input *iam.GetUserInput) (*iam.GetUserOutput, error) {
	return c.iamClient.GetUser(input)
}

func (c *awsClient) ListAccessKeys(input *iam.ListAccessKeysInput) (*iam.ListAccessKeysOutput, error) {
	return c.iamClient.ListAccessKeys(input)
}

func (c *awsClient) ListUserPolicies(input *iam.ListUserPoliciesInput) (*iam.ListUserPoliciesOutput, error) {
	return c.iamClient.ListUserPolicies(input)
}

func (c *awsClient) PutRolePolicy(input *iam.PutRolePolicyInput) (*iam.PutRolePolicyOutput, error) {
	return c.iamClient.PutRolePolicy(input)
}

func (c *awsClient) PutUserPolicy(input *iam.PutUserPolicyInput) (*iam.PutUserPolicyOutput, error) {
	return c.iamClient.PutUserPolicy(input)
}

func (c *awsClient) GetRole(input *iam.GetRoleInput) (*iam.GetRoleOutput, error) {
	return c.iamClient.GetRole(input)
}

func (c *awsClient) GetUserPolicy(input *iam.GetUserPolicyInput) (*iam.GetUserPolicyOutput, error) {
	return c.iamClient.GetUserPolicy(input)
}

func (c *awsClient) SimulatePrincipalPolicy(input *iam.SimulatePrincipalPolicyInput) (*iam.SimulatePolicyResponse, error) {
	return c.iamClient.SimulatePrincipalPolicy(input)
}

func (c *awsClient) SimulatePrincipalPolicyPages(input *iam.SimulatePrincipalPolicyInput, fn func(*iam.SimulatePolicyResponse, bool) bool) error {
	return c.iamClient.SimulatePrincipalPolicyPages(input, fn)
}

func (c *awsClient) TagUser(input *iam.TagUserInput) (*iam.TagUserOutput, error) {
	return c.iamClient.TagUser(input)
}

func (c *awsClient) ListOpenIDConnectProviders(input *iam.ListOpenIDConnectProvidersInput) (*iam.ListOpenIDConnectProvidersOutput, error) {
	return c.iamClient.ListOpenIDConnectProviders(input)
}

func (c *awsClient) CreateOpenIDConnectProvider(input *iam.CreateOpenIDConnectProviderInput) (*iam.CreateOpenIDConnectProviderOutput, error) {
	return c.iamClient.CreateOpenIDConnectProvider(input)
}

func (c *awsClient) TagOpenIDConnectProvider(input *iam.TagOpenIDConnectProviderInput) (*iam.TagOpenIDConnectProviderOutput, error) {
	return c.iamClient.TagOpenIDConnectProvider(input)
}

func (c *awsClient) UpdateAssumeRolePolicy(input *iam.UpdateAssumeRolePolicyInput) (*iam.UpdateAssumeRolePolicyOutput, error) {
	return c.iamClient.UpdateAssumeRolePolicy(input)
}

func (c *awsClient) GetOpenIDConnectProvider(input *iam.GetOpenIDConnectProviderInput) (*iam.GetOpenIDConnectProviderOutput, error) {
	return c.iamClient.GetOpenIDConnectProvider(input)
}

func (c *awsClient) DeleteOpenIDConnectProvider(input *iam.DeleteOpenIDConnectProviderInput) (*iam.DeleteOpenIDConnectProviderOutput, error) {
	return c.iamClient.DeleteOpenIDConnectProvider(input)
}

func (c *awsClient) ListRoles(input *iam.ListRolesInput) (*iam.ListRolesOutput, error) {
	return c.iamClient.ListRoles(input)
}

func (c *awsClient) DeleteRole(input *iam.DeleteRoleInput) (*iam.DeleteRoleOutput, error) {
	return c.iamClient.DeleteRole(input)
}

func (c *awsClient) ListRolePolicies(input *iam.ListRolePoliciesInput) (*iam.ListRolePoliciesOutput, error) {
	return c.iamClient.ListRolePolicies(input)
}

func (c *awsClient) DeleteRolePolicy(input *iam.DeleteRolePolicyInput) (*iam.DeleteRolePolicyOutput, error) {
	return c.iamClient.DeleteRolePolicy(input)
}

func (c *awsClient) CreateBucket(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	return c.s3Client.CreateBucket(input)
}

func (c *awsClient) PutBucketTagging(input *s3.PutBucketTaggingInput) (*s3.PutBucketTaggingOutput, error) {
	return c.s3Client.PutBucketTagging(input)
}

func (c *awsClient) GetBucketTagging(input *s3.GetBucketTaggingInput) (*s3.GetBucketTaggingOutput, error) {
	return c.s3Client.GetBucketTagging(input)
}

func (c *awsClient) DeleteBucket(input *s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error) {
	return c.s3Client.DeleteBucket(input)
}

func (c *awsClient) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return c.s3Client.PutObject(input)
}

func (c *awsClient) ListObjects(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	return c.s3Client.ListObjects(input)
}

func (c *awsClient) GetObjectTagging(input *s3.GetObjectTaggingInput) (*s3.GetObjectTaggingOutput, error) {
	return c.s3Client.GetObjectTagging(input)
}

func (c *awsClient) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return c.s3Client.DeleteObject(input)
}

func (c *awsClient) PutPublicAccessBlock(input *s3.PutPublicAccessBlockInput) (*s3.PutPublicAccessBlockOutput, error) {
	return c.s3Client.PutPublicAccessBlock(input)
}

func (c *awsClient) PutBucketPolicy(input *s3.PutBucketPolicyInput) (*s3.PutBucketPolicyOutput, error) {
	return c.s3Client.PutBucketPolicy(input)
}

func (c *awsClient) CreateCloudFrontOriginAccessIdentity(input *cloudfront.CreateCloudFrontOriginAccessIdentityInput) (*cloudfront.CreateCloudFrontOriginAccessIdentityOutput, error) {
	return c.cloudFrontClient.CreateCloudFrontOriginAccessIdentity(input)
}

func (c *awsClient) DeleteCloudFrontOriginAccessIdentity(input *cloudfront.DeleteCloudFrontOriginAccessIdentityInput) (*cloudfront.DeleteCloudFrontOriginAccessIdentityOutput, error) {
	return c.cloudFrontClient.DeleteCloudFrontOriginAccessIdentity(input)
}

func (c *awsClient) GetCloudFrontOriginAccessIdentity(input *cloudfront.GetCloudFrontOriginAccessIdentityInput) (*cloudfront.GetCloudFrontOriginAccessIdentityOutput, error) {
	return c.cloudFrontClient.GetCloudFrontOriginAccessIdentity(input)
}

func (c *awsClient) ListCloudFrontOriginAccessIdentities(input *cloudfront.ListCloudFrontOriginAccessIdentitiesInput) (*cloudfront.ListCloudFrontOriginAccessIdentitiesOutput, error) {
	return c.cloudFrontClient.ListCloudFrontOriginAccessIdentities(input)
}
func (c *awsClient) CreateCloudFrontDistributionWithTags(input *cloudfront.CreateDistributionWithTagsInput) (*cloudfront.CreateDistributionWithTagsOutput, error) {
	return c.cloudFrontClient.CreateDistributionWithTags(input)
}

func (c *awsClient) DeleteCloudFrontDistribution(input *cloudfront.DeleteDistributionInput) (*cloudfront.DeleteDistributionOutput, error) {
	return c.cloudFrontClient.DeleteDistribution(input)
}

func (c *awsClient) GetCloudFrontDistribution(input *cloudfront.GetDistributionInput) (*cloudfront.GetDistributionOutput, error) {
	return c.cloudFrontClient.GetDistribution(input)
}

func (c *awsClient) UpdateCloudFrontDistribution(input *cloudfront.UpdateDistributionInput) (*cloudfront.UpdateDistributionOutput, error) {
	return c.cloudFrontClient.UpdateDistribution(input)
}

func (c *awsClient) ListCloudFrontDistributions(input *cloudfront.ListDistributionsInput) (*cloudfront.ListDistributionsOutput, error) {
	return c.cloudFrontClient.ListDistributions(input)
}

func (c *awsClient) ListTagsForCloudFrontResource(input *cloudfront.ListTagsForResourceInput) (*cloudfront.ListTagsForResourceOutput, error) {
	return c.cloudFrontClient.ListTagsForResource(input)
}

// NewClient creates our client wrapper object for the actual AWS clients we use.
func NewClient(accessKeyID, secretAccessKey []byte, params *ClientParams) (Client, error) {
	var awsOpts session.Options

	agentText := "defaultAgent"

	if params != nil {
		if params.Region != "" {
			awsOpts.Config.Region = aws.String(params.Region)
		}

		if params.Endpoint != "" {
			awsOpts.Config.Endpoint = aws.String(params.Endpoint)
		}

		if params.CABundle != "" {
			awsOpts.CustomCABundle = strings.NewReader(params.CABundle)
		}

		if params.InfraName != "" {
			agentText = params.InfraName
		}
	}

	awsOpts.Config.Credentials = credentials.NewStaticCredentials(
		string(accessKeyID), string(secretAccessKey), "")

	s, err := session.NewSessionWithOptions(awsOpts)
	if err != nil {
		return nil, err
	}

	s.Handlers.Build.PushBackNamed(request.NamedHandler{
		Name: "openshift.io/cloud-credential-operator",
		Fn:   request.MakeAddToUserAgentHandler("openshift.io cloud-credential-operator", version.Get().String(), agentText),
	})

	return NewClientFromSession(s), nil
}

// NewClientFromSession will return a basic Client using only the provided awsSession
func NewClientFromSession(sess *session.Session) Client {
	return &awsClient{
		iamClient:        iam.New(sess),
		s3Client:         s3.New(sess),
		cloudFrontClient: cloudfront.New(sess),
	}
}
