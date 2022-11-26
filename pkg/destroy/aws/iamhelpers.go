package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
)

type iamRoleSearch struct {
	client    IAMAPI
	filters   []Filter
	logger    logrus.FieldLogger
	unmatched map[string]struct{}
}

func (search *iamRoleSearch) find(ctx context.Context) (arns []string, names []string, returnErr error) {
	if search.unmatched == nil {
		search.unmatched = map[string]struct{}{}
	}

	var lastError error
	err := search.client.ListRolesPages(ctx,
		func(results *iam.ListRolesOutput, lastPage bool) bool {
			search.logger.Debugf("iterating over a page of %d IAM roles", len(results.Roles))
			for _, role := range results.Roles {
				if _, ok := search.unmatched[*role.Arn]; ok {
					continue
				}

				// Unfortunately role.Tags is empty from ListRoles, so we need to query each one
				response, err := search.client.GetRole(ctx, *role.RoleName)
				if err != nil {
					if err.(awserr.Error).Code() == iam.ErrCodeNoSuchEntityException {
						search.unmatched[*role.Arn] = exists
					} else {
						if lastError != nil {
							search.logger.Debug(lastError)
						}
						lastError = errors.Wrapf(err, "get tags for %s", *role.Arn)
					}
				} else {
					role = response.Role
					tags := make(map[string]string, len(role.Tags))
					for _, tag := range role.Tags {
						tags[*tag.Key] = *tag.Value
					}
					if tagMatch(search.filters, tags) {
						arns = append(arns, *role.Arn)
						names = append(names, *role.RoleName)
					} else {
						search.unmatched[*role.Arn] = exists
					}
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return arns, names, lastError
	}
	return arns, names, err
}

type iamUserSearch struct {
	client    IAMAPI
	filters   []Filter
	logger    logrus.FieldLogger
	unmatched map[string]struct{}
}

func (search *iamUserSearch) arns(ctx context.Context) ([]string, error) {
	if search.unmatched == nil {
		search.unmatched = map[string]struct{}{}
	}

	arns := []string{}
	var lastError error
	err := search.client.ListUsersPages(ctx,
		func(results *iam.ListUsersOutput, lastPage bool) bool {
			search.logger.Debugf("iterating over a page of %d IAM users", len(results.Users))
			for _, user := range results.Users {
				if _, ok := search.unmatched[*user.Arn]; ok {
					continue
				}

				// Unfortunately user.Tags is empty from ListUsers, so we need to query each one
				response, err := search.client.GetUser(ctx, *user.UserName)
				if err != nil {
					if err.(awserr.Error).Code() == iam.ErrCodeNoSuchEntityException {
						search.unmatched[*user.Arn] = exists
					} else {
						if lastError != nil {
							search.logger.Debug(lastError)
						}
						lastError = errors.Wrapf(err, "get tags for %s", *user.Arn)
					}
				} else {
					user = response.User
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

			return !lastPage
		},
	)

	if lastError != nil {
		return arns, lastError
	}
	return arns, err
}

// findIAMRoles returns the IAM roles for the cluster.
//
//	deleted - the resources that have already been deleted. Any resources specified in this set will be ignored.
func findIAMRoles(ctx context.Context, search *iamRoleSearch, deleted sets.String, logger logrus.FieldLogger) (sets.String, error) {
	logger.Debug("search for IAM roles")
	resources, _, err := search.find(ctx)
	if err != nil {
		logger.Info(err)
		return nil, err
	}
	return sets.NewString(resources...).Difference(deleted), nil
}

// findIAMUsers returns the IAM users for the cluster.
//
//	deleted - the resources that have already been deleted. Any resources specified in this set will be ignored.
func findIAMUsers(ctx context.Context, search *iamUserSearch, deleted sets.String, logger logrus.FieldLogger) (sets.String, error) {
	logger.Debug("search for IAM users")
	resources, err := search.arns(ctx)
	if err != nil {
		logger.Info(err)
		return nil, err
	}
	return sets.NewString(resources...).Difference(deleted), nil
}

func deleteIAM(ctx context.Context, session *session.Session, arn arn.ARN, logger logrus.FieldLogger) error {
	client := NewIAMClient(session, logger)
	return deleteIAMWithClient(ctx, client, arn, logger)
}

// This is defined as a separate function so it can be tested with a mocked client
func deleteIAMWithClient(ctx context.Context, client IAMAPI, arn arn.ARN, logger logrus.FieldLogger) error {
	resourceType, id, err := splitSlash("resource", arn.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("id", id)

	switch resourceType {
	case "instance-profile":
		return deleteIAMInstanceProfile(ctx, client, arn, logger)
	case "role":
		return deleteIAMRole(ctx, client, arn, logger)
	case "user":
		return deleteIAMUser(ctx, client, id, logger)
	default:
		return errors.Errorf("unrecognized EC2 resource type %s", resourceType)
	}
}

func deleteIAMInstanceProfileByName(ctx context.Context, client IAMAPI, name *string, logger logrus.FieldLogger) error {
	err := client.DeleteInstanceProfile(ctx, *name)
	if err != nil {
		if err.(awserr.Error).Code() == iam.ErrCodeNoSuchEntityException {
			return nil
		}
		return err
	}
	logger.WithField("InstanceProfileName", *name).Info("Deleted")
	return err
}

func deleteIAMInstanceProfile(ctx context.Context, client IAMAPI, profileARN arn.ARN, logger logrus.FieldLogger) error {
	resourceType, name, err := splitSlash("resource", profileARN.Resource)
	if err != nil {
		return err
	}

	if resourceType != "instance-profile" {
		return errors.Errorf("%s ARN passed to deleteIAMInstanceProfile: %s", resourceType, profileARN.String())
	}

	response, err := client.GetInstanceProfile(ctx, name)
	if err != nil {
		if err.(awserr.Error).Code() == iam.ErrCodeNoSuchEntityException {
			return nil
		}
		return err
	}
	profile := response.InstanceProfile

	for _, role := range profile.Roles {
		err = client.RemoveRoleFromInstanceProfile(ctx, *profile.InstanceProfileName, *role.RoleName)
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

func deleteRolePolicies(ctx context.Context, client IAMAPI, role string, logger logrus.FieldLogger) error {
	var lastError error
	err := client.ListRolePoliciesPages(ctx, role,
		func(results *iam.ListRolePoliciesOutput, lastPage bool) bool {
			for _, policy := range results.PolicyNames {
				err := client.DeleteRolePolicy(ctx, role, *policy)
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrapf(err, "deleting IAM role policy %s", *policy)
				}
				logger.WithField("policy", *policy).Info("Deleted")
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	if err != nil {
		return errors.Wrap(err, "listing IAM role policies")
	}

	return nil
}

func deleteAttachedRolePolicies(ctx context.Context, client IAMAPI, role string, logger logrus.FieldLogger) error {
	var lastError error
	err := client.ListAttachedRolePoliciesPages(ctx, role,
		func(results *iam.ListAttachedRolePoliciesOutput, lastPage bool) bool {
			for _, policy := range results.AttachedPolicies {
				err := client.DetachRolePolicy(ctx, role, *policy.PolicyArn)
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrapf(err, "detaching IAM role policy %s", *policy.PolicyName)
				}
				logger.WithField("policy", *policy.PolicyName).Info("Detached")
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	if err != nil {
		return errors.Wrap(err, "listing attached IAM role policies")
	}

	return nil
}

func deleteInstanceProfilesForRole(ctx context.Context, client IAMAPI, role string, logger logrus.FieldLogger) error {
	var lastError error
	err := client.ListInstanceProfilesForRolePages(ctx, role,
		func(results *iam.ListInstanceProfilesForRoleOutput, lastPage bool) bool {
			for _, profile := range results.InstanceProfiles {
				parsed, err := arn.Parse(*profile.Arn)
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrap(err, "parse ARN for IAM instance profile")
					continue
				}

				err = deleteIAMInstanceProfile(ctx, client, parsed, logger)
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrapf(err, "deleting %s", parsed.String())
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	if err != nil {
		return errors.Wrap(err, "listing IAM instance profiles")
	}

	return nil
}

func deleteIAMRole(ctx context.Context, client IAMAPI, roleARN arn.ARN, logger logrus.FieldLogger) error {
	resourceType, name, err := splitSlash("resource", roleARN.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("name", name)

	if resourceType != "role" {
		return errors.Errorf("%s ARN passed to deleteIAMRole: %s", resourceType, roleARN.String())
	}

	err = deleteRolePolicies(ctx, client, name, logger)
	if err != nil {
		return err
	}

	err = deleteAttachedRolePolicies(ctx, client, name, logger)
	if err != nil {
		return err
	}

	err = deleteInstanceProfilesForRole(ctx, client, name, logger)
	if err != nil {
		return err
	}

	err = client.DeleteRole(ctx, name)
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteUserPolicies(ctx context.Context, client IAMAPI, id string, logger logrus.FieldLogger) error {
	var lastError error
	err := client.ListUserPoliciesPages(ctx, id,
		func(results *iam.ListUserPoliciesOutput, lastPage bool) bool {
			for _, policy := range results.PolicyNames {
				err := client.DeleteUserPolicy(ctx, id, *policy)
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrapf(err, "deleting IAM user policy %s", *policy)
				}
				logger.WithField("policy", *policy).Info("Deleted")
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	if err != nil {
		return errors.Wrap(err, "listing IAM user policies")
	}

	return nil
}

func deleteAccessKeys(ctx context.Context, client IAMAPI, id string, logger logrus.FieldLogger) error {
	var lastError error
	err := client.ListAccessKeysPages(ctx, id,
		func(results *iam.ListAccessKeysOutput, lastPage bool) bool {
			for _, key := range results.AccessKeyMetadata {
				err := client.DeleteAccessKey(ctx, id, *key.AccessKeyId)
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrapf(err, "deleting IAM access key %s", *key.AccessKeyId)
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	if err != nil {
		return errors.Wrap(err, "listing IAM access keys")
	}

	return nil
}

func deleteIAMUser(ctx context.Context, client IAMAPI, id string, logger logrus.FieldLogger) error {
	err := deleteUserPolicies(ctx, client, id, logger)
	if err != nil {
		return err
	}

	err = deleteAccessKeys(ctx, client, id, logger)
	if err != nil {
		return err
	}

	err = client.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}
