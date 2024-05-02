package gcp

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

func (o *ClusterUninstaller) listHealthChecks(ctx context.Context, typeName string, listFunc healthCheckListFunc) ([]cloudResource, error) {
	return o.listHealthChecksWithFilter(ctx, typeName, "items(name),nextPageToken", o.clusterIDFilter(), listFunc)
}

// listHealthChecksWithFilter lists health checks in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listHealthChecksWithFilter(ctx context.Context, typeName, fields, filter string, listFunc healthCheckListFunc) ([]cloudResource, error) {
	o.Logger.Debugf("Listing health checks")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}
	list, err := listFunc(ctx, filter, fields)
	if err != nil {
		return nil, fmt.Errorf("failed to list health checks: %w", err)
	}

	for _, item := range list.Items {
		o.Logger.Debugf("Found health check: %s", item.Name)
		result = append(result, cloudResource{
			key:      item.Name,
			name:     item.Name,
			typeName: typeName,
			quota: []gcp.QuotaUsage{{
				Metric: &gcp.Metric{
					Service: gcp.ServiceComputeEngineAPI,
					Limit:   "health_checks",
				},
				Amount: 1,
			}},
		})
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteHealthCheck(ctx context.Context, item cloudResource, deleteFunc healthCheckDestroyFunc) error {
	o.Logger.Debugf("Deleting health check %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	op, err := deleteFunc(ctx, item)
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Wrapf(err, "failed to delete health check %s", item.name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Errorf("failed to delete health check %s with error: %s", item.name, operationErrorMessage(op))
	}
	if (err != nil && isNoOp(err)) || (op != nil && op.Status == "DONE") {
		o.resetRequestID(item.typeName, item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted health check %s", item.name)
	}
	return nil
}

// destroyHealthChecks removes all health check resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyHealthChecks(ctx context.Context) error {
	for _, hcd := range []healthCheckDestroyer{
		{
			itemTypeName: "healthcheck",
			destroyFunc:  o.healthCheckDelete,
			listFunc:     o.healthCheckList,
		},
		{
			itemTypeName: "regionHealthCheck",
			destroyFunc:  o.regionHealthCheckDelete,
			listFunc:     o.regionHealthCheckList,
		},
	} {
		found, err := o.listHealthChecks(ctx, hcd.itemTypeName, hcd.listFunc)
		if err != nil {
			return err
		}
		items := o.insertPendingItems(hcd.itemTypeName, found)
		for _, item := range items {
			err := o.deleteHealthCheck(ctx, item, hcd.destroyFunc)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}
		if items = o.getPendingItems(hcd.itemTypeName); len(items) > 0 {
			return errors.Errorf("%d items pending", len(items))
		}
	}
	return nil
}

type healthCheckListFunc func(ctx context.Context, filter, fields string) (*compute.HealthCheckList, error)
type healthCheckDestroyFunc func(ctx context.Context, item cloudResource) (*compute.Operation, error)
type healthCheckDestroyer struct {
	itemTypeName string
	destroyFunc  healthCheckDestroyFunc
	listFunc     healthCheckListFunc
}

func (o *ClusterUninstaller) healthCheckDelete(ctx context.Context, item cloudResource) (*compute.Operation, error) {
	return o.computeSvc.HealthChecks.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
}

func (o *ClusterUninstaller) healthCheckList(ctx context.Context, filter, fields string) (*compute.HealthCheckList, error) {
	return o.computeSvc.HealthChecks.List(o.ProjectID).Filter(filter).Fields(googleapi.Field(fields)).Context(ctx).Do()
}

func (o *ClusterUninstaller) regionHealthCheckDelete(ctx context.Context, item cloudResource) (*compute.Operation, error) {
	return o.computeSvc.RegionHealthChecks.Delete(o.ProjectID, o.Region, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
}

func (o *ClusterUninstaller) regionHealthCheckList(ctx context.Context, filter, fields string) (*compute.HealthCheckList, error) {
	return o.computeSvc.RegionHealthChecks.List(o.ProjectID, o.Region).Filter(filter).Fields(googleapi.Field(fields)).Context(ctx).Do()
}
