package gcp

import (
	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func (o *ClusterUninstaller) listTargetPools() ([]cloudResource, error) {
	return o.listTargetPoolsWithFilter("items(name),nextPageToken", o.clusterIDFilter(), nil)
}

// listTargetPoolsWithFilter lists target pools in the project. The field parameter allows
// specifying which fields to return. The filter parameter specifies a server-side filter for the
// GCP API (preferred). The filterFunc specifies a client-side filtering function for each TargetPool.
func (o *ClusterUninstaller) listTargetPoolsWithFilter(field string, filter string, filterFunc func(*compute.TargetPool) bool) ([]cloudResource, error) {
	o.Logger.Debugf("Listing target pools")
	result := []cloudResource{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.TargetPools.List(o.ProjectID, o.Region).Fields(googleapi.Field(field)).Filter(filter)
	err := req.Pages(ctx, func(list *compute.TargetPoolList) error {
		for _, targetPool := range list.Items {
			if filterFunc == nil || (filterFunc != nil && filterFunc(targetPool)) {
				o.Logger.Debugf("Found target pool: %s", targetPool.Name)
				result = append(result, cloudResource{
					key:      targetPool.Name,
					name:     targetPool.Name,
					typeName: "targetpool",
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

func (o *ClusterUninstaller) deleteTargetPool(item cloudResource) error {
	o.Logger.Debugf("Deleting target pool %s", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.TargetPools.Delete(o.ProjectID, o.Region, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Wrapf(err, "failed to delete target pool %s", item.name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Errorf("failed to delete route %s with error: %s", item.name, operationErrorMessage(op))
	}
	o.Logger.Infof("Deleted target pool %s", item.name)
	return nil
}

func (o *ClusterUninstaller) clearTargetPoolHealthChecks(name string) error {
	o.Logger.Debugf("Clearing target pool %s health checks", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	targetPool, err := o.computeSvc.TargetPools.Get(o.ProjectID, o.Region, name).Context(ctx).Do()
	if isNotFound(err) {
		return nil
	}
	if err != nil {
		return errors.Wrapf(err, "cannot retrieve target pool %s", name)
	}
	if len(targetPool.HealthChecks) == 0 {
		o.Logger.Debugf("Target pool %s has no health checks to clear", name)
		return nil
	}
	hcRemoveRequest := &compute.TargetPoolsRemoveHealthCheckRequest{}
	for _, hc := range targetPool.HealthChecks {
		hcRemoveRequest.HealthChecks = append(hcRemoveRequest.HealthChecks, &compute.HealthCheckReference{
			HealthCheck: hc,
		})
	}
	op, err := o.computeSvc.TargetPools.RemoveHealthCheck(o.ProjectID, o.Region, name, hcRemoveRequest).Context(ctx).RequestId(o.requestID("cleartargetpool", name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("cleartargetpool", name)
		return errors.Wrapf(err, "failed to clear target pool %s health checks", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("cleartargetpool", name)
		return errors.Errorf("failed to clear target pool %s health checks with error: %s", name, operationErrorMessage(op))
	}
	if op != nil && op.Status != "DONE" {
		return errors.Errorf("target pool pending to be cleared of health checks")
	}
	return nil
}

// destroyTargetPools removes target pools created for external load balancers that have
// a name that starts with the cluster infra ID. These are load balancers created by the
// installer or cluster operators.
func (o *ClusterUninstaller) destroyTargetPools() error {
	targetPools, err := o.listTargetPools()
	if err != nil {
		return err
	}
	found := cloudResources{}
	errs := []error{}
	for _, targetPool := range targetPools {
		found.insert(targetPool)
		err := o.deleteTargetPool(targetPool)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("targetpool", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted target pool %s", item.name)
	}
	return aggregateError(errs, len(found))
}
