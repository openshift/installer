package gcp

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

const (
	targetPoolResourceName = "targetpool"
)

type targetPoolFilterFunc func(pool *compute.TargetPool) bool

func (o *ClusterUninstaller) listTargetPools(ctx context.Context) ([]cloudResource, error) {
	return o.listTargetPoolsWithFilter(ctx, "items(name),nextPageToken", func(item *compute.TargetPool) bool {
		return o.isClusterResource(item.Name)
	})
}

// listTargetPoolsWithFilter lists target pools in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listTargetPoolsWithFilter(ctx context.Context, fields string, filterFunc targetPoolFilterFunc) ([]cloudResource, error) {
	o.Logger.Debugf("Listing target pools")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.TargetPools.List(o.ProjectID, o.Region).Fields(googleapi.Field(fields))

	err := req.Pages(ctx, func(list *compute.TargetPoolList) error {
		for _, item := range list.Items {
			if filterFunc(item) {
				o.Logger.Debugf("Found target pool: %s", item.Name)
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: targetPoolResourceName,
					quota: []gcp.QuotaUsage{{
						Metric: &gcp.Metric{
							Service: gcp.ServiceComputeEngineAPI,
							Limit:   "target_pools",
						},
						Amount: 1,
					}},
				})
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list target pools")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteTargetPool(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting target pool %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	op, err := o.computeSvc.TargetPools.Delete(o.ProjectID, o.Region, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	if err = o.handleOperation(op, err, item, "target pool"); err != nil {
		return err
	}
	return nil
}

// destroyTargetPools removes target pools resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyTargetPools(ctx context.Context) error {
	found, err := o.listTargetPools(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems(targetPoolResourceName, found)
	for _, item := range items {
		err := o.deleteTargetPool(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems(targetPoolResourceName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
