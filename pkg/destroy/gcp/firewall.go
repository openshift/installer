package gcp

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

func (o *ClusterUninstaller) listFirewalls(ctx context.Context) ([]cloudResource, error) {
	return o.listFirewallsWithFilter(ctx, "items(name),nextPageToken", o.clusterIDFilter(), nil)
}

// listFirewallsWithFilter lists firewall rules in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listFirewallsWithFilter(ctx context.Context, fields string, filter string, filterFunc func(*compute.Firewall) bool) ([]cloudResource, error) {
	o.Logger.Debugf("Listing firewall rules")
	results := []cloudResource{}

	findFirewallRules := func(projectID string) ([]cloudResource, error) {
		ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
		defer cancel()
		result := []cloudResource{}
		req := o.computeSvc.Firewalls.List(projectID).Fields(googleapi.Field(fields))
		if len(filter) > 0 {
			req = req.Filter(filter)
		}
		err := req.Pages(ctx, func(list *compute.FirewallList) error {
			for _, item := range list.Items {
				if filterFunc == nil || filterFunc != nil && filterFunc(item) {
					o.Logger.Debugf("Found firewall rule: %s", item.Name)
					result = append(result, cloudResource{
						key:      item.Name,
						name:     item.Name,
						project:  projectID,
						typeName: "firewall",
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
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Wrapf(err, "failed to delete firewall %s", item.name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Errorf("failed to delete firewall %s with error: %s", item.name, operationErrorMessage(op))
	}
	if (err != nil && isNoOp(err)) || (op != nil && op.Status == "DONE") {
		o.resetRequestID(item.typeName, item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted firewall rule %s", item.name)
	}
	return nil
}

// destroyFirewalls removes all firewall resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyFirewalls(ctx context.Context) error {
	found, err := o.listFirewalls(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems("firewall", found)
	for _, item := range items {
		err := o.deleteFirewall(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems("firewall"); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
