package gcp

import (
	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func (o *ClusterUninstaller) listForwardingRules() ([]cloudResource, error) {
	return o.listForwardingRulesWithFilter("items(name),nextPageToken", o.clusterIDFilter(), nil)
}

// listForwardingRulesWithFilter lists forwarding rules in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listForwardingRulesWithFilter(fields string, filter string, filterFunc func(*compute.ForwardingRule) bool) ([]cloudResource, error) {
	o.Logger.Debugf("Listing forwarding rules")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.ForwardingRules.List(o.ProjectID, o.Region).Fields(googleapi.Field(fields))
	if len(filter) > 0 {
		req = req.Filter(filter)
	}
	err := req.Pages(ctx, func(list *compute.ForwardingRuleList) error {
		for _, item := range list.Items {
			if filterFunc == nil || filterFunc != nil && filterFunc(item) {
				o.Logger.Debugf("Found forwarding rule: %s", item.Name)
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: "forwardingrule",
				})
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list forwarding rules")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteForwardingRule(item cloudResource) error {
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
	if (err != nil && isNoOp(err)) || (op != nil && op.Status == "DONE") {
		o.resetRequestID(item.typeName, item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted forwarding rule %s", item.name)
	}
	return nil
}

// destroyForwardingRules removes all forwarding rules with a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyForwardingRules() error {
	found, err := o.listForwardingRules()
	if err != nil {
		return err
	}
	items := o.insertPendingItems("forwardingrule", found)
	errs := []error{}
	for _, item := range items {
		err := o.deleteForwardingRule(item)
		if err != nil {
			errs = append(errs, err)
		}
	}
	items = o.getPendingItems("forwardingrule")
	return aggregateError(errs, len(items))
}
