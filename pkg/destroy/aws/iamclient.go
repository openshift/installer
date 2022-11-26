package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=./iamclient.go -destination=mock/iamclient_generated.go -package=mock

// IAMAPI represents the calls made to the AWS EC2 API.
type IAMAPI interface {
	DeleteAccessKey(ctx context.Context, userName, id string) error
	DeleteInstanceProfile(ctx context.Context, name string) error
	DeleteRole(ctx context.Context, name string) error
	DeleteRolePolicy(ctx context.Context, roleName, name string) error
	DeleteUser(ctx context.Context, name string) error
	DeleteUserPolicy(ctx context.Context, userName, name string) error
	DetachRolePolicy(ctx context.Context, roleName, arn string) error
	GetInstanceProfile(ctx context.Context, name string) (*iam.GetInstanceProfileOutput, error)
	GetRole(ctx context.Context, name string) (*iam.GetRoleOutput, error)
	GetUser(ctx context.Context, name string) (*iam.GetUserOutput, error)
	ListAccessKeysPages(ctx context.Context, userName string, fn func(results *iam.ListAccessKeysOutput, lastPage bool) bool) error
	ListAttachedRolePoliciesPages(ctx context.Context, roleName string, fn func(results *iam.ListAttachedRolePoliciesOutput, lastPage bool) bool) error
	ListInstanceProfilesForRolePages(ctx context.Context, roleName string, fn func(results *iam.ListInstanceProfilesForRoleOutput, lastPage bool) bool) error
	ListRolePoliciesPages(ctx context.Context, roleName string, fn func(results *iam.ListRolePoliciesOutput, lastPage bool) bool) error
	ListRolesPages(ctx context.Context, fn func(*iam.ListRolesOutput, bool) bool) error
	ListUserPoliciesPages(ctx context.Context, userName string, fn func(results *iam.ListUserPoliciesOutput, lastPage bool) bool) error
	ListUsersPages(ctx context.Context, fn func(results *iam.ListUsersOutput, lastPage bool) bool) error
	UntagRole(ctx context.Context, roleName, tagKey string) error
	RemoveRoleFromInstanceProfile(ctx context.Context, name, roleName string) error
}

// IAMClient makes calls to the AWS IAM API.
type IAMClient struct {
	client *iam.IAM
	logger logrus.FieldLogger
}

// NewIAMClient returns a new IAM Client.
func NewIAMClient(awsSession *session.Session, logger logrus.FieldLogger) *IAMClient {
	return &IAMClient{
		client: iam.New(awsSession),
		logger: logger,
	}
}

// DeleteAccessKey deletes the access key with id `id` and user name `userName`.
func (c *IAMClient) DeleteAccessKey(ctx context.Context, userName, id string) error {
	_, err := c.client.DeleteAccessKeyWithContext(ctx, &iam.DeleteAccessKeyInput{
		UserName:    aws.String(userName),
		AccessKeyId: aws.String(id),
	})
	return err
}

// DeleteInstanceProfile deletes an instance profile named `name`.
func (c *IAMClient) DeleteInstanceProfile(ctx context.Context, name string) error {
	_, err := c.client.DeleteInstanceProfileWithContext(ctx, &iam.DeleteInstanceProfileInput{InstanceProfileName: aws.String(name)})
	return err
}

// DeleteRole deletes a role named `name`.
func (c *IAMClient) DeleteRole(ctx context.Context, name string) error {
	_, err := c.client.DeleteRoleWithContext(ctx, &iam.DeleteRoleInput{RoleName: aws.String(name)})
	return err
}

// DeleteRolePolicy deletes a policy named `name` from role named `roleName`.
func (c *IAMClient) DeleteRolePolicy(ctx context.Context, roleName, name string) error {
	_, err := c.client.DeleteRolePolicyWithContext(ctx, &iam.DeleteRolePolicyInput{
		RoleName:   aws.String(name),
		PolicyName: aws.String(name),
	})
	return err
}

// DeleteUser deletes a user named `name`.
func (c *IAMClient) DeleteUser(ctx context.Context, name string) error {
	_, err := c.client.DeleteUserWithContext(ctx, &iam.DeleteUserInput{UserName: aws.String(name)})
	return err
}

// DeleteUserPolicy deletes a user policy named `name` from user named `username`.
func (c *IAMClient) DeleteUserPolicy(ctx context.Context, userName, name string) error {
	_, err := c.client.DeleteUserPolicyWithContext(ctx, &iam.DeleteUserPolicyInput{
		UserName:   aws.String(userName),
		PolicyName: aws.String(name),
	})
	return err
}

// DetachRolePolicy detaches policy with arn `arn` from role named `roleName`.
func (c *IAMClient) DetachRolePolicy(ctx context.Context, roleName, arn string) error {
	_, err := c.client.DetachRolePolicyWithContext(ctx, &iam.DetachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String(arn),
	})
	return err
}

// ListAccessKeysPages runs `fn` for each AccessKeys page from user named `userName`.
func (c *IAMClient) ListAccessKeysPages(ctx context.Context, userName string, fn func(results *iam.ListAccessKeysOutput, lastPage bool) bool) error {
	return c.client.ListAccessKeysPagesWithContext(ctx, &iam.ListAccessKeysInput{UserName: aws.String(userName)}, fn)
}

// ListAttachedRolePoliciesPages runs `fn` for each Attached Role Policies page from role named `roleName“.
func (c *IAMClient) ListAttachedRolePoliciesPages(ctx context.Context, roleName string, fn func(results *iam.ListAttachedRolePoliciesOutput, lastPage bool) bool) error {
	return c.client.ListAttachedRolePoliciesPagesWithContext(ctx, &iam.ListAttachedRolePoliciesInput{RoleName: aws.String(roleName)}, fn)
}

// ListInstanceProfilesForRolePages runs `fn` for each Instance Profiles page from role named `roleName`.
func (c *IAMClient) ListInstanceProfilesForRolePages(ctx context.Context, roleName string, fn func(results *iam.ListInstanceProfilesForRoleOutput, lastPage bool) bool) error {
	return c.client.ListInstanceProfilesForRolePagesWithContext(ctx, &iam.ListInstanceProfilesForRoleInput{RoleName: aws.String(roleName)}, fn)
}

// ListUsersPages runs `fn` for each Users page.
func (c *IAMClient) ListUsersPages(ctx context.Context, fn func(results *iam.ListUsersOutput, lastPage bool) bool) error {
	return c.client.ListUsersPagesWithContext(ctx, &iam.ListUsersInput{}, fn)
}

// ListUserPoliciesPages runs `fn` for each User Policies page from user named `userName`.
func (c *IAMClient) ListUserPoliciesPages(ctx context.Context, userName string, fn func(results *iam.ListUserPoliciesOutput, lastPage bool) bool) error {
	return c.client.ListUserPoliciesPagesWithContext(ctx, &iam.ListUserPoliciesInput{UserName: aws.String(userName)}, fn)
}

// ListRolesPages runs `fn` for each Roles page.
func (c *IAMClient) ListRolesPages(ctx context.Context, fn func(results *iam.ListRolesOutput, lastPage bool) bool) error {
	return c.client.ListRolesPagesWithContext(ctx, &iam.ListRolesInput{}, fn)
}

// ListRolePoliciesPages runs `fn` for each Policies page from role named `roleName`.
func (c *IAMClient) ListRolePoliciesPages(ctx context.Context, roleName string, fn func(results *iam.ListRolePoliciesOutput, lastPage bool) bool) error {
	return c.client.ListRolePoliciesPagesWithContext(ctx, &iam.ListRolePoliciesInput{RoleName: aws.String(roleName)}, fn)
}

// GetInstanceProfile returns the Instance Profile named `name`.
func (c *IAMClient) GetInstanceProfile(ctx context.Context, name string) (*iam.GetInstanceProfileOutput, error) {
	return c.client.GetInstanceProfileWithContext(ctx, &iam.GetInstanceProfileInput{InstanceProfileName: &name})
}

// GetUser returns the User named `name`.
func (c *IAMClient) GetUser(ctx context.Context, name string) (*iam.GetUserOutput, error) {
	return c.client.GetUserWithContext(ctx, &iam.GetUserInput{UserName: aws.String(name)})
}

// GetRole returns the Role named `name`.
func (c *IAMClient) GetRole(ctx context.Context, name string) (*iam.GetRoleOutput, error) {
	return c.client.GetRoleWithContext(ctx, &iam.GetRoleInput{RoleName: aws.String(name)})
}

// UntagRole untags the role named `roleName` with tag `key`.
func (c *IAMClient) UntagRole(ctx context.Context, roleName, tagKey string) error {
	input := &iam.UntagRoleInput{
		RoleName: aws.String(roleName),
		TagKeys:  []*string{aws.String(tagKey)},
	}
	_, err := c.client.UntagRoleWithContext(ctx, input)
	return err
}

// RemoveRoleFromInstanceProfile removes Role named `roleName` from Instance Profile named `name`.
func (c *IAMClient) RemoveRoleFromInstanceProfile(ctx context.Context, name, roleName string) error {
	_, err := c.client.RemoveRoleFromInstanceProfileWithContext(ctx, &iam.RemoveRoleFromInstanceProfileInput{
		InstanceProfileName: aws.String(name),
		RoleName:            aws.String(roleName),
	})
	return err
}
