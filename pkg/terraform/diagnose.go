package terraform

import (
	"regexp"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/diagnostics"
)

// Diagnose accepts an error from terraform runs and tries to diagnose the
// underlying cause.
func Diagnose(message string) error {
	for _, cand := range conditions {
		if cand.match.MatchString(message) {
			return &diagnostics.Err{
				Source:  "Infrastructure Provider",
				Reason:  cand.reason,
				Message: cand.message,
			}
		}
	}

	return errors.New("failed to complete the change")
}

type condition struct {
	match *regexp.Regexp

	reason  string
	message string
}

// conditions is a list matches for the error string from terraform.
// specific on the top, generic matches on the bottom.
var conditions = []condition{{
	match: regexp.MustCompile(`Error: Error creating Blob .*: Error copy/waiting`),

	reason:  "Timeout",
	message: `Copying the VHD to user environment was too slow, and timeout was reached for the success.`,
}, {
	match: regexp.MustCompile(`Error: Error Creating/Updating Subnet .*: network.SubnetsClient#CreateOrUpdate: .* Code="AnotherOperationInProgress" Message="Another operation on this or dependent resource is in progress`),

	reason:  "AzureMultiOperationFailure",
	message: `Creating Subnets failed because Azure could not process multiple operations.`,
}, {
	match: regexp.MustCompile(`Error: Error Creating/Updating Public IP .*: network.PublicIPAddressesClient#CreateOrUpdate: .* Code="PublicIPCountLimitReached" Message="Cannot create more than .* public IP addresses for this subscription in this region`),

	reason:  "AzureQuotaLimitExceeded",
	message: `Service limits exceeded for Public IPs in the the subscriptions for the region. Requesting increase in quota should fix the error.`,
}, {
	match: regexp.MustCompile(`Error: compute\.VirtualMachinesClient#CreateOrUpdate: .* Code="OperationNotAllowed" Message="Operation could not be completed as it results in exceeding approved Total Regional Cores quota`),

	reason:  "AzureQuotaLimitExceeded",
	message: `Service limits exceeded for Virtual Machine cores in the the subscriptions for the region. Requesting increase in quota should fix the error.`,
}, {
	match: regexp.MustCompile(`Error: Code="OSProvisioningTimedOut"`),

	reason:  "AzureVirtualMachineFailure",
	message: `Some virtual machines failed to provision in alloted time. Virtual machines can fail to provision if the bootstap virtual machine has failing services.`,
}, {
	match: regexp.MustCompile(`Status=404 Code="ResourceGroupNotFound"`),

	reason:  "AzureEventualConsistencyFailure",
	message: `Failed to find a resource that was recently created usualy caused by Azure's eventual consistency delays.`,
}, {
	match: regexp.MustCompile(`Error: Error applying IAM policy to project .*: Too many conflicts`),

	reason:  "GCPTooManyIAMUpdatesInFlight",
	message: `There are a lot of IAM updates to the project in flight. Failed after reaching a limit of read-modify-write on conflict backoffs.`,
}, {
	match: regexp.MustCompile(`Error: .*: googleapi: Error 503: .*, backendError`),

	reason:  "GCPBackendInternalError",
	message: `GCP is experiencing backend service interuptions. Please try again or contact Google Support`,
}, {
	match: regexp.MustCompile(`Error: Error waiting for instance to create: Internal error`),

	reason:  "GCPComputeBackendTimeout",
	message: `GCP is experiencing backend service interuptions, the compute instance failed to create in reasonable time.`,
}}
