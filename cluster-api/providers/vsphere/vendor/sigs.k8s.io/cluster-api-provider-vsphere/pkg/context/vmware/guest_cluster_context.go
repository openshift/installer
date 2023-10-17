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

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GuestClusterContext is the context used for GuestClusterControllers.
type GuestClusterContext struct {
	*ClusterContext

	// GuestClient can be used to access the guest cluster.
	GuestClient client.Client
}

// String returns ClusterGroupVersionKind ClusterNamespace/ClusterName.
func (c *GuestClusterContext) String() string {
	return fmt.Sprintf("%s %s/%s", c.Cluster.GroupVersionKind(), c.Cluster.Namespace, c.Cluster.Name)
}
