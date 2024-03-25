/*
Copyright 2023 The Kubernetes Authors.

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

package shared

import (
	"errors"
	"fmt"
	"strings"

	clusterv1exp "sigs.k8s.io/cluster-api/exp/api/v1beta1"

	"sigs.k8s.io/cluster-api-provider-gcp/cloud"
	infrav1exp "sigs.k8s.io/cluster-api-provider-gcp/exp/api/v1beta1"
)

// ManagedMachinePoolPreflightCheck will perform checks against the machine pool before its created.
func ManagedMachinePoolPreflightCheck(managedPool *infrav1exp.GCPManagedMachinePool, machinePool *clusterv1exp.MachinePool, location string) error {
	if machinePool.Spec.Template.Spec.InfrastructureRef.Name != managedPool.Name {
		return fmt.Errorf("expect machinepool infraref (%s) to match managed machine pool name (%s)", machinePool.Spec.Template.Spec.InfrastructureRef.Name, managedPool.Name)
	}

	if IsRegional(location) {
		var numRegionsPerZone int32
		if len(managedPool.Spec.NodeLocations) != 0 {
			numRegionsPerZone = int32(len(managedPool.Spec.NodeLocations))
		} else {
			numRegionsPerZone = cloud.DefaultNumRegionsPerZone
		}
		if *machinePool.Spec.Replicas%numRegionsPerZone != 0 {
			return fmt.Errorf("a machine pool (%s) in a regional cluster with %d regions, must have replicas with a multiple of %d", machinePool.Name, numRegionsPerZone, numRegionsPerZone)
		}
	}

	return nil
}

// ManagedMachinePoolsPreflightCheck will perform checks against a slice of machine pool before they are created.
func ManagedMachinePoolsPreflightCheck(managedPools []infrav1exp.GCPManagedMachinePool, machinePools []clusterv1exp.MachinePool, location string) error {
	if len(machinePools) != len(managedPools) {
		return errors.New("each machinepool must have a matching gcpmanagedmachinepool")
	}

	for i := range machinePools {
		machinepool := machinePools[i]
		managedPool := managedPools[i]

		if err := ManagedMachinePoolPreflightCheck(&managedPool, &machinepool, location); err != nil {
			return err
		}
	}

	return nil
}

// IsRegional will check if a given location is a region (if not its a zone).
func IsRegional(location string) bool {
	return strings.Count(location, "-") == 1
}
