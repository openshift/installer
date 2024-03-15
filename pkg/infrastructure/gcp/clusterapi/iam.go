package clusterapi

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	resourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	iam "google.golang.org/api/iam/v1"
	"google.golang.org/api/option"

	gcp "github.com/openshift/installer/pkg/asset/installconfig/gcp"
)

const (
	retryTime  = 10 * time.Second
	retryCount = 6
)

func defaultServiceAccountID(infraID, projectID, role string) string {
	// The account id is used to generate the service account email address,
	// it should not contain the email suffixi. It is unique within a project,
	// must be 6-30 characters long, and match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])`
	return fmt.Sprintf("%s-%s", infraID, role[0:1])
}

// GetMasterRoles returns the pre-defined roles for a master node.
// Roles are described here https://cloud.google.com/iam/docs/understanding-roles#predefined_roles.
func GetMasterRoles() []string {
	return []string{
		"roles/compute.instanceAdmin",
		"roles/compute.networkAdmin",
		"roles/compute.securityAdmin",
		"roles/storage.admin",
	}
}

// GetWorkerRoles returns the pre-defined roles for a worker node.
func GetWorkerRoles() []string {
	return []string{
		"roles/compute.viewer",
		"roles/storage.admin",
	}
}

// CreateServiceAccount is used to create a service account for a compute instance.
func CreateServiceAccount(ctx context.Context, infraID, projectID, role string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	ssn, err := gcp.GetSession(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get session: %w", err)
	}
	service, err := iam.NewService(ctx, option.WithCredentials(ssn.Credentials))
	if err != nil {
		return "", fmt.Errorf("failed to create IAM service: %w", err)
	}

	accountID := defaultServiceAccountID(infraID, projectID, role)
	displayName := fmt.Sprintf("%s-%s-node", infraID, role)

	request := &iam.CreateServiceAccountRequest{
		AccountId: accountID,
		ServiceAccount: &iam.ServiceAccount{
			Description: "The service account used by the instances.",
			DisplayName: displayName,
		},
	}

	sa, err := service.Projects.ServiceAccounts.Create("projects/"+projectID, request).Do()
	if err != nil {
		return "", fmt.Errorf("Projects.ServiceAccounts.Create: %w", err)
	}

	// Poll for service account
	for i := 0; i < retryCount; i++ {
		_, err := service.Projects.ServiceAccounts.Get(sa.Name).Do()
		if err == nil {
			logrus.Debugf("Service account created for %s", accountID)
			return accountID, nil
		}
		time.Sleep(retryTime)
	}

	return "", fmt.Errorf("failure creating service account: %w", err)
}

// AddServiceAccountRoles adds predefined roles for service account.
func AddServiceAccountRoles(ctx context.Context, projectID, serviceAccountID string, roles []string) error {
	policy, err := getProjectIAMPolicy(ctx, projectID)
	if err != nil {
		return err
	}

	for _, role := range roles {
		err = addMemberToRole(policy, role, serviceAccountID)
		if err != nil {
			return fmt.Errorf("failed to add role %s to %s: %w", role, serviceAccountID, err)
		}
	}

	return nil
}

func getProjectIAMPolicy(ctx context.Context, projectID string) (*resourcemanager.Policy, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()
	req := &resourcemanager.GetIamPolicyRequest{}

	ssn, err := gcp.GetSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	service, err := resourcemanager.NewService(ctx, option.WithCredentials(ssn.Credentials))
	if err != nil {
		return nil, fmt.Errorf("failed to create resourcemanager service: %w", err)
	}

	policy, err := service.Projects.GetIamPolicy(projectID, req).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch project IAM policy: %w", err)
	}
	return policy, nil
}

// addMemberToRole adds a member to a role binding.
func addMemberToRole(policy *resourcemanager.Policy, role, member string) error {
	for _, binding := range policy.Bindings {
		if binding.Role == role {
			for _, m := range binding.Members {
				if m == member {
					logrus.Debugf("found %s role, member %s already exists", role, member)
					return nil
				}
			}
			binding.Members = append(binding.Members, member)
			logrus.Debugf("found %s role, added %s member", role, member)
			return nil
		}
	}

	return fmt.Errorf("role %s not found, member %s not added", role, member)
}
