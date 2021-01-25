package gcp

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	resourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

func (o *ClusterUninstaller) getProjectIAMPolicy() (*resourcemanager.Policy, error) {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := &resourcemanager.GetIamPolicyRequest{}
	policy, err := o.rmSvc.Projects.GetIamPolicy(o.ProjectID, req).Context(ctx).Do()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch project IAM policy")
	}
	return policy, nil
}

func (o *ClusterUninstaller) setProjectIAMPolicy(policy *resourcemanager.Policy) error {
	o.Logger.Debug("Setting project IAM policy")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := &resourcemanager.SetIamPolicyRequest{Policy: policy}
	_, err := o.rmSvc.Projects.SetIamPolicy(o.ProjectID, req).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "failed to set project IAM policy")
	}
	return nil
}

// removeIAMPolicyBindings reads through the IAM policy and updates a local copy to remove any
// roles with a service account prefixed with the clusterID. The list of pending items for deletion
// is updated to reflect the current state of the policy.
// returns true if the policy was updated to remove a role.
func (o *ClusterUninstaller) removeIAMPolicyBindings(policy *resourcemanager.Policy, logger logrus.FieldLogger) bool {
	removedBindings := false

	// Clear any pending role bindings so they can be updated to reflect the current IAM policy,
	// which is the source of truth. We shouldn't expect any rolebindings to be pending unless
	// there was an error when writing the updated policy on a previous run.
	o.clearPendingRoleBindings(false)

	logger.Debug("Listing IAM role bindings")

	for _, binding := range policy.Bindings {
		members := []string{}
		for _, member := range binding.Members {
			email := policyMemberToEmail(member)
			if o.isClusterResource(email) {
				logger.Debugf("IAM: updating local policy to remove %s from role %s", member, binding.Role)
				bindingName := fmt.Sprintf("%s-%s", member, binding.Role)
				o.insertPendingItems("iamrolebindings", []cloudResource{{key: bindingName, name: bindingName, typeName: "iamrolebindings"}})
				removedBindings = true
				continue
			}
			members = append(members, member)
		}
		binding.Members = members
	}
	return removedBindings
}

// destroyIAMPolicyBindings removes any role bindings from the project policy to
// service accounts that start with the cluster's infra ID.
func (o *ClusterUninstaller) destroyIAMPolicyBindings() error {
	policy, err := o.getProjectIAMPolicy()
	if err != nil {
		return err
	}

	if o.removeIAMPolicyBindings(policy, o.Logger) {
		if err = o.setProjectIAMPolicy(policy); err != nil {
			pendingPolicy := o.getPendingItems("iamrolebindings")
			errMsgPendingItems := fmt.Sprintf("unable to update IAM policy to remove %d pending roles", len(pendingPolicy))
			return errors.Wrapf(err, errMsgPendingItems)
		}

		// Clear pending items, because updating policy was successful.
		o.clearPendingRoleBindings(true)
		return nil
	}

	// No role bindings were found in the policy.
	return nil
}

// policyMemberToEmail takes member of IAM policy binding and converts it to service account email.
// https://cloud.google.com/iam/docs/reference/rest/v1/Policy#Binding
// see members[]
func policyMemberToEmail(member string) string {
	email := strings.TrimPrefix(strings.TrimPrefix(member, "deleted:"), "serviceAccount:")
	if idx := strings.Index(email, "?uid"); idx != -1 {
		email = email[:idx]
	}
	return email
}

// clearPendingRoleBindings removes all currently pending role bindings.
// expected toggles whether a debugging statement is displayed if items are removed from the queue.
func (o *ClusterUninstaller) clearPendingRoleBindings(expected bool) {
	pendingPolicy := o.getPendingItems("iamrolebindings")
	if len(pendingPolicy) > 0 {
		if !expected {
			o.Logger.Debugf("Found %d leftover IAM rolebindings when clearing pending items", len(pendingPolicy))
		}
		o.deletePendingItems("iamrolebindings", pendingPolicy)
	}
}
