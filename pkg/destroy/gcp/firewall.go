package gcp

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

const (
	firewallResourceName = "firewall"
)

func (o *ClusterUninstaller) listFirewalls(ctx context.Context) ([]cloudResource, error) {
	// The firewall rules that the destroyer is searching for here include a
	// pattern before and after the cluster ID. Use a regular expression that allows
	// wildcard values before and after the cluster ID.
	return o.listFirewallsWithFilter(ctx, "items(name),nextPageToken", func(item string) bool {
		return strings.Contains(item, o.ClusterID)
	})
}

// listFirewallsWithFilter lists firewall rules in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listFirewallsWithFilter(ctx context.Context, fields string, filterFunc resourceFilterFunc) ([]cloudResource, error) {
	o.Logger.Debugf("Listing firewall rules")
	results := []cloudResource{}

	findFirewallRules := func(projectID string) ([]cloudResource, error) {
		ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
		defer cancel()
		result := []cloudResource{}
		req := o.computeSvc.Firewalls.List(projectID).Fields(googleapi.Field(fields))

		err := req.Pages(ctx, func(list *compute.FirewallList) error {
			for _, item := range list.Items {
				if filterFunc(item.Name) {
					o.Logger.Debugf("Found firewall rule: %s", item.Name)
					result = append(result, cloudResource{
						key:      item.Name,
						name:     item.Name,
						project:  projectID,
						typeName: firewallResourceName,
						quota: []gcp.QuotaUsage{{
							Metric: &gcp.Metric{
								Service: gcp.ServiceComputeEngineAPI,
								Limit:   "firewalls",
							},
							Amount: 1,
						}},
					})
				}
			}
			return nil
		})

		if err != nil {
			return nil, errors.Wrapf(err, "failed to list firewall rules for project %s", projectID)
		}
		return result, nil
	}

	findResults, err := findFirewallRules(o.ProjectID)
	if err != nil {
		return results, err
	}
	results = append(results, findResults...)

	if o.NetworkProjectID != "" {
		o.Logger.Debugf("Listing firewall rules for network project %s", o.NetworkProjectID)
		findResults, err := findFirewallRules(o.NetworkProjectID)
		if err != nil {
			return results, err
		}
		results = append(results, findResults...)
	}

	return results, nil
}

func (o *ClusterUninstaller) deleteFirewall(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting firewall rule %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	op, err := o.computeSvc.Firewalls.Delete(item.project, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	item.scope = global
	return o.handleOperation(ctx, op, err, item, "firewall")
}

// destroyFirewalls removes all firewall resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyFirewalls(ctx context.Context) error {
	found, err := o.listFirewalls(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems(firewallResourceName, found)
	for _, item := range items {
		err := o.deleteFirewall(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems(firewallResourceName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
