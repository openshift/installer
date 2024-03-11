package gcp

import (
	"context"
	"fmt"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

func (o *ClusterUninstaller) listTargetTCPProxies(ctx context.Context) ([]cloudResource, error) {
	return o.listTargetTCPProxiesWithFilter(ctx, "items(name),nextPageToken", o.clusterIDFilter(), nil)
}

// listTargetTCPProxiesWithFilter lists target TCP Proxies in the project that satisfy the filter criteria.
func (o *ClusterUninstaller) listTargetTCPProxiesWithFilter(ctx context.Context, fields string, filter string, filterFunc func(list *compute.TargetTcpProxy) bool) ([]cloudResource, error) {
	o.Logger.Debugf("Listing target tcp proxies")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.TargetTcpProxies.List(o.ProjectID).Fields(googleapi.Field(fields))
	if len(filter) > 0 {
		req = req.Filter(filter)
	}
	err := req.Pages(ctx, func(list *compute.TargetTcpProxyList) error {
		for _, item := range list.Items {
			if filterFunc == nil || (filterFunc != nil && filterFunc(item)) {
				o.Logger.Debugf("Found target TCP proxy: %s", item.Name)
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: "targettcpproxy",
					quota: []gcp.QuotaUsage{{
						Metric: &gcp.Metric{
							Service: gcp.ServiceComputeEngineAPI,
							Limit:   "target_tcp_proxy",
						},
						Amount: 1,
					}},
				})
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list target tcp proxies: %w", err)
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteTargetTCPProxy(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting target TCP Proxies %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	op, err := o.computeSvc.TargetTcpProxies.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return fmt.Errorf("failed to delete target TCP proxy %s: %w", item.name, err)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return fmt.Errorf("failed to delete target TCP proxy %s with error: %s: %w", item.name, operationErrorMessage(op), err)
	}
	if (err != nil && isNoOp(err)) || (op != nil && op.Status == "DONE") {
		o.resetRequestID(item.typeName, item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted target TCP proxy %s", item.name)
	}
	return nil
}

// destroyTargetTCPProxies removes target tcp proxies with a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyTargetTCPProxies(ctx context.Context) error {
	found, err := o.listTargetTCPProxies(ctx)
	if err != nil {
		return fmt.Errorf("failed to list target TCP proxies: %w", err)
	}
	items := o.insertPendingItems("targettcpproxy", found)
	for _, item := range items {
		err := o.deleteTargetTCPProxy(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems("targettcpproxy"); len(items) > 0 {
		return fmt.Errorf("%d items pending", len(items))
	}
	return nil
}
