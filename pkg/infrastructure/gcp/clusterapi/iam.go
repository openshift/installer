package clusterapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/iam/apiv1/iampb"
	kms "cloud.google.com/go/kms/apiv1"
	"github.com/sirupsen/logrus"
	resourcemanager "google.golang.org/api/cloudresourcemanager/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// GetKMSRoles returns the KMS-related roles needed when customer-managed encryption keys are configured.
// These roles allow the service account to decrypt data encrypted with the specified KMS keys.
// Only grant these to service accounts that need to access KMS-encrypted resources.
func GetKMSRoles() []string {
	return []string{
		"roles/cloudkms.cryptoKeyEncrypterDecrypter",
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
	var lastPollErr error
	for i := 0; i < retryCount; i++ {
		_, lastPollErr = service.Projects.ServiceAccounts.Get(sa.Name).Do()
		if lastPollErr == nil {
			logrus.Infof("Service account %s created successfully", accountID)
			return sa.Email, nil
		}
		logrus.Debugf("Service account %s not yet available (attempt %d/%d): %v", accountID, i+1, retryCount, lastPollErr)
		time.Sleep(retryTime)
	}

	logrus.Errorf("Service account %s failed to become available after %d attempts over %v",
		accountID, retryCount, time.Duration(retryCount)*retryTime)
	return "", fmt.Errorf("service account %s not available after %d retries: %w", accountID, retryCount, lastPollErr)
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
			logrus.Warnf("IAM API quota exceeded for project %s - this may require waiting for quota reset or requesting a quota increase", projectID)
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
				logrus.Warnf("Failed to set IAM policy for project %s, service account %s, roles %v: %v. This may indicate the service account is not yet active. Retrying...",
					projectID, serviceAccountID, roles, err)
				return false, nil
			}
			return false, fmt.Errorf("failed to set IAM policy, unexpected error: %w", err)
		}
		logrus.Infof("Successfully set IAM policy for project %s, service account %s, roles %v", projectID, serviceAccountID, roles)
		return true, nil
	}); waitErr != nil {
		if wait.Interrupted(waitErr) {
			logrus.Errorf("IAM policy update exhausted retries for project %s, service account %s, roles %v. Last error: %v",
				projectID, serviceAccountID, roles, lastErr)
			return fmt.Errorf("failed to set IAM policy for project %s, service account %s after retries (last error: %w)",
				projectID, serviceAccountID, lastErr)
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

// GetCloudStorageServiceAccount returns the email address of the Google Cloud Storage service account
// for the specified project. This is the Google-managed service account that performs bucket operations,
// including encryption/decryption when customer-managed KMS keys are used.
//
// The Cloud Storage service account email format is: service-{PROJECT_NUMBER}@gs-project-accounts.iam.gserviceaccount.com.
func GetCloudStorageServiceAccount(ctx context.Context, projectID string, endpoint *gcptypes.PSCEndpoint) (string, error) {
	opts := []option.ClientOption{}
	if gcptypes.ShouldUseEndpointForInstaller(endpoint) {
		opts = append(opts, gcp.CreateEndpointOption(endpoint.Name, gcp.ServiceNameGCPCloudResource))
	}

	service, err := gcp.GetCloudResourceService(ctx, opts...)
	if err != nil {
		return "", fmt.Errorf("failed to create resourcemanager service: %w", err)
	}

	project, err := service.Projects.Get(fmt.Sprintf("projects/%s", projectID)).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get project %s: %w", projectID, err)
	}

	// Extract project number from the project name (format: "projects/123456789")
	projectNumber := strings.TrimPrefix(project.Name, "projects/")

	return fmt.Sprintf("service-%s@gs-project-accounts.iam.gserviceaccount.com", projectNumber), nil
}

// GetComputeEngineServiceAccount returns the email address of the Google Compute Engine service account
// for the specified project. This is the Google-managed service account that performs compute operations,
// including disk encryption/decryption when customer-managed KMS keys are used.
//
// The Compute Engine service account email format is: service-{PROJECT_NUMBER}@compute-system.iam.gserviceaccount.com.
func GetComputeEngineServiceAccount(ctx context.Context, projectID string, endpoint *gcptypes.PSCEndpoint) (string, error) {
	opts := []option.ClientOption{}
	if gcptypes.ShouldUseEndpointForInstaller(endpoint) {
		opts = append(opts, gcp.CreateEndpointOption(endpoint.Name, gcp.ServiceNameGCPCloudResource))
	}

	service, err := gcp.GetCloudResourceService(ctx, opts...)
	if err != nil {
		return "", fmt.Errorf("failed to create resourcemanager service: %w", err)
	}

	project, err := service.Projects.Get(fmt.Sprintf("projects/%s", projectID)).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get project %s: %w", projectID, err)
	}

	// Extract project number from the project name (format: "projects/123456789")
	projectNumber := strings.TrimPrefix(project.Name, "projects/")

	return fmt.Sprintf("service-%s@compute-system.iam.gserviceaccount.com", projectNumber), nil
}

// GrantKMSKeyIAMPermission grants a service account permission to use a KMS key by updating
// the key's IAM policy. This is required to allow Google-managed service accounts (like Cloud Storage
// and Compute Engine) to encrypt/decrypt data using customer-managed KMS keys.
//
// The function uses a read-modify-write pattern with exponential backoff to handle concurrent
// policy updates. It requires the installer service account to have cloudkms.cryptoKeys.getIamPolicy
// and cloudkms.cryptoKeys.setIamPolicy permissions.
func GrantKMSKeyIAMPermission(ctx context.Context, kmsKey *gcptypes.KMSKeyReference, projectID, serviceAccountEmail, role string) error {
	// The context timeout must be greater in time than the exponential backoff below
	ctx, cancel := context.WithTimeout(ctx, time.Minute*2)
	defer cancel()

	// Get session for credentials - CRITICAL: use installer's session credentials instead of
	// default application credentials to ensure we access the correct GCP project
	ssn, err := gcp.GetSession(ctx)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	// Create KMS client with session credentials
	kmsClient, err := kms.NewKeyManagementClient(ctx, option.WithCredentials(ssn.Credentials))
	if err != nil {
		return fmt.Errorf("failed to create KMS client: %w", err)
	}
	defer kmsClient.Close()

	// Build the full KMS key resource name
	keyProjectID := projectID
	if kmsKey.ProjectID != "" {
		keyProjectID = kmsKey.ProjectID
	}
	resourceName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		keyProjectID, kmsKey.Location, kmsKey.KeyRing, kmsKey.Name)

	member := fmt.Sprintf("serviceAccount:%s", serviceAccountEmail)

	// Use exponential backoff for IAM policy updates to handle concurrent modifications
	backoff := wait.Backoff{
		Duration: 2 * time.Second,
		Factor:   2.0,
		Jitter:   1.0,
		Steps:    retryCount,
	}

	var lastErr error
	if waitErr := wait.ExponentialBackoffWithContext(ctx, backoff, func(ctx context.Context) (bool, error) {
		// Get current IAM policy
		req := &iampb.GetIamPolicyRequest{
			Resource: resourceName,
		}
		policy, err := kmsClient.GetIamPolicy(ctx, req)
		if err != nil {
			return false, fmt.Errorf("failed to get IAM policy for KMS key %s: %w", resourceName, err)
		}

		// Add the member to the role
		if err := addMemberToKMSPolicy(policy, role, member); err != nil {
			return false, fmt.Errorf("failed to add member to KMS policy: %w", err)
		}

		// Set the updated IAM policy
		setReq := &iampb.SetIamPolicyRequest{
			Resource: resourceName,
			Policy:   policy,
		}
		_, err = kmsClient.SetIamPolicy(ctx, setReq)
		if err != nil {
			// Check for concurrent modification errors (gRPC status codes)
			if st, ok := status.FromError(err); ok && (st.Code() == codes.Aborted || st.Code() == codes.FailedPrecondition) {
				lastErr = err
				logrus.Debugf("Concurrent KMS IAM policy changes, restarting read/modify/write")
				return false, nil
			}
			return false, fmt.Errorf("failed to set IAM policy for KMS key %s: %w", resourceName, err)
		}

		logrus.Infof("Successfully granted %s permission to %s on KMS key %s", role, serviceAccountEmail, resourceName)
		return true, nil
	}); waitErr != nil {
		if wait.Interrupted(waitErr) {
			logrus.Errorf("KMS IAM policy update exhausted retries after %d attempts for key %s, service account %s. Last error: %v",
				retryCount, resourceName, serviceAccountEmail, lastErr)
			return fmt.Errorf("failed to set KMS IAM policy for key %s after %d retries (last error: %w)",
				resourceName, retryCount, lastErr)
		}
		return waitErr
	}

	return nil
}

// addMemberToKMSPolicy adds a member to a KMS IAM policy binding for the specified role.
// If the member already exists, no changes are made. If the role doesn't exist in the policy,
// a new binding is created.
func addMemberToKMSPolicy(policy *iampb.Policy, role, member string) error {
	// Check if member already has the role
	for _, binding := range policy.Bindings {
		if binding.Role == role {
			for _, m := range binding.Members {
				if m == member {
					logrus.Debugf("Member %s already has role %s on KMS key", member, role)
					return nil
				}
			}
			// Role exists, add member to it
			binding.Members = append(binding.Members, member)
			logrus.Debugf("Added %s to existing role %s on KMS key", member, role)
			return nil
		}
	}

	// Role doesn't exist, create new binding
	policy.Bindings = append(policy.Bindings, &iampb.Binding{
		Role:    role,
		Members: []string{member},
	})
	logrus.Debugf("Created new role %s with member %s on KMS key", role, member)
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
