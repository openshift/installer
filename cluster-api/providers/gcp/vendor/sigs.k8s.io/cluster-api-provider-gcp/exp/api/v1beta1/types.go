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

package v1beta1

import (
	"strings"

	"cloud.google.com/go/container/apiv1/containerpb"
)

// TaintEffect is the effect for a Kubernetes taint.
type TaintEffect string

// Taint represents a Kubernetes taint.
type Taint struct {
	// Effect specifies the effect for the taint.
	// +kubebuilder:validation:Enum=NoSchedule;NoExecute;PreferNoSchedule
	Effect TaintEffect `json:"effect"`
	// Key is the key of the taint
	Key string `json:"key"`
	// Value is the value of the taint
	Value string `json:"value"`
}

// Taints is an array of Taints.
type Taints []Taint

func convertToSdkTaintEffect(effect TaintEffect) containerpb.NodeTaint_Effect {
	switch effect {
	case "NoSchedule":
		return containerpb.NodeTaint_NO_SCHEDULE
	case "NoExecute":
		return containerpb.NodeTaint_NO_EXECUTE
	case "PreferNoSchedule":
		return containerpb.NodeTaint_PREFER_NO_SCHEDULE
	default:
		return containerpb.NodeTaint_EFFECT_UNSPECIFIED
	}
}

// ConvertToSdkTaint converts taints to format that is used by GCP SDK.
func ConvertToSdkTaint(taints Taints) []*containerpb.NodeTaint {
	if taints == nil {
		return nil
	}
	res := []*containerpb.NodeTaint{}
	for _, taint := range taints {
		res = append(res, &containerpb.NodeTaint{
			Key:    taint.Key,
			Value:  taint.Value,
			Effect: convertToSdkTaintEffect(taint.Effect),
		})
	}
	return res
}

// convertToSdkLocationPolicy converts node pool autoscaling location policy to a value that is used by GCP SDK.
func convertToSdkLocationPolicy(locationPolicy ManagedNodePoolLocationPolicy) containerpb.NodePoolAutoscaling_LocationPolicy {
	switch locationPolicy {
	case ManagedNodePoolLocationPolicyBalanced:
		return containerpb.NodePoolAutoscaling_BALANCED
	case ManagedNodePoolLocationPolicyAny:
		return containerpb.NodePoolAutoscaling_ANY
	}
	return containerpb.NodePoolAutoscaling_LOCATION_POLICY_UNSPECIFIED
}

// ConvertToSdkAutoscaling converts node pool autoscaling config to a value that is used by GCP SDK.
func ConvertToSdkAutoscaling(autoscaling *NodePoolAutoScaling) *containerpb.NodePoolAutoscaling {
	sdkAutoscaling := containerpb.NodePoolAutoscaling{
		Enabled:           true, // enable autoscaling by default
		TotalMinNodeCount: 0,
		TotalMaxNodeCount: 1,
		LocationPolicy:    convertToSdkLocationPolicy(ManagedNodePoolLocationPolicyBalanced),
	}
	if autoscaling != nil {
		// set fields
		if autoscaling.MinCount != nil {
			sdkAutoscaling.TotalMinNodeCount = *autoscaling.MinCount
		}
		if autoscaling.MaxCount != nil {
			sdkAutoscaling.TotalMaxNodeCount = *autoscaling.MaxCount
		}
		if autoscaling.LocationPolicy != nil {
			sdkAutoscaling.LocationPolicy = convertToSdkLocationPolicy(*autoscaling.LocationPolicy)
		}
		if autoscaling.EnableAutoscaling != nil {
			if !*autoscaling.EnableAutoscaling {
				sdkAutoscaling = containerpb.NodePoolAutoscaling{
					Enabled: false,
				}
			}
		}
	}

	return &sdkAutoscaling
}

// ConvertFromSdkNodeVersion converts GCP SDK node version to k8s version.
func ConvertFromSdkNodeVersion(sdkNodeVersion string) string {
	// For example, the node version returned from GCP SDK can be 1.27.2-gke.2100, we want to convert it to 1.27.2
	return strings.Replace(strings.Split(sdkNodeVersion, "-")[0], "v", "", 1)
}

// ConvertToSdkCgroupMode converts GCP SDK node version to k8s version.
func ConvertToSdkCgroupMode(cgroupMode ManagedNodePoolCgroupMode) containerpb.LinuxNodeConfig_CgroupMode {
	switch cgroupMode {
	case 1:
		return containerpb.LinuxNodeConfig_CGROUP_MODE_V1
	case 2:
		return containerpb.LinuxNodeConfig_CGROUP_MODE_V2
	}
	return containerpb.LinuxNodeConfig_CGROUP_MODE_UNSPECIFIED
}

// ConvertToSdkLinuxNodeConfig converts GCP SDK node version to k8s version.
func ConvertToSdkLinuxNodeConfig(linuxNodeConfig *LinuxNodeConfig) *containerpb.LinuxNodeConfig {
	sdkLinuxNodeConfig := containerpb.LinuxNodeConfig{}
	if linuxNodeConfig != nil {
		if linuxNodeConfig.Sysctls != nil {
			sdkSysctl := make(map[string]string)
			for _, sysctl := range linuxNodeConfig.Sysctls {
				sdkSysctl[sysctl.Parameter] = sysctl.Value
			}
			sdkLinuxNodeConfig.Sysctls = sdkSysctl
		}
		if linuxNodeConfig.CgroupMode != nil {
			sdkLinuxNodeConfig.CgroupMode = ConvertToSdkCgroupMode(*linuxNodeConfig.CgroupMode)
		}
	}
	return &sdkLinuxNodeConfig
}
