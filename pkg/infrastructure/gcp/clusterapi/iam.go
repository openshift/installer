package clusterapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	resourcemanager "google.golang.org/api/cloudresourcemanager/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	"k8s.io/apimachinery/pkg/util/wait"

	gcp "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
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
		"roles/artifactregistry.reader",
	}
}

// GetSharedVPCRoles returns the pre-defined roles for a shared VPC installation.
func GetSharedVPCRoles() []string {
	return []string{
		"roles/compute.networkAdmin",
		"roles/compute.securityAdmin",
	}
}

// CreateServiceAccount is used to create a service account for a compute instance.
func CreateServiceAccount(ctx context.Context, infraID, projectID, role string, endpoint *gcptypes.PSCEndpoint) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	opts := []option.ClientOption{}
	if gcptypes.ShouldUseEndpointForInstaller(endpoint) {
		opts = append(opts, gcp.CreateEndpointOption(endpoint.Name, gcp.ServiceNameGCPIAM))
	}
	service, err := gcp.GetIAMService(ctx, opts...)
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
			return sa.Email, nil
		}
		time.Sleep(retryTime)
	}

	return "", fmt.Errorf("failure creating service account: %w", err)
}

// AddServiceAccountRoles adds predefined roles for service account.
func AddServiceAccountRoles(ctx context.Context, projectID, serviceAccountID string, roles []string, endpoint *gcptypes.PSCEndpoint) error {
	// Get cloudresourcemanager service
	// The context timeout must be greater in time than the exponential backoff below
	ctx, cancel := context.WithTimeout(ctx, time.Minute*2)
	defer cancel()

	opts := []option.ClientOption{}
	if gcptypes.ShouldUseEndpointForInstaller(endpoint) {
		opts = append(opts, gcp.CreateEndpointOption(endpoint.Name, gcp.ServiceNameGCPCloudResource))
	}

	service, err := gcp.GetCloudResourceService(ctx, opts...)
	if err != nil {
		return fmt.Errorf("failed to create resourcemanager service: %w", err)
	}

	backoff := wait.Backoff{
		Duration: 2 * time.Second,
		Factor:   2.0,
		Jitter:   1.0,
		Steps:    retryCount,
	}
	// Get and set the policy in a backoff loop.
	// If the policy set fails, the policy must be retrieved again via the get before retrying the set.
	var lastErr error
	if waitErr := wait.ExponentialBackoffWithContext(ctx, backoff, func(ctx context.Context) (bool, error) {
		policy, err := getPolicy(ctx, service, projectID)
		if isQuotaExceededError(err) {
			lastErr = err
			logrus.Debugf("Failed to get IAM policy, retrying after backoff")
			return false, nil
		} else if err != nil {
			return false, fmt.Errorf("failed to get IAM policy, unexpected error: %w", err)
		}

		member := fmt.Sprintf("serviceAccount:%s", serviceAccountID)
		for _, role := range roles {
			if err := addMemberToRole(policy, role, member); err != nil {
				return false, fmt.Errorf("failed to add role %s to %s: %w", role, member, err)
			}
		}

		err = setPolicy(ctx, service, projectID, policy)
		if err != nil {
			if isConflictError(err) {
				lastErr = err
				logrus.Debugf("Concurrent IAM policy changes, restarting read/modify/write")
				return false, nil
			} else if isBadStatusError(err) {
				// Documented here, https://cloud.google.com/iam/docs/retry-strategy, google
				// indicates that a service account may be created but not active for up to
				// 60 seconds. This behavior was causing a failure here when setting the policy
				// resulting in a 400 error from the API. If this error occurs retry with an
				// exponential backoff.
				lastErr = err
				logrus.Debugf("bad request, unexpected error: %s", err.Error())
				return false, nil
			}
			return false, fmt.Errorf("failed to set IAM policy, unexpected error: %w", err)
		}
		logrus.Debugf("Successfully set IAM policy")
		return true, nil
	}); waitErr != nil {
		if wait.Interrupted(waitErr) {
			return fmt.Errorf("failed to set IAM policy: %w", lastErr)
		}
		return waitErr
	}
	return nil
}

// getPolicy gets the project's IAM policy.
func getPolicy(ctx context.Context, crmService *resourcemanager.Service, projectID string) (*resourcemanager.Policy, error) {
	request := &resourcemanager.GetIamPolicyRequest{
		Options: &resourcemanager.GetPolicyOptions{
			RequestedPolicyVersion: 3,
		},
	}
	policy, err := crmService.Projects.GetIamPolicy(fmt.Sprintf("projects/%s", projectID), request).Context(ctx).Do()
	return policy, err
}

// setPolicy sets the project's IAM policy.
func setPolicy(ctx context.Context, crmService *resourcemanager.Service, projectID string, policy *resourcemanager.Policy) error {
	request := &resourcemanager.SetIamPolicyRequest{
		Policy: policy,
	}
	_, err := crmService.Projects.SetIamPolicy(fmt.Sprintf("projects/%s", projectID), request).Context(ctx).Do()
	return err
}

// addMemberToRole adds a member to a role binding.
func addMemberToRole(policy *resourcemanager.Policy, role, member string) error {
	var policyBinding *resourcemanager.Binding

	for _, binding := range policy.Bindings {
		if binding.Role == role {
			for _, m := range binding.Members {
				if m == member {
					logrus.Debugf("found %s role, member %s already exists", role, member)
					return nil
				}
			}
			policyBinding = binding
		}
	}

	if policyBinding == nil {
		policyBinding = &resourcemanager.Binding{
			Role:    role,
			Members: []string{member},
		}
		logrus.Debugf("creating new policy binding for %s role and %s member", role, member)
		policy.Bindings = append(policy.Bindings, policyBinding)
	}

	policyBinding.Members = append(policyBinding.Members, member)
	logrus.Debugf("adding %s role, added %s member", role, member)
	return nil
}

// isConflictError returns true if error matches conflict on concurrent policy sets.
func isConflictError(err error) bool {
	var ae *googleapi.Error
	if errors.As(err, &ae) && (ae.Code == http.StatusConflict || ae.Code == http.StatusPreconditionFailed) {
		return true
	}
	return false
}

// isQuotaExceededError returns true if the error matches quota exceeded.
func isQuotaExceededError(err error) bool {
	var ae *googleapi.Error
	if errors.As(err, &ae) && (ae.Code == http.StatusTooManyRequests) {
		return true
	}
	return false
}

func isBadStatusError(err error) bool {
	var ae *googleapi.Error
	return errors.As(err, &ae) && (ae.Code == http.StatusBadRequest)
}
