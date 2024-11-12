package gcp

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

const (
	globalForwardingRuleResourceName = "forwardingrule"
	regionForwardingRuleResourceName = "regionforwardingrule"
)

func (o *ClusterUninstaller) listForwardingRules(ctx context.Context, typeName string) ([]cloudResource, error) {
	resources, err := o.listForwardingRulesWithFilter(ctx, typeName, "items(name,region,loadBalancingScheme),nextPageToken", o.clusterIDFilter())
	if err == nil {
		o.filteredResources.Insert(typeName)
	}
	return resources, err
}

// listForwardingRulesWithFilter lists forwarding rules in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listForwardingRulesWithFilter(ctx context.Context, typeName, fields, filter string) ([]cloudResource, error) {
	o.Logger.Debugf("Listing forwarding rules")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var err error
	var list *compute.ForwardingRuleList
	switch typeName {
	case globalForwardingRuleResourceName:
		list, err = o.computeSvc.GlobalForwardingRules.List(o.ProjectID).Filter(filter).Fields(googleapi.Field(fields)).Context(ctx).Do()
	case regionForwardingRuleResourceName:
		list, err = o.computeSvc.ForwardingRules.List(o.ProjectID, o.Region).Filter(filter).Fields(googleapi.Field(fields)).Context(ctx).Do()
	default:
		return nil, fmt.Errorf("invalid forwarding rule type %q", typeName)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list forwarding rules: %w", err)
	}

	result := []cloudResource{}
	for _, item := range list.Items {
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
			typeName: typeName,
			quota:    quota,
		})
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteForwardingRule(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting forwarding rule %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var err error
	var op *compute.Operation
	switch item.typeName {
	case globalForwardingRuleResourceName:
		op, err = o.computeSvc.GlobalForwardingRules.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	case regionForwardingRuleResourceName:
		op, err = o.computeSvc.ForwardingRules.Delete(o.ProjectID, o.Region, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	default:
		return fmt.Errorf("invalid forwarding rule type %q", item.typeName)
	}

	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return fmt.Errorf("failed to delete forwarding rule %s: %w", item.name, err)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return fmt.Errorf("failed to delete forwarding rule %s with error: %s: %w", item.name, operationErrorMessage(op), err)
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
	items := []cloudResource{}

	if cachedResources, ok := o.findCachedResources(globalForwardingRuleResourceName); !ok {
		found, err := o.listForwardingRules(ctx, globalForwardingRuleResourceName)
		if err != nil {
			return err
		}
		items = o.insertPendingItems(globalForwardingRuleResourceName, found)
	} else {
		items = append(items, cachedResources...)
	}

	if cachedResources, ok := o.findCachedResources(regionForwardingRuleResourceName); !ok {
		found, err := o.listForwardingRules(ctx, regionForwardingRuleResourceName)
		if err != nil {
			return err
		}
		items = append(items, o.insertPendingItems(regionForwardingRuleResourceName, found)...)
	} else {
		items = append(items, cachedResources...)
	}

	for _, item := range items {
		err := o.deleteForwardingRule(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}

	if items = o.getPendingItems(globalForwardingRuleResourceName); len(items) > 0 {
		return fmt.Errorf("%d global forwarding rules pending", len(items))
	}
	o.Logger.Debugf("Adding Destroyed Resource %s", globalForwardingRuleResourceName)
	o.destroyedResources.Insert(globalForwardingRuleResourceName)

	if items = o.getPendingItems(regionForwardingRuleResourceName); len(items) > 0 {
		return fmt.Errorf("%d regional forwarding rules pending", len(items))
	}
	o.Logger.Debugf("Adding Destroyed Resource %s", regionForwardingRuleResourceName)
	o.destroyedResources.Insert(regionForwardingRuleResourceName)

	return nil
}
