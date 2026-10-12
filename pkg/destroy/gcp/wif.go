package gcp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/api/googleapi"
)

// wifPoolID is the ID of the workload identity pool owned by the cluster. WIF
// resources carry no cluster label, so they are discovered by name convention.
func (o *ClusterUninstaller) wifPoolID() string {
	return fmt.Sprintf("%s-wif-pool", o.ClusterID)
}

// isWIFUnavailable reports whether an error means there are no cluster-owned WIF
// resources to look at: either nothing exists, or the credentials cannot see
// workload identity pools at all. Clusters that never used WIF are the common
// case, so neither should fail the destroy.
func isWIFUnavailable(err error) bool {
	if isNoOp(err) {
		return true
	}
	var ae *googleapi.Error
	return errors.As(err, &ae) && (ae.Code == http.StatusForbidden || ae.Code == http.StatusUnauthorized)
}

func (o *ClusterUninstaller) listWIFProviders(ctx context.Context) ([]cloudResource, error) {
	o.Logger.Debugf("Listing WIF providers")
	result := []cloudResource{}

	parent := fmt.Sprintf("projects/%s/locations/global/workloadIdentityPools/%s", o.ProjectID, o.wifPoolID())
	resp, err := o.iamSvc.Projects.Locations.WorkloadIdentityPools.Providers.List(parent).Context(ctx).Do()
	if err != nil {
		if isWIFUnavailable(err) {
			o.Logger.Debugf("Skipping WIF providers: %v", err)
			return nil, nil
		}
		return nil, fmt.Errorf("failed to list WIF providers: %w", err)
	}

	for _, provider := range resp.WorkloadIdentityPoolProviders {
		if provider.State == "DELETED" {
			continue
		}
		o.Logger.Debugf("Found WIF provider: %s", provider.Name)
		result = append(result, cloudResource{
			key:      provider.Name,
			name:     provider.Name,
			typeName: "wifprovider",
		})
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteWIFProvider(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting WIF provider %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	op, err := o.iamSvc.Projects.Locations.WorkloadIdentityPools.Providers.Delete(item.name).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return fmt.Errorf("failed to delete WIF provider %s: %w", item.name, err)
	}
	if op != nil && !op.Done {
		o.Logger.Debugf("WIF provider deletion in progress: %s", op.Name)
		return fmt.Errorf("WIF provider %s deletion still in progress", item.name)
	}
	o.deletePendingItems(item.typeName, []cloudResource{item})
	o.Logger.Infof("Deleted WIF provider %s", item.name)
	return nil
}

func (o *ClusterUninstaller) destroyWIFProviders(ctx context.Context) error {
	found, err := o.listWIFProviders(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems("wifprovider", found)
	for _, item := range items {
		if err := o.deleteWIFProvider(ctx, item); err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems("wifprovider"); len(items) > 0 {
		return fmt.Errorf("%d items pending", len(items))
	}
	return nil
}

func (o *ClusterUninstaller) listWIFPools(ctx context.Context) ([]cloudResource, error) {
	o.Logger.Debugf("Listing WIF pools")
	result := []cloudResource{}

	parent := fmt.Sprintf("projects/%s/locations/global", o.ProjectID)
	resp, err := o.iamSvc.Projects.Locations.WorkloadIdentityPools.List(parent).Context(ctx).Do()
	if err != nil {
		if isWIFUnavailable(err) {
			o.Logger.Debugf("Skipping WIF pools: %v", err)
			return nil, nil
		}
		return nil, fmt.Errorf("failed to list WIF pools: %w", err)
	}

	expectedPoolID := o.wifPoolID()
	for _, pool := range resp.WorkloadIdentityPools {
		if pool.State == "DELETED" {
			continue
		}
		poolID := pool.Name[strings.LastIndex(pool.Name, "/")+1:]
		if poolID == expectedPoolID {
			o.Logger.Debugf("Found WIF pool: %s", pool.Name)
			result = append(result, cloudResource{
				key:      pool.Name,
				name:     pool.Name,
				typeName: "wifpool",
			})
		}
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteWIFPool(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting WIF pool %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	op, err := o.iamSvc.Projects.Locations.WorkloadIdentityPools.Delete(item.name).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return fmt.Errorf("failed to delete WIF pool %s: %w", item.name, err)
	}
	if op != nil && !op.Done {
		o.Logger.Debugf("WIF pool deletion in progress: %s", op.Name)
		return fmt.Errorf("WIF pool %s deletion still in progress", item.name)
	}
	o.deletePendingItems(item.typeName, []cloudResource{item})
	o.Logger.Infof("Deleted WIF pool %s", item.name)
	return nil
}

func (o *ClusterUninstaller) destroyWIFPools(ctx context.Context) error {
	found, err := o.listWIFPools(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems("wifpool", found)
	for _, item := range items {
		if err := o.deleteWIFPool(ctx, item); err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems("wifpool"); len(items) > 0 {
		return fmt.Errorf("%d items pending", len(items))
	}
	return nil
}
