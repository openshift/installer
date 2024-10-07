package gcp

import (
	"context"
	"fmt"

	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

const (
	globalTargetTCPProxyResource = "targettcpproxy"
)

func (o *ClusterUninstaller) listTargetTCPProxies(ctx context.Context, typeName string) ([]cloudResource, error) {
	return o.listTargetTCPProxiesWithFilter(ctx, typeName, "items(name),nextPageToken", o.clusterIDFilter())
}

// listTargetTCPProxiesWithFilter lists target TCP Proxies in the project that satisfy the filter criteria.
func (o *ClusterUninstaller) listTargetTCPProxiesWithFilter(ctx context.Context, typeName, fields, filter string) ([]cloudResource, error) {
	o.Logger.Debugf("Listing target tcp proxies")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	list, err := o.computeSvc.TargetTcpProxies.List(o.ProjectID).Filter(filter).Fields(googleapi.Field(fields)).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list target tcp proxies: %w", err)
	}

	result := []cloudResource{}
	for _, item := range list.Items {
		o.Logger.Debugf("Found target TCP proxy: %s", item.Name)
		result = append(result, cloudResource{
			key:      item.Name,
			name:     item.Name,
			typeName: typeName,
			quota: []gcp.QuotaUsage{{
				Metric: &gcp.Metric{
					Service: gcp.ServiceComputeEngineAPI,
					Limit:   "target_tcp_proxy",
				},
				Amount: 1,
			}},
		})
	}

	return result, nil
}

func (o *ClusterUninstaller) deleteTargetTCPProxy(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting target tcp proxy %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	op, err := o.computeSvc.TargetTcpProxies.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return fmt.Errorf("failed to target tcp proxy %s: %w", item.name, err)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return fmt.Errorf("failed to delete target tcp proxy %s with error: %s", item.name, operationErrorMessage(op))
	}
	if (err != nil && isNoOp(err)) || (op != nil && op.Status == "DONE") {
		o.resetRequestID(item.typeName, item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted target tcp proxy %s", item.name)
	}
	return nil
}

// destroyTargetTCPProxies removes all target tcp proxy resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyTargetTCPProxies(ctx context.Context) error {
	found, err := o.listTargetTCPProxies(ctx, globalTargetTCPProxyResource)
	if err != nil {
		return err
	}
	items := o.insertPendingItems(globalTargetTCPProxyResource, found)

	for _, item := range items {
		err := o.deleteTargetTCPProxy(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}

	if items = o.getPendingItems(globalTargetTCPProxyResource); len(items) > 0 {
		return fmt.Errorf("%d global target tcp proxy pending", len(items))
	}

	return nil
}
