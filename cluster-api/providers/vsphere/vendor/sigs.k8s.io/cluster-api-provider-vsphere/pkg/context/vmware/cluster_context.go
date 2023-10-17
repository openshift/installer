/*
Copyright 2021 The Kubernetes Authors.

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

package vmware

import (
	"fmt"

	"github.com/go-logr/logr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"

	vmwarev1b1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
)

// ClusterContext is a Go context used with a CAPI cluster.
type ClusterContext struct {
	*context.ControllerContext
	Cluster        *clusterv1.Cluster
	VSphereCluster *vmwarev1b1.VSphereCluster
	PatchHelper    *patch.Helper
	Logger         logr.Logger
}

// String returns ControllerManagerName/ControllerName/ClusterAPIVersion/ClusterNamespace/ClusterName.
func (c *ClusterContext) String() string {
	return fmt.Sprintf("%s/%s/%s/%s", c.ControllerContext.String(), c.VSphereCluster.APIVersion, c.VSphereCluster.Namespace, c.VSphereCluster.Name)
}

// Patch updates the object and its status on the API server.
func (c *ClusterContext) Patch() error {
	// always update the readyCondition.
	conditions.SetSummary(c.VSphereCluster,
		conditions.WithConditions(
			vmwarev1b1.ResourcePolicyReadyCondition,
			vmwarev1b1.ClusterNetworkReadyCondition,
			vmwarev1b1.LoadBalancerReadyCondition,
		),
	)
	return c.PatchHelper.Patch(c, c.VSphereCluster)
}
