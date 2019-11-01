package gcp

import (
	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)

func (o *ClusterUninstaller) listFirewalls() ([]string, error) {
	o.Logger.Debugf("Listing firewall rules")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Firewalls.List(o.ProjectID).Fields("items(name),nextPageToken").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(list *compute.FirewallList) error {
		for _, firewall := range list.Items {
			o.Logger.Debugf("Found firewall rule: %s", firewall.Name)
			result = append(result, firewall.Name)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list firewall rules")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteFirewall(name string, errorOnPending bool) error {
	o.Logger.Debugf("Deleting firewall rule %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.Firewalls.Delete(o.ProjectID, name).RequestId(o.requestID("firewall", name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("firewall", name)
		return errors.Wrapf(err, "failed to delete firewall %s", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("firewall", name)
		return errors.Errorf("failed to delete firewall %s with error: %s", name, operationErrorMessage(op))
	}
	if errorOnPending && op != nil && op.Status != "DONE" {
		return errors.Errorf("deletion of firewall %s is pending", name)
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
	found := make([]string, 0, len(firewalls))
	errs := []error{}
	for _, firewall := range firewalls {
		found = append(found, firewall)
		err := o.deleteFirewall(firewall, false)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("firewall", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted firewall %s", item)
	}
	return aggregateError(errs, len(found))
}
