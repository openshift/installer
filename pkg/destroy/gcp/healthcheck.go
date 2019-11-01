package gcp

import (
	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)

func (o *ClusterUninstaller) listHealthChecks() ([]string, error) {
	o.Logger.Debugf("Listing health checks")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.HealthChecks.List(o.ProjectID).Fields("items(name),nextPageToken").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(list *compute.HealthCheckList) error {
		for _, healthCheck := range list.Items {
			o.Logger.Debugf("Found health check: %s", healthCheck.Name)
			result = append(result, healthCheck.Name)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list health checks")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteHealthCheck(name string, errorOnPending bool) error {
	o.Logger.Debugf("Deleting health check %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.HealthChecks.Delete(o.ProjectID, name).RequestId(o.requestID("healthcheck", name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("healthcheck", name)
		return errors.Wrapf(err, "failed to delete health check %s", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("healthcheck", name)
		return errors.Errorf("failed to delete health check %s with error: %s", name, operationErrorMessage(op))
	}
	if errorOnPending && op != nil && op.Status != "DONE" {
		return errors.Errorf("deletion of health check %s is pending", name)
	}
	return nil
}

// destroyHealthChecks removes all health check resources that have a name prefixed
// with the cluster's infra ID
func (o *ClusterUninstaller) destroyHealthChecks() error {
	healthChecks, err := o.listHealthChecks()
	if err != nil {
		return err
	}
	found := make([]string, 0, len(healthChecks))
	errs := []error{}
	for _, healthCheck := range healthChecks {
		found = append(found, healthCheck)
		err := o.deleteHealthCheck(healthCheck, false)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("healthcheck", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted health check %s", item)
	}
	return aggregateError(errs, len(found))
}
