package aws

import (
	"context"
	"fmt"
	"strings"

	iamv2 "github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
)

// IamRoleSearch holds data to search for IAM roles.
type IamRoleSearch struct {
	Client    *iamv2.Client
	Filters   []Filter
	Logger    logrus.FieldLogger
	Unmatched map[string]struct{}
}

func (search *IamRoleSearch) find(ctx context.Context) (arns []string, names []string, returnErr error) {
	if search.Unmatched == nil {
		search.Unmatched = map[string]struct{}{}
	}

	var lastError error

	paginator := iamv2.NewListRolesPaginator(search.Client, &iamv2.ListRolesInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to list roles: %w", err)
		}

		for _, role := range page.Roles {
			if _, ok := search.Unmatched[*role.Arn]; ok {
				continue
			}

			// Unfortunately role.Tags is empty from ListRoles, so we need to query each one
			response, err := search.Client.GetRole(ctx, &iamv2.GetRoleInput{RoleName: role.RoleName})
			if err != nil {
				var awsErr awserr.Error
				if errors.As(err, &awsErr) {
					switch {
					case awsErr.Code() == iam.ErrCodeNoSuchEntityException:
						// The role does not exist.
						// Ignore this IAM Role and donot report this error via
						// lastError
						search.Unmatched[*role.Arn] = exists
					case strings.Contains(err.Error(), "AccessDenied"):
						// Installer does not have access to this IAM role
						// Ignore this IAM Role and donot report this error via
						// lastError
						search.Unmatched[*role.Arn] = exists
					default:
						if lastError != nil {
							search.Logger.Debug(lastError)
						}
						lastError = fmt.Errorf("get tags for %s: %w", *role.Arn, err)
					}
				}
			} else {
				role = *response.Role
				tags := make(map[string]string, len(role.Tags))
				for _, tag := range role.Tags {
					tags[*tag.Key] = *tag.Value
				}
				if tagMatch(search.Filters, tags) {
					arns = append(arns, *role.Arn)
					names = append(names, *role.RoleName)
				} else {
					search.Unmatched[*role.Arn] = exists
				}
			}
		}
	}

	return arns, names, lastError
}

// IamUserSearch holds data to search for IAM users.
type IamUserSearch struct {
	client    *iamv2.Client
	filters   []Filter
	logger    logrus.FieldLogger
	unmatched map[string]struct{}
}

func (search *IamUserSearch) arns(ctx context.Context) ([]string, error) {
	if search.unmatched == nil {
		search.unmatched = map[string]struct{}{}
	}

	arns := []string{}
	var lastError error
	paginator := iamv2.NewListUsersPaginator(search.client, &iamv2.ListUsersInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list users: %w", err)
		}

		search.logger.Debugf("iterating over a page of %d IAM users", len(page.Users))
		for _, user := range page.Users {
			if _, ok := search.unmatched[*user.Arn]; ok {
				continue
			}

			// Unfortunately user.Tags is empty from ListUsers, so we need to query each one
			response, err := search.client.GetUser(ctx, &iamv2.GetUserInput{UserName: aws.String(*user.UserName)})
			if err != nil {

				var awsErr awserr.Error
				if errors.As(err, &awsErr) {
					switch {
					case awsErr.Code() == iam.ErrCodeNoSuchEntityException:
						// The role does not exist.
						// Ignore this IAM Role and do not report this error via lastError.
						search.unmatched[*user.Arn] = exists
					case strings.Contains(err.Error(), "AccessDenied"):
						// Installer does not have access to this IAM role.
						// Ignore this IAM Role and do not report this error via lastError.
						search.unmatched[*user.Arn] = exists
					default:
						if lastError != nil {
							search.logger.Debug(lastError)
						}
						lastError = errors.Wrapf(err, "get tags for %s", *user.Arn)
					}
				}
			} else {
				user = *response.User
				tags := make(map[string]string, len(user.Tags))
				for _, tag := range user.Tags {
					tags[*tag.Key] = *tag.Value
				}
				if tagMatch(search.filters, tags) {
					arns = append(arns, *user.Arn)
				} else {
					search.unmatched[*user.Arn] = exists
				}
			}
		}
	}

	return arns, lastError
}

// findIAMRoles returns the IAM roles for the cluster.
//
//	deleted - the resources that have already been deleted. Any resources specified in this set will be ignored.
func findIAMRoles(ctx context.Context, search *IamRoleSearch, deleted sets.Set[string], logger logrus.FieldLogger) (sets.Set[string], error) {
	logger.Debug("search for IAM roles")
	resources, _, err := search.find(ctx)
	if err != nil {
		logger.Info(err)
		return nil, err
	}
	return sets.New[string](resources...).Difference(deleted), nil
}

// findIAMUsers returns the IAM users for the cluster.
//
//	deleted - the resources that have already been deleted. Any resources specified in this set will be ignored.
func findIAMUsers(ctx context.Context, search *IamUserSearch, deleted sets.Set[string], logger logrus.FieldLogger) (sets.Set[string], error) {
	logger.Debug("search for IAM users")
	resources, err := search.arns(ctx)
	if err != nil {
		logger.Info(err)
		return nil, err
	}
	return sets.New[string](resources...).Difference(deleted), nil
}

func (o *ClusterUninstaller) deleteIAM(ctx context.Context, arn arn.ARN) error {
	resourceType, id, err := splitSlash("resource", arn.Resource)
	if err != nil {
		return err
	}
	logger := o.Logger.WithField("id", id)

	switch resourceType {
	case "instance-profile":
		return deleteIAMInstanceProfile(ctx, o.IAMClient, arn, logger)
	case "role":
		return deleteIAMRole(ctx, o.IAMClient, arn, logger)
	case "user":
		return deleteIAMUser(ctx, o.IAMClient, id, logger)
	default:
		return errors.Errorf("unrecognized EC2 resource type %s", resourceType)
	}
}

func deleteIAMInstanceProfileByName(ctx context.Context, client *iamv2.Client, name *string, logger logrus.FieldLogger) error {
	_, err := client.DeleteInstanceProfile(ctx, &iamv2.DeleteInstanceProfileInput{
		InstanceProfileName: name,
	})
	if err != nil {
		if err.(awserr.Error).Code() == iam.ErrCodeNoSuchEntityException {
			return nil
		}
		return err
	}
	logger.WithField("InstanceProfileName", *name).Info("Deleted")
	return err
}

func deleteIAMInstanceProfile(ctx context.Context, client *iamv2.Client, profileARN arn.ARN, logger logrus.FieldLogger) error {
	resourceType, name, err := splitSlash("resource", profileARN.Resource)
	if err != nil {
		return err
	}

	if resourceType != "instance-profile" {
		return errors.Errorf("%s ARN passed to deleteIAMInstanceProfile: %s", resourceType, profileARN.String())
	}

	response, err := client.GetInstanceProfile(ctx, &iamv2.GetInstanceProfileInput{
		InstanceProfileName: &name,
	})
	if err != nil {
		if err.(awserr.Error).Code() == iam.ErrCodeNoSuchEntityException {
			return nil
		}
		return err
	}
	profile := response.InstanceProfile

	for _, role := range profile.Roles {
		_, err = client.RemoveRoleFromInstanceProfile(ctx, &iamv2.RemoveRoleFromInstanceProfileInput{
			InstanceProfileName: profile.InstanceProfileName,
			RoleName:            role.RoleName,
		})
		if err != nil {
			return errors.Wrapf(err, "dissociating %s", *role.RoleName)
		}
		logger.WithField("name", name).WithField("role", *role.RoleName).Info("Disassociated")
	}

	logger = logger.WithField("arn", profileARN.String())
	if err := deleteIAMInstanceProfileByName(ctx, client, profile.InstanceProfileName, logger); err != nil {
		return err
	}

	return nil
}

func deleteIAMRole(ctx context.Context, client *iamv2.Client, roleARN arn.ARN, logger logrus.FieldLogger) error {
	resourceType, name, err := splitSlash("resource", roleARN.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("name", name)

	if resourceType != "role" {
		return errors.Errorf("%s ARN passed to deleteIAMRole: %s", resourceType, roleARN.String())
	}

	var lastError error
	paginator := iamv2.NewListRolePoliciesPaginator(client, &iamv2.ListRolePoliciesInput{
		RoleName: aws.String(name),
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("listing IAM role policies: %w", err)
		}

		for _, policy := range page.PolicyNames {
			_, err := client.DeleteRolePolicy(ctx, &iamv2.DeleteRolePolicyInput{
				RoleName:   &name,
				PolicyName: &policy,
			})
			if err != nil {
				if lastError != nil {
					logger.Debug(lastError)
				}
				lastError = fmt.Errorf("deleting IAM role policy %s: %w", policy, err)
			}
			logger.WithField("policy", policy).Info("Deleted")
		}
	}

	if lastError != nil {
		return lastError
	}

	attachedPaginator := iamv2.NewListAttachedRolePoliciesPaginator(client, &iamv2.ListAttachedRolePoliciesInput{
		RoleName: aws.String(name),
	})
	for attachedPaginator.HasMorePages() {
		page, err := attachedPaginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("listing attached IAM role policies: %w", err)
		}

		for _, policy := range page.AttachedPolicies {
			_, err := client.DetachRolePolicy(ctx, &iamv2.DetachRolePolicyInput{
				RoleName:  &name,
				PolicyArn: policy.PolicyArn,
			})
			if err != nil {
				if lastError != nil {
					logger.Debug(lastError)
				}
				lastError = fmt.Errorf("detaching IAM role policy %s: %w", *policy.PolicyName, err)
			}
			logger.WithField("policy", *policy.PolicyName).Info("Detached")
		}
	}

	return lastError
}

func deleteIAMUser(ctx context.Context, client *iamv2.Client, id string, logger logrus.FieldLogger) error {
	var lastError error

	policiesPaginator := iamv2.NewListUserPoliciesPaginator(client, &iamv2.ListUserPoliciesInput{
		UserName: aws.String(id),
	})
	for policiesPaginator.HasMorePages() {
		page, err := policiesPaginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("listing IAM user policies for %s: %w", id, err)
		}

		for _, policy := range page.PolicyNames {
			_, err := client.DeleteUserPolicy(ctx, &iamv2.DeleteUserPolicyInput{
				UserName:   &id,
				PolicyName: &policy,
			})
			if err != nil {
				if lastError != nil {
					logger.Debug(lastError)
				}
				lastError = fmt.Errorf("deleting IAM user policy %s: %w", policy, err)
			}
			logger.WithField("policy", policy).Info("Deleted")
		}
	}

	if lastError != nil {
		return lastError
	}

	keysPaginator := iamv2.NewListAccessKeysPaginator(client, &iamv2.ListAccessKeysInput{
		UserName: aws.String(id),
	})
	for keysPaginator.HasMorePages() {
		page, err := keysPaginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("listing access keys for %s: %w", id, err)
		}

		for _, key := range page.AccessKeyMetadata {
			_, err := client.DeleteAccessKey(ctx, &iamv2.DeleteAccessKeyInput{
				UserName:    &id,
				AccessKeyId: key.AccessKeyId,
			})
			if err != nil {
				if lastError != nil {
					logger.Debug(lastError)
				}
				lastError = fmt.Errorf("deleting IAM access key %s: %w", *key.AccessKeyId, err)
			}
		}
	}

	if lastError != nil {
		return lastError
	}

	_, err := client.DeleteUser(ctx, &iamv2.DeleteUserInput{
		UserName: &id,
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}
