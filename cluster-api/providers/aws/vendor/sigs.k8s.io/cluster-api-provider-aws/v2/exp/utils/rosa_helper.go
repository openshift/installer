/*
Copyright 2025 The Kubernetes Authors.

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

// Package utils provide helper functions.
package utils

import (
	"time"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"

	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
)

// NodePoolToRosaMachinePoolSpec convert ocm nodePool to rosaMachinePool spec.
func NodePoolToRosaMachinePoolSpec(nodePool *cmv1.NodePool) expinfrav1.RosaMachinePoolSpec {
	spec := expinfrav1.RosaMachinePoolSpec{
		NodePoolName:             nodePool.ID(),
		Version:                  rosa.RawVersionID(nodePool.Version()),
		AvailabilityZone:         nodePool.AvailabilityZone(),
		Subnet:                   nodePool.Subnet(),
		Labels:                   nodePool.Labels(),
		AutoRepair:               nodePool.AutoRepair(),
		InstanceType:             nodePool.AWSNodePool().InstanceType(),
		TuningConfigs:            nodePool.TuningConfigs(),
		AdditionalSecurityGroups: nodePool.AWSNodePool().AdditionalSecurityGroupIds(),
		VolumeSize:               nodePool.AWSNodePool().RootVolume().Size(),
		CapacityReservationID:    nodePool.AWSNodePool().CapacityReservation().Id(),
		// nodePool.AWSNodePool().Tags() returns all tags including "system" tags if "fetchUserTagsOnly" parameter is not specified.
		// TODO: enable when AdditionalTags day2 changes is supported.
		// AdditionalTags:           nodePool.AWSNodePool().Tags(),
	}

	if nodePool.Autoscaling() != nil {
		spec.Autoscaling = &rosacontrolplanev1.AutoScaling{
			MinReplicas: nodePool.Autoscaling().MinReplica(),
			MaxReplicas: nodePool.Autoscaling().MaxReplica(),
		}
	}
	if nodePool.Taints() != nil {
		rosaTaints := make([]expinfrav1.RosaTaint, 0, len(nodePool.Taints()))
		for _, taint := range nodePool.Taints() {
			rosaTaints = append(rosaTaints, expinfrav1.RosaTaint{
				Key:    taint.Key(),
				Value:  taint.Value(),
				Effect: corev1.TaintEffect(taint.Effect()),
			})
		}
		spec.Taints = rosaTaints
	}
	if nodePool.NodeDrainGracePeriod() != nil {
		spec.NodeDrainGracePeriod = &metav1.Duration{
			Duration: time.Minute * time.Duration(nodePool.NodeDrainGracePeriod().Value()),
		}
	}
	if nodePool.ManagementUpgrade() != nil {
		spec.UpdateConfig = &expinfrav1.RosaUpdateConfig{
			RollingUpdate: &expinfrav1.RollingUpdate{},
		}
		if nodePool.ManagementUpgrade().MaxSurge() != "" {
			spec.UpdateConfig.RollingUpdate.MaxSurge = ptr.To(intstr.Parse(nodePool.ManagementUpgrade().MaxSurge()))
		}
		if nodePool.ManagementUpgrade().MaxUnavailable() != "" {
			spec.UpdateConfig.RollingUpdate.MaxUnavailable = ptr.To(intstr.Parse(nodePool.ManagementUpgrade().MaxUnavailable()))
		}
	}

	return spec
}
