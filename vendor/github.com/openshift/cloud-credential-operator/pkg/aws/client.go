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
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/openshift/cloud-credential-operator/version"
)

//go:generate mockgen -source=./client.go -destination=./mock/client_generated.go -package=mock

// Client is a wrapper object for actual AWS SDK clients to allow for easier testing.
type Client interface {
	//IAM
	CreateAccessKey(*iam.CreateAccessKeyInput) (*iam.CreateAccessKeyOutput, error)
	CreateUser(*iam.CreateUserInput) (*iam.CreateUserOutput, error)
	DeleteAccessKey(*iam.DeleteAccessKeyInput) (*iam.DeleteAccessKeyOutput, error)
	DeleteUser(*iam.DeleteUserInput) (*iam.DeleteUserOutput, error)
	DeleteUserPolicy(*iam.DeleteUserPolicyInput) (*iam.DeleteUserPolicyOutput, error)
	GetUser(*iam.GetUserInput) (*iam.GetUserOutput, error)
	ListAccessKeys(*iam.ListAccessKeysInput) (*iam.ListAccessKeysOutput, error)
	ListUserPolicies(*iam.ListUserPoliciesInput) (*iam.ListUserPoliciesOutput, error)
	PutUserPolicy(*iam.PutUserPolicyInput) (*iam.PutUserPolicyOutput, error)
	GetUserPolicy(*iam.GetUserPolicyInput) (*iam.GetUserPolicyOutput, error)
	SimulatePrincipalPolicy(*iam.SimulatePrincipalPolicyInput) (*iam.SimulatePolicyResponse, error)
	SimulatePrincipalPolicyPages(*iam.SimulatePrincipalPolicyInput, func(*iam.SimulatePolicyResponse, bool) bool) error
	TagUser(*iam.TagUserInput) (*iam.TagUserOutput, error)
}

type awsClient struct {
	iamClient iamiface.IAMAPI
}

func (c *awsClient) CreateAccessKey(input *iam.CreateAccessKeyInput) (*iam.CreateAccessKeyOutput, error) {
	return c.iamClient.CreateAccessKey(input)
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

func (c *awsClient) PutUserPolicy(input *iam.PutUserPolicyInput) (*iam.PutUserPolicyOutput, error) {
	return c.iamClient.PutUserPolicy(input)
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

// NewClient creates our client wrapper object for the actual AWS clients we use.
func NewClient(accessKeyID, secretAccessKey []byte, region, infraName string) (Client, error) {
	awsConfig := &awssdk.Config{}

	if region != "" {
		awsConfig.Region = &region
	}

	awsConfig.Credentials = credentials.NewStaticCredentials(
		string(accessKeyID), string(secretAccessKey), "")

	s, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}
	s.Handlers.Build.PushBackNamed(request.NamedHandler{
		Name: "openshift.io/cloud-credential-operator",
		Fn:   request.MakeAddToUserAgentHandler("openshift.io cloud-credential-operator", version.Version, infraName),
	})

	return NewClientFromIAMClient(iam.New(s))
}

// NewClientFromIAMClient create a client from AWS IAM client.
func NewClientFromIAMClient(client iamiface.IAMAPI) (Client, error) {
	return &awsClient{
		iamClient: client,
	}, nil
}
