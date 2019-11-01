package gcp

import (
	"fmt"

	"github.com/pkg/errors"

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
