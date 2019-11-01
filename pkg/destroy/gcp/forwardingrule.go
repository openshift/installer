package gcp

import (
	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)

func (o *ClusterUninstaller) listForwardingRules() ([]cloudResource, error) {
	o.Logger.Debugf("Listing forwarding rules")
	result := []cloudResource{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.ForwardingRules.List(o.ProjectID, o.Region).Fields("items(name),nextPageToken").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(list *compute.ForwardingRuleList) error {
		for _, forwardingRule := range list.Items {
			o.Logger.Debugf("Found forwarding rule: %s", forwardingRule.Name)
			result = append(result, cloudResource{
				key:      forwardingRule.Name,
				name:     forwardingRule.Name,
				typeName: "forwardingrule",
			})
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list forwarding rules")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteForwardingRule(item cloudResource, errorOnPending bool) error {
	o.Logger.Debugf("Deleting forwarding rule %s", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.ForwardingRules.Delete(o.ProjectID, o.Region, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Wrapf(err, "failed to delete forwarding rule %s", item.name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Errorf("failed to delete forwarding rule %s with error: %s", item.name, operationErrorMessage(op))
	}
	if errorOnPending && op != nil && op.Status != "DONE" {
		return errors.Errorf("deletion of forwarding rule %s is pending", item.name)
	}
	return nil
}

// destroyForwardingRules removes all forwarding rules with a name prefixed with the
// cluster's infra ID.
func (o *ClusterUninstaller) destroyForwardingRules() error {
	forwardingRules, err := o.listForwardingRules()
	if err != nil {
		return err
	}
	found := cloudResources{}
	errs := []error{}
	for _, forwardingRule := range forwardingRules {
		found.insert(forwardingRule)
		err := o.deleteForwardingRule(forwardingRule, false)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("forwardingrule", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted forwarding rule %s", item.name)
	}
	return aggregateError(errs, len(found))
}
