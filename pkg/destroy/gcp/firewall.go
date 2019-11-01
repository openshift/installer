package gcp

import (
	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)

func (o *ClusterUninstaller) listFirewalls() ([]cloudResource, error) {
	o.Logger.Debugf("Listing firewall rules")
	result := []cloudResource{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Firewalls.List(o.ProjectID).Fields("items(name),nextPageToken").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(list *compute.FirewallList) error {
		for _, firewall := range list.Items {
			o.Logger.Debugf("Found firewall rule: %s", firewall.Name)
			result = append(result, cloudResource{
				key:      firewall.Name,
				name:     firewall.Name,
				typeName: "firewallrule",
			})
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list firewall rules")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteFirewall(item cloudResource, errorOnPending bool) error {
	o.Logger.Debugf("Deleting firewall rule %s", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.Firewalls.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Wrapf(err, "failed to delete firewall %s", item.name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Errorf("failed to delete firewall %s with error: %s", item.name, operationErrorMessage(op))
	}
	if errorOnPending && op != nil && op.Status != "DONE" {
		return errors.Errorf("deletion of firewall %s is pending", item.name)
	}
	return nil
}

// destroyFirewalls removes all firewall resources that have a name prefixed
// with the cluster's infra ID
func (o *ClusterUninstaller) destroyFirewalls() error {
	firewalls, err := o.listFirewalls()
	if err != nil {
		return err
	}
	found := cloudResources{}
	errs := []error{}
	for _, firewall := range firewalls {
		found.insert(firewall)
		err := o.deleteFirewall(firewall, false)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("firewall", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted firewall %s", item.name)
	}
	return aggregateError(errs, len(found))
}
