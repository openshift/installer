package gcp

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

func (o *ClusterUninstaller) listForwardingRules(ctx context.Context, scope resourceScope) ([]cloudResource, error) {
	return o.listForwardingRulesWithFilter(ctx, "items(name,region,loadBalancingScheme),nextPageToken", o.clusterIDFilter(), nil, scope)
}

func createForwardingRuleResources(filterFunc func(*compute.ForwardingRule) bool, list *compute.ForwardingRuleList) []cloudResource {
	result := []cloudResource{}

	for _, item := range list.Items {
		if filterFunc == nil || filterFunc(item) {
			logrus.Debugf("Found forwarding rule: %s", item.Name)
			var quota []gcp.QuotaUsage
			if item.LoadBalancingScheme == "EXTERNAL" {
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

	return result
}

// listForwardingRulesWithFilter lists forwarding rules in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listForwardingRulesWithFilter(ctx context.Context, fields string, filter string, filterFunc func(*compute.ForwardingRule) bool, scope resourceScope) ([]cloudResource, error) {
	o.Logger.Debugf("Listing %s forwarding rules", scope)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	result := []cloudResource{}

	if scope == gcpGlobalResource {
		req := o.computeSvc.GlobalForwardingRules.List(o.ProjectID).Fields(googleapi.Field(fields))
		if len(filter) > 0 {
			req = req.Filter(filter)
		}
		err := req.Pages(ctx, func(list *compute.ForwardingRuleList) error {
			result = append(result, createForwardingRuleResources(filterFunc, list)...)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list global forwarding rules: %w", err)
		}

		return result, nil
	}

	// Regional forwarding rules
	req := o.computeSvc.ForwardingRules.List(o.ProjectID, o.Region).Fields(googleapi.Field(fields))
	if len(filter) > 0 {
		req = req.Filter(filter)
	}
	err := req.Pages(ctx, func(list *compute.ForwardingRuleList) error {
		result = append(result, createForwardingRuleResources(filterFunc, list)...)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list regional forwarding rules: %w", err)
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteForwardingRule(ctx context.Context, item cloudResource, scope resourceScope) error {
	o.Logger.Debugf("Deleting %s forwarding rule %s", scope, item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var op *compute.Operation
	var err error
	if scope == gcpGlobalResource {
		op, err = o.computeSvc.GlobalForwardingRules.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	} else {
		op, err = o.computeSvc.ForwardingRules.Delete(o.ProjectID, o.Region, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	}

	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return fmt.Errorf("failed to delete forwarding rule %s: %w", item.name, err)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return fmt.Errorf("failed to delete forwarding rule %s with error: %s: %w", item.name, operationErrorMessage(op), err)
	}
	if op != nil && op.Status == "DONE" {
		o.resetRequestID(item.typeName, item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted forwarding rule %s", item.name)
	}
	return nil
}

// destroyForwardingRules removes all forwarding rules with a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyForwardingRules(ctx context.Context) error {
	for _, scope := range []resourceScope{gcpRegionalResource, gcpGlobalResource} {
		found, err := o.listForwardingRules(ctx, scope)
		if err != nil {
			return fmt.Errorf("failed to list forwarding rules: %w", err)
		}
		items := o.insertPendingItems("forwardingrule", found)
		for _, item := range items {
			if err := o.deleteForwardingRule(ctx, item, scope); err != nil {
				o.Logger.Errorf("error deleting forwarding rule %s: %w", item.name, err)
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}
		if items = o.getPendingItems("forwardingrule"); len(items) > 0 {
			for _, item := range items {
				if err := o.deleteForwardingRule(ctx, item, scope); err != nil {
					return fmt.Errorf("error deleting pending forwarding rule %s: %w", item.name, err)
				}
			}
		}
	}
	return nil
}
