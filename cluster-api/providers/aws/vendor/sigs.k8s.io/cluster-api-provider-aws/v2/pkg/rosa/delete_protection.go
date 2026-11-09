/*
Copyright 2026 The Kubernetes Authors.

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

package rosa

import (
	"fmt"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

// DeleteProtectionEnabledFromCluster reads delete protection from the inline cluster field.
func DeleteProtectionEnabledFromCluster(cluster *cmv1.Cluster) (enabled bool, ok bool) {
	if cluster == nil {
		return false, false
	}
	if dp, present := cluster.GetDeleteProtection(); present && dp != nil {
		return dp.Enabled(), true
	}
	return false, false
}

// UpdateClusterDeletionProtection patches the cluster delete protection sub-resource.
func UpdateClusterDeletionProtection(ocmClient OCMClient, clusterID string, enabled bool) error {
	return ocmClient.UpdateClusterDeletionProtection(clusterID, enabled)
}

// ReconcileDeleteProtection syncs spec.deleteProtection to OCM.
func ReconcileDeleteProtection(
	rosaScope *scope.ROSAControlPlaneScope,
	ocmClient OCMClient,
	cluster *cmv1.Cluster,
) error {
	liveEnabled, ok := DeleteProtectionEnabledFromCluster(cluster)
	if !ok {
		return fmt.Errorf("failed to read delete protection for cluster '%s'", cluster.ID())
	}

	specEnabled := rosaScope.ControlPlane.Spec.DeleteProtection == rosacontrolplanev1.DeleteProtectionEnabled
	if specEnabled == liveEnabled {
		return nil
	}

	if err := UpdateClusterDeletionProtection(ocmClient, cluster.ID(), specEnabled); err != nil {
		return fmt.Errorf("failed to update delete protection for cluster '%s': %w", cluster.ID(), err)
	}

	return nil
}

// IsDeleteProtectionBlocking reports whether cluster deletion should be blocked.
func IsDeleteProtectionBlocking(
	rosaScope *scope.ROSAControlPlaneScope,
	cluster *cmv1.Cluster,
) bool {
	if rosaScope.ControlPlane.Spec.DeleteProtection == rosacontrolplanev1.DeleteProtectionEnabled {
		return true
	}

	enabled, ok := DeleteProtectionEnabledFromCluster(cluster)
	if !ok {
		return false
	}

	return enabled
}
