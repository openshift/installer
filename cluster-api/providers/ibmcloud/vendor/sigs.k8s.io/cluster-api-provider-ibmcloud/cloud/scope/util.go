/*
Copyright 2024 The Kubernetes Authors.

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

package scope

import (
	"context"
	"fmt"
	"strconv"

	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
)

// GetClusterByName finds and return a Cluster object using the specified params.
func GetClusterByName(ctx context.Context, c client.Client, namespace, name string) (*infrav1beta2.IBMPowerVSCluster, error) {
	cluster := &infrav1beta2.IBMPowerVSCluster{}
	key := client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}

	if err := c.Get(ctx, key, cluster); err != nil {
		return nil, fmt.Errorf("failed to get Cluster/%s: %w", name, err)
	}

	return cluster, nil
}

// CheckCreateInfraAnnotation checks for annotations set on IBMPowerVSCluster object to determine cluster creation workflow.
func CheckCreateInfraAnnotation(cluster infrav1beta2.IBMPowerVSCluster) bool {
	annotations := cluster.GetAnnotations()
	if len(annotations) == 0 {
		return false
	}
	value, found := annotations[infrav1beta2.CreateInfrastructureAnnotation]
	if !found {
		return false
	}
	createInfra, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}
	return createInfra
}
