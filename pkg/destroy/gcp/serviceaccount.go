package gcp

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/api/iam/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/types/gcp"
)

// listServiceAccounts retrieves all service accounts with a display name prefixed with the cluster's
// infra ID. Filtering is done client side because the API doesn't offer filtering for service accounts.
func (o *ClusterUninstaller) listServiceAccounts(ctx context.Context) ([]cloudResource, error) {
	o.Logger.Debugf("Listing service accounts")

	result := []cloudResource{}
	sas, err := o.listClusterServiceAccount(ctx)
	if err != nil {
		errors.Wrapf(err, "failed to fetch service accounts for the cluster")
	}
	for _, item := range sas {
		o.Logger.Debugf("Found service account: %s", item.Name)
		result = append(result, cloudResource{
			key:      item.Name,
			name:     item.Name,
			url:      item.Email,
			typeName: "serviceaccount",
			quota: []gcp.QuotaUsage{{
				Metric: &gcp.Metric{
					Service: gcp.ServiceIAMAPI,
					Limit:   "quota/service-account-count",
				},
				Amount: 1,
			}},
		})
	}
	return result, nil
}

func (o *ClusterUninstaller) listClusterServiceAccount(ctx context.Context) ([]*iam.ServiceAccount, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []*iam.ServiceAccount{}
	req := o.iamSvc.Projects.ServiceAccounts.List(fmt.Sprintf("projects/%s", o.ProjectID)).Fields("accounts(name,displayName,email),nextPageToken")
	err := req.Pages(ctx, func(list *iam.ListServiceAccountsResponse) error {
		for idx, item := range list.Accounts {
			if o.isClusterResource(item.Email) || o.isClusterResource(item.DisplayName) {
				result = append(result, list.Accounts[idx])
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch service accounts")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteServiceAccount(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting service account %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	_, err := o.iamSvc.Projects.ServiceAccounts.Delete(item.name).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete service account %s", item.name)
	}
	o.deletePendingItems(item.typeName, []cloudResource{item})
	o.Logger.Infof("Deleted service account %s", item.name)
	return nil
}

// destroyServiceAccounts removes service accounts with a display name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyServiceAccounts(ctx context.Context) error {
	found, err := o.listServiceAccounts(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems("serviceaccount", found)
	if len(items) == 0 {
		return nil
	}

	// Remove service accounts from project policy
	policy, err := o.getProjectIAMPolicy(ctx)
	if err != nil {
		return err
	}
	emails := sets.NewString()
	for _, item := range items {
		emails.Insert(item.url)
	}
	if o.clearIAMPolicyBindings(policy, emails, o.Logger) {
		err = o.setProjectIAMPolicy(ctx, policy)
		if err != nil {
			o.errorTracker.suppressWarning("iampolicy", err, o.Logger)
			return errors.Errorf("%d items pending", len(items))
		}
		o.Logger.Infof("Deleted IAM project role bindings")
	}

	for _, item := range items {
		err := o.deleteServiceAccount(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems("serviceaccount"); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
