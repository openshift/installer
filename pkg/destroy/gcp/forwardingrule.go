package gcp

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

func (o *ClusterUninstaller) listForwardingRules(ctx context.Context) ([]cloudResource, error) {
	return o.listForwardingRulesWithFilter(ctx, "items(name,region,loadBalancingScheme),nextPageToken", o.clusterIDFilter(), nil)
}

// listForwardingRulesWithFilter lists forwarding rules in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listForwardingRulesWithFilter(ctx context.Context, fields string, filter string, filterFunc func(*compute.ForwardingRule) bool) ([]cloudResource, error) {
	o.Logger.Debugf("Listing forwarding rules")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
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
				var quota []gcp.QuotaUsage
				switch item.LoadBalancingScheme {
				case "EXTERNAL":
					quota = []gcp.QuotaUsage{{
						Metric: &gcp.Metric{
							Service: gcp.ServiceComputeEngineAPI,
							Limit:   "external_network_lb_forwarding_rules",
							Dimensions: map[string]string{
								"region": getNameFromURL("regions", item.Region),
							},
						},
						Amount: 1,
					}}
				}
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: "forwardingrule",
					quota:    quota,
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

func (o *ClusterUninstaller) deleteForwardingRule(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting forwarding rule %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
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
func (o *ClusterUninstaller) destroyForwardingRules(ctx context.Context) error {
	found, err := o.listForwardingRules(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems("forwardingrule", found)
	for _, item := range items {
		err := o.deleteForwardingRule(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems("forwardingrule"); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
