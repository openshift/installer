package gcp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/iam/v1"
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

	// Unlike the pool, provider names are not ours to predict, so every provider
	// in the cluster-owned pool is claimed. The parent already scopes the search
	// to this cluster.
	parent := fmt.Sprintf("projects/%s/locations/global/workloadIdentityPools/%s", o.ProjectID, o.wifPoolID())
	err := o.iamSvc.Projects.Locations.WorkloadIdentityPools.Providers.List(parent).
		Pages(ctx, func(resp *iam.ListWorkloadIdentityPoolProvidersResponse) error {
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
			return nil
		})
	if err != nil {
		if isWIFUnavailable(err) {
			o.Logger.Debugf("Skipping WIF providers in %s: %v", parent, err)
			return nil, nil
		}
		return nil, fmt.Errorf("failed to list WIF providers: %w", err)
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

	// The pool name is fully determined by the cluster ID, so fetch it directly
	// rather than paging through every pool in the project.
	name := fmt.Sprintf("projects/%s/locations/global/workloadIdentityPools/%s", o.ProjectID, o.wifPoolID())
	pool, err := o.iamSvc.Projects.Locations.WorkloadIdentityPools.Get(name).Context(ctx).Do()
	if err != nil {
		if isWIFUnavailable(err) {
			o.Logger.Debugf("Skipping WIF pool %s: %v", name, err)
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get WIF pool %s: %w", name, err)
	}

	// A deleted pool lingers in a soft-deleted state until GCP purges it.
	if pool.State == "DELETED" {
		return nil, nil
	}

	o.Logger.Debugf("Found WIF pool: %s", pool.Name)
	return []cloudResource{{
		key:      pool.Name,
		name:     pool.Name,
		typeName: "wifpool",
	}}, nil
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
