/*
Copyright 2019 The Kubernetes Authors.

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

package context

import (
	"fmt"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
)

// ClusterContext is a Go context used with a VSphereCluster.
type ClusterContext struct {
	Cluster        *clusterv1.Cluster
	VSphereCluster *infrav1.VSphereCluster
	PatchHelper    *patch.Helper
}

// String returns VSphereClusterGroupVersionKind VSphereClusterNamespace/VSphereClusterName.
func (c *ClusterContext) String() string {
	return fmt.Sprintf("%s %s/%s", c.VSphereCluster.GroupVersionKind(), c.VSphereCluster.Namespace, c.VSphereCluster.Name)
}
