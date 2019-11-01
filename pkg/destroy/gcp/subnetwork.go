package gcp

import (
	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)

func (o *ClusterUninstaller) listSubNetworks() ([]string, error) {
	o.Logger.Debugf("Listing subnetworks")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Subnetworks.List(o.ProjectID, o.Region).Fields("items(name),nextPageToken").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(list *compute.SubnetworkList) error {
		for _, subNetwork := range list.Items {
			o.Logger.Debugf("Found subnetwork: %s", subNetwork.Name)
			result = append(result, subNetwork.Name)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list subnetworks")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteSubNetwork(name string) error {
	o.Logger.Debugf("Deleting subnetwork %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.Subnetworks.Delete(o.ProjectID, o.Region, name).RequestId(o.requestID("subnetwork", name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("subnetwork", name)
		return errors.Wrapf(err, "failed to delete subnetwork %s", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("subnetwork", name)
		return errors.Errorf("failed to delete subnetwork %s with error: %s", name, operationErrorMessage(op))
	}
	return nil
}

// destroySubNetworks removes all subnetwork resources that have a name prefixed
// with the cluster's infra ID
func (o *ClusterUninstaller) destroySubNetworks() error {
	subNetworks, err := o.listSubNetworks()
	if err != nil {
		return err
	}
	found := make([]string, 0, len(subNetworks))
	errs := []error{}
	for _, subNetwork := range subNetworks {
		found = append(found, subNetwork)
		err := o.deleteSubNetwork(subNetwork)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("subnetwork", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted subnetwork %s", item)
	}
	return aggregateError(errs, len(found))
}
