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
	globalForwardingRuleResource = "fowardingrule"
	regionForwardingRuleResource = "regionalforwardingrule"
)

func (o *ClusterUninstaller) listForwardingRules(ctx context.Context, typeName string) ([]cloudResource, error) {
	return o.listForwardingRulesWithFilter(ctx, typeName, "items(name,region,loadBalancingScheme),nextPageToken", o.isClusterResource)
}

// listForwardingRulesWithFilter lists forwarding rules in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listForwardingRulesWithFilter(ctx context.Context, typeName, fields string, filterFunc resourceFilterFunc) ([]cloudResource, error) {
	o.Logger.Debugf("Listing forwarding rules")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	result := []cloudResource{}

	pagesFunc := func(list *compute.ForwardingRuleList) error {
		for _, item := range list.Items {
			if filterFunc(item.Name) {
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
		}
		return nil
	}

	var err error
	switch typeName {
	case globalForwardingRuleResource:
		err = o.computeSvc.GlobalForwardingRules.List(o.ProjectID).Fields(googleapi.Field(fields)).Pages(ctx, pagesFunc)
	case regionForwardingRuleResource:
		err = o.computeSvc.ForwardingRules.List(o.ProjectID, o.Region).Fields(googleapi.Field(fields)).Pages(ctx, pagesFunc)
	default:
		return nil, fmt.Errorf("invalid forwarding rule type %q", typeName)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list forwarding rule: %w", err)
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
	case globalForwardingRuleResource:
		op, err = o.computeSvc.GlobalForwardingRules.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	case regionForwardingRuleResource:
		op, err = o.computeSvc.ForwardingRules.Delete(o.ProjectID, o.Region, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	default:
		return fmt.Errorf("invalid forwarding rule type %q", item.typeName)
	}

	if err = o.handleOperation(op, err, item, "forwarding rule"); err != nil {
		return err
	}
	return nil
}

// destroyForwardingRules removes all forwarding rules with a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyForwardingRules(ctx context.Context) error {
	forwardingRuleTypes := []string{globalForwardingRuleResource, regionForwardingRuleResource}
	for _, forwardingRuleType := range forwardingRuleTypes {
		found, err := o.listForwardingRules(ctx, forwardingRuleType)
		if err != nil {
			return err
		}
		items := o.insertPendingItems(forwardingRuleType, found)

		for _, item := range items {
			err := o.deleteForwardingRule(ctx, item)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}

		if items := o.getPendingItems(forwardingRuleType); len(items) > 0 {
			return fmt.Errorf("%d %s resources pending", len(items), forwardingRuleType)
		}
	}
	return nil
}
