package gcp

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	resourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	iam "google.golang.org/api/iam/v1"
)

// listServiceAccounts retrieves all service accounts with a display name prefixed with the cluster's
// infra ID. Filtering is done client side because the API doesn't offer filtering for service accounts.
func (o *ClusterUninstaller) listServiceAccounts() ([]string, error) {
	o.Logger.Debugf("Listing service accounts")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.iamSvc.Projects.ServiceAccounts.List(fmt.Sprintf("projects/%s", o.ProjectID)).Fields("accounts(name,email),nextPageToken")
	err := req.Pages(ctx, func(response *iam.ListServiceAccountsResponse) error {
		for _, account := range response.Accounts {
			if o.isClusterResource(account.Email) {
				o.Logger.Debugf("Found service account %s", account.Name)
				result = append(result, account.Name)
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch service accounts")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteServiceAccount(serviceAccount string) error {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	o.Logger.Debugf("Deleting service account %s", serviceAccount)
	_, err := o.iamSvc.Projects.ServiceAccounts.Delete(serviceAccount).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete service account %s", serviceAccount)
	}
	return nil
}

// destroyServiceAccounts removes service accounts with a display name that starts
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyServiceAccounts() error {
	serviceAccounts, err := o.listServiceAccounts()
	if err != nil {
		return err
	}
	errs := []error{}
	for _, serviceAccount := range serviceAccounts {
		err := o.deleteServiceAccount(serviceAccount)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deletedItems := o.setPendingItems("serviceaccount", serviceAccounts)
	for _, item := range deletedItems {
		o.Logger.Infof("Deleted service account %s", item)
	}
	return aggregateError(errs, len(serviceAccounts))
}

func (o *ClusterUninstaller) getProjectIAMPolicy() (*resourcemanager.Policy, error) {
	o.Logger.Debug("Fetching project IAM policy")
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

func clearIAMPolicyBindings(policy *resourcemanager.Policy, clusterID string, logger logrus.FieldLogger) bool {
	removedBindings := false
	for _, binding := range policy.Bindings {
		members := []string{}
		for _, member := range binding.Members {
			if strings.HasPrefix(member, fmt.Sprintf("serviceAccount:%s", clusterID)) {
				logger.Debugf("IAM: removing %s from role %s", member, binding.Role)
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
	if !clearIAMPolicyBindings(policy, o.ClusterID, o.Logger) {
		deletedPolicy := o.setPendingItems("iampolicy", []string{})
		if len(deletedPolicy) > 0 {
			o.Logger.Infof("Deleted IAM project role bindings")
		}
		return nil
	}
	o.setPendingItems("iampolicy", []string{"policy"})
	err = o.setProjectIAMPolicy(policy)
	return aggregateError([]error{err}, 1)
}
