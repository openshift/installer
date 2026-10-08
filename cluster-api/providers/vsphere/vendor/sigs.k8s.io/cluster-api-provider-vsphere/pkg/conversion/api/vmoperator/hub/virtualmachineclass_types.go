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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	conversionmeta "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/api/meta"
)

// VirtualMachineClassHardware describes a virtual hardware resource
// specification.
type VirtualMachineClassHardware struct {
	// +optional
	Cpus int64 `json:"cpus,omitempty"`

	// +optional
	Memory resource.Quantity `json:"memory,omitempty"`
}

// VirtualMachineClassSpec defines the desired state of VirtualMachineClass.
type VirtualMachineClassSpec struct {
	// +optional

	// Hardware describes the configuration of the VirtualMachineClass
	// attributes related to virtual hardware. The configuration specified in
	// this field is used to customize the virtual hardware characteristics of
	// any VirtualMachine associated with this VirtualMachineClass.
	Hardware VirtualMachineClassHardware `json:"hardware,omitempty"`
}

// VirtualMachineClassStatus defines the observed state of VirtualMachineClass.
type VirtualMachineClassStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=vmclass
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="CPU",type="string",JSONPath=".spec.hardware.cpus"
// +kubebuilder:printcolumn:name="Memory",type="string",JSONPath=".spec.hardware.memory"
// +kubebuilder:printcolumn:name="Capabilities",type="string",priority=1,JSONPath=".status.capabilities"

// VirtualMachineClass is the schema for the virtualmachineclasses API and
// represents the desired state and observed status of a virtualmachineclasses
// resource.
type VirtualMachineClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineClassSpec   `json:"spec,omitempty"`
	Status VirtualMachineClassStatus `json:"status,omitempty"`

	Source conversionmeta.SourceTypeMeta `json:"source,omitempty,omitzero"`
}

// GetSource returns the Source for this object.
func (c *VirtualMachineClass) GetSource() conversionmeta.SourceTypeMeta {
	return c.Source
}

// SetSource sets Source for an API object.
func (c *VirtualMachineClass) SetSource(source conversionmeta.SourceTypeMeta) {
	c.Source = source
}

// +kubebuilder:object:root=true

// VirtualMachineClassList contains a list of VirtualMachineClass.
type VirtualMachineClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineClass `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &VirtualMachineClass{}, &VirtualMachineClassList{})
}
