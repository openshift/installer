package gcp

import (
	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)

func (o *ClusterUninstaller) listHTTPHealthChecks() ([]cloudResource, error) {
	o.Logger.Debugf("Listing HTTP health checks")
	result := []cloudResource{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.HttpHealthChecks.List(o.ProjectID).Fields("items(name),nextPageToken").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(list *compute.HttpHealthCheckList) error {
		for _, healthCheck := range list.Items {
			o.Logger.Debugf("Found HTTP health check: %s", healthCheck.Name)
			result = append(result, cloudResource{
				key:      healthCheck.Name,
				name:     healthCheck.Name,
				typeName: "httphealthcheck",
			})
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list HTTP health checks")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteHTTPHealthCheck(item cloudResource, errorOnPending bool) error {
	o.Logger.Debugf("Deleting HTTP health check %s", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.HttpHealthChecks.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Wrapf(err, "failed to delete HTTP health check %s", item.name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Errorf("failed to delete HTTP health check %s with error: %s", item.name, operationErrorMessage(op))
	}
	if errorOnPending && op != nil && op.Status != "DONE" {
		return errors.Errorf("deletion of HTTP health check %s is pending", item.name)
	}
	return nil
}

// destroyHTTPHealthChecks removes all HTTP health check resources that have a name prefixed
// with the cluster's infra ID
func (o *ClusterUninstaller) destroyHTTPHealthChecks() error {
	healthChecks, err := o.listHTTPHealthChecks()
	if err != nil {
		return err
	}
	found := cloudResources{}
	errs := []error{}
	for _, healthCheck := range healthChecks {
		found.insert(healthCheck)
		err := o.deleteHTTPHealthCheck(healthCheck, false)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("httphealthcheck", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted HTTP health check %s", item.name)
	}
	return aggregateError(errs, len(found))
}
