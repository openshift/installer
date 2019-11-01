package gcp

import (
	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)

func (o *ClusterUninstaller) listHTTPHealthChecks() ([]string, error) {
	o.Logger.Debugf("Listing HTTP health checks")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.HttpHealthChecks.List(o.ProjectID).Fields("items(name),nextPageToken").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(list *compute.HttpHealthCheckList) error {
		for _, healthCheck := range list.Items {
			o.Logger.Debugf("Found HTTP health check: %s", healthCheck.Name)
			result = append(result, healthCheck.Name)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list HTTP health checks")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteHTTPHealthCheck(name string, errorOnPending bool) error {
	o.Logger.Debugf("Deleting HTTP health check %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.HttpHealthChecks.Delete(o.ProjectID, name).RequestId(o.requestID("httphealthcheck", name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("httphealthcheck", name)
		return errors.Wrapf(err, "failed to delete HTTP health check %s", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("httphealthcheck", name)
		return errors.Errorf("failed to delete HTTP health check %s with error: %s", name, operationErrorMessage(op))
	}
	if errorOnPending && op != nil && op.Status != "DONE" {
		return errors.Errorf("deletion of HTTP health check %s is pending", name)
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
	found := make([]string, 0, len(healthChecks))
	errs := []error{}
	for _, healthCheck := range healthChecks {
		found = append(found, healthCheck)
		err := o.deleteHTTPHealthCheck(healthCheck, false)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("httphealthcheck", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted HTTP health check %s", item)
	}
	return aggregateError(errs, len(found))
}
