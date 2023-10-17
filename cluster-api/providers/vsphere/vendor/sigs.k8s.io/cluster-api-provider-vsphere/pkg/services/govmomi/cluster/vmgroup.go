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

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

func FindVMGroup(ctx computeClusterContext, clusterName, vmGroupName string) (*VMGroup, error) {
	ccr, err := ctx.GetSession().Finder.ClusterComputeResource(ctx, clusterName)
	if err != nil {
		return nil, err
	}

	clusterConfigInfoEx, err := ccr.Configuration(ctx)
	if err != nil {
		return nil, err
	}
	for _, group := range clusterConfigInfoEx.Group {
		if clusterVMGroup, ok := group.(*types.ClusterVmGroup); ok {
			if clusterVMGroup.Name == vmGroupName {
				return &VMGroup{ccr, clusterVMGroup}, nil
			}
		}
	}
	return nil, errors.Errorf("cannot find VM group %s", vmGroupName)
}

// VMGroup represents a VSphere VM Group object.
type VMGroup struct {
	*object.ClusterComputeResource
	*types.ClusterVmGroup
}

// Add a VSphere VM object to the VM Group.
func (vg VMGroup) Add(ctx context.Context, vmObj types.ManagedObjectReference) (*object.Task, error) {
	vms := vg.listVMs()
	vg.ClusterVmGroup.Vm = append(vms, vmObj) //nolint:gocritic

	spec := &types.ClusterConfigSpecEx{
		GroupSpec: []types.ClusterGroupSpec{
			{
				ArrayUpdateSpec: types.ArrayUpdateSpec{
					Operation: types.ArrayUpdateOperationEdit,
				},
				Info: vg.ClusterVmGroup,
			},
		},
	}
	return vg.ClusterComputeResource.Reconfigure(ctx, spec, true)
}

// HasVM returns whether a VSphere VM object is a member of the VM Group.
func (vg VMGroup) HasVM(vmObj types.ManagedObjectReference) (bool, error) {
	vms := vg.listVMs()

	for _, vm := range vms {
		if vm == vmObj {
			return true, nil
		}
	}
	return false, nil
}

func (vg VMGroup) listVMs() []types.ManagedObjectReference {
	return vg.ClusterVmGroup.Vm
}
