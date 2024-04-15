/*
Copyright 2022 The Kubernetes Authors.

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

package clustermodule

import (
	"sort"

	"github.com/blang/semver"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
)

// Compare returns whether both the cluster module slices are the same.
func Compare(oldMods, newMods []infrav1.ClusterModule) bool {
	if len(oldMods) != len(newMods) {
		return false
	}

	sort.SliceStable(oldMods, func(i, j int) bool {
		return oldMods[i].TargetObjectName < oldMods[j].TargetObjectName
	})
	sort.SliceStable(newMods, func(i, j int) bool {
		return newMods[i].TargetObjectName < newMods[j].TargetObjectName
	})

	for i := range oldMods {
		if oldMods[i].ControlPlane == newMods[i].ControlPlane &&
			oldMods[i].TargetObjectName == newMods[i].TargetObjectName &&
			oldMods[i].ModuleUUID == newMods[i].ModuleUUID {
			continue
		}
		return false
	}
	return true
}

// IsClusterCompatible checks if the VCenterVersion is compatibly with CAPV. Only version 7 and over are supported.
func IsClusterCompatible(clusterCtx *capvcontext.ClusterContext) bool {
	version := clusterCtx.VSphereCluster.Status.VCenterVersion
	if version == "" {
		return false
	}
	apiVersion, err := semver.New(string(version))
	if err != nil {
		return false
	}

	return apiVersion.Major >= 7
}
