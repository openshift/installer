package gcp

import (
	"context"
	"fmt"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

const (
	globalTargetTCPProxyResource = "targettcpproxy"
)

func (o *ClusterUninstaller) deleteTargetTCPProxyByName(ctx context.Context, resourceName string) error {
	items, err := o.listTargetTCPProxiesWithFilter(ctx, globalTargetTCPProxyResource, "items(name),nextPageToken", func(item string) bool {
		return item == resourceName
	})
	if err != nil {
		return fmt.Errorf("failed to list target TCP Proxy by name: %w", err)
	}
	for _, item := range items {
		if err := o.deleteTargetTCPProxy(ctx, item); err != nil {
			return fmt.Errorf("failed to delete target TCP Proxy by name: %w", err)
		}
	}
	return nil
}

func (o *ClusterUninstaller) listTargetTCPProxies(ctx context.Context, typeName string) ([]cloudResource, error) {
	return o.listTargetTCPProxiesWithFilter(ctx, typeName, "items(name),nextPageToken", o.isClusterResource)
}

// listTargetTCPProxiesWithFilter lists target TCP Proxies in the project that satisfy the filter criteria.
func (o *ClusterUninstaller) listTargetTCPProxiesWithFilter(ctx context.Context, typeName, fields string, filterFunc resourceFilterFunc) ([]cloudResource, error) {
	o.Logger.Debugf("Listing target tcp proxies")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	result := []cloudResource{}

	pagesFunc := func(list *compute.TargetTcpProxyList) error {
		for _, item := range list.Items {
			if filterFunc(item.Name) {
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
		}
		return nil
	}

	err := o.computeSvc.TargetTcpProxies.List(o.ProjectID).Fields(googleapi.Field(fields)).Pages(ctx, pagesFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to list target tcp proxies: %w", err)
	}

	return result, nil
}

func (o *ClusterUninstaller) deleteTargetTCPProxy(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting target tcp proxy %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	op, err := o.computeSvc.TargetTcpProxies.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	item.scope = global
	return o.handleOperation(ctx, op, err, item, "target tcp proxy")
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
