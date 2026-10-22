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

package hub

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	conversionmeta "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/api/meta"
)

// ResourcePoolSpec defines a Logical Grouping of workloads that share resource
// policies.
type ResourcePoolSpec struct {
	// +optional

	// Name describes the name of the ResourcePool grouping.
	Name string `json:"name,omitempty"`
}

// VirtualMachineSetResourcePolicySpec defines the desired state of
// VirtualMachineSetResourcePolicy.
type VirtualMachineSetResourcePolicySpec struct {
	ResourcePool        ResourcePoolSpec `json:"resourcePool,omitempty"`
	Folder              string           `json:"folder,omitempty"`
	ClusterModuleGroups []string         `json:"clusterModuleGroups,omitempty"`
}

// VirtualMachineSetResourcePolicyStatus defines the observed state of
// VirtualMachineSetResourcePolicy.
type VirtualMachineSetResourcePolicyStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// VirtualMachineSetResourcePolicy is the Schema for the virtualmachinesetresourcepolicies API.
type VirtualMachineSetResourcePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineSetResourcePolicySpec   `json:"spec,omitempty"`
	Status VirtualMachineSetResourcePolicyStatus `json:"status,omitempty"`

	Source conversionmeta.SourceTypeMeta `json:"source,omitempty,omitzero"`
}

// GetSource returns the Source for this object.
func (p *VirtualMachineSetResourcePolicy) GetSource() conversionmeta.SourceTypeMeta {
	return p.Source
}

// SetSource sets Source for an API object.
func (p *VirtualMachineSetResourcePolicy) SetSource(source conversionmeta.SourceTypeMeta) {
	p.Source = source
}

// +kubebuilder:object:root=true

// VirtualMachineSetResourcePolicyList contains a list of VirtualMachineSetResourcePolicy.
type VirtualMachineSetResourcePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineSetResourcePolicy `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &VirtualMachineSetResourcePolicy{}, &VirtualMachineSetResourcePolicyList{})
}
