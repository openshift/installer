/*
Copyright (c) 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

import (
	"net/http"
	"path"
)

// MetricQueriesClient is the client of the 'metric_queries' resource.
//
// Manages telemetry queries for a cluster.
type MetricQueriesClient struct {
	transport http.RoundTripper
	path      string
}

// NewMetricQueriesClient creates a new client for the 'metric_queries'
// resource using the given transport to send the requests and receive the
// responses.
func NewMetricQueriesClient(transport http.RoundTripper, path string) *MetricQueriesClient {
	return &MetricQueriesClient{
		transport: transport,
		path:      path,
	}
}

// CPUTotalByNodeRolesOS returns the target 'CPU_total_by_node_roles_OS_metric_query' resource.
//
// Reference to the resource that retrieves the total cpu
// capacity in the cluster by node role and operating system.
func (c *MetricQueriesClient) CPUTotalByNodeRolesOS() *CPUTotalByNodeRolesOSMetricQueryClient {
	return NewCPUTotalByNodeRolesOSMetricQueryClient(
		c.transport,
		path.Join(c.path, "cpu_total_by_node_roles_os"),
	)
}

// Alerts returns the target 'alerts_metric_query' resource.
//
// Reference to the resource that retrieves the firing alerts in the cluster.
func (c *MetricQueriesClient) Alerts() *AlertsMetricQueryClient {
	return NewAlertsMetricQueryClient(
		c.transport,
		path.Join(c.path, "alerts"),
	)
}

// ClusterOperators returns the target 'cluster_operators_metric_query' resource.
//
// Reference to the resource that retrieves the cluster operator status metrics.
func (c *MetricQueriesClient) ClusterOperators() *ClusterOperatorsMetricQueryClient {
	return NewClusterOperatorsMetricQueryClient(
		c.transport,
		path.Join(c.path, "cluster_operators"),
	)
}

// Nodes returns the target 'nodes_metric_query' resource.
//
// Reference to the resource that retrieves the nodes in the cluster.
func (c *MetricQueriesClient) Nodes() *NodesMetricQueryClient {
	return NewNodesMetricQueryClient(
		c.transport,
		path.Join(c.path, "nodes"),
	)
}

// SocketTotalByNodeRolesOS returns the target 'socket_total_by_node_roles_OS_metric_query' resource.
//
// Reference to the resource that retrieves the total socket
// capacity in the cluster by node role and operating system.
func (c *MetricQueriesClient) SocketTotalByNodeRolesOS() *SocketTotalByNodeRolesOSMetricQueryClient {
	return NewSocketTotalByNodeRolesOSMetricQueryClient(
		c.transport,
		path.Join(c.path, "socket_total_by_node_roles_os"),
	)
}
