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

package cluster

import (
	"context"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"

	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
)

type computeClusterContext interface {
	context.Context

	GetSession() *session.Session
}

func ListHostsFromGroup(ctx context.Context, ccr *object.ClusterComputeResource, hostGroup string) ([]object.Reference, error) {
	clusterConfigInfoEx, err := ccr.Configuration(ctx)
	if err != nil {
		return nil, err
	}

	var refs []object.Reference
	for _, group := range clusterConfigInfoEx.Group {
		if clusterHostGroup, ok := group.(*types.ClusterHostGroup); ok {
			if clusterHostGroup.Name == hostGroup {
				for _, managedRef := range clusterHostGroup.Host {
					refs = append(refs, managedRef)
				}
				return refs, nil
			}
		}
	}
	return refs, nil
}
