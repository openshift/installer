package gcp

import (
	"fmt"

	"github.com/pkg/errors"

	iam "google.golang.org/api/iam/v1"
)

// listServiceAccounts retrieves all service accounts with a display name prefixed with the cluster's
// infra ID. Filtering is done client side because the API doesn't offer filtering for service accounts.
func (o *ClusterUninstaller) listServiceAccounts() ([]cloudResource, error) {
	o.Logger.Debugf("Listing service accounts")

	result := []cloudResource{}
	sas, err := o.listClusterServiceAccount()
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
		})
	}
	return result, nil
}

func (o *ClusterUninstaller) listClusterServiceAccount() ([]*iam.ServiceAccount, error) {
	ctx, cancel := o.contextWithTimeout()
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

func (o *ClusterUninstaller) deleteServiceAccount(item cloudResource) error {
	o.Logger.Debugf("Deleting service account %s", item.name)
	ctx, cancel := o.contextWithTimeout()
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
func (o *ClusterUninstaller) destroyServiceAccounts() error {
	found, err := o.listServiceAccounts()
	if err != nil {
		return err
	}
	o.insertPendingItems("serviceaccount_binding", found) // store service accounts to remove project IAM binding

	items := o.insertPendingItems("serviceaccount", found)
	errs := []error{}
	for _, item := range items {
		err := o.deleteServiceAccount(item)
		if err != nil {
			errs = append(errs, err)
		}
	}
	items = o.getPendingItems("serviceaccount")
	return aggregateError(errs, len(items))
}
