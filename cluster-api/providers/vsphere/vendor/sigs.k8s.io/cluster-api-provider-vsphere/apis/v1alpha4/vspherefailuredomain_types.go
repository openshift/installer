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

//nolint:godot
package v1alpha4

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type FailureDomainType string

const (
	HostGroupFailureDomain      FailureDomainType = "HostGroup"
	ComputeClusterFailureDomain FailureDomainType = "ComputeCluster"
	DatacenterFailureDomain     FailureDomainType = "Datacenter"
)

// VSphereFailureDomainSpec defines the desired state of VSphereFailureDomain
type VSphereFailureDomainSpec struct {

	// Region defines the name and type of a region
	Region FailureDomain `json:"region"`

	// Zone defines the name and type of a zone
	Zone FailureDomain `json:"zone"`

	// Topology describes a given failure domain using vSphere constructs
	Topology Topology `json:"topology"`
}

type FailureDomain struct {
	// Name is the name of the tag that represents this failure domain
	Name string `json:"name"`

	// Type is the type of failure domain, the current values are "Datacenter", "ComputeCluster" and "HostGroup"
	// +kubebuilder:validation:Enum=Datacenter;ComputeCluster;HostGroup
	Type FailureDomainType `json:"type"`

	// TagCategory is the category used for the tag
	TagCategory string `json:"tagCategory"`

	// AutoConfigure tags the Type which is specified in the Topology
	AutoConfigure *bool `json:"autoConfigure,omitempty"`
}

type Topology struct {
	// The underlying infrastructure for this failure domain
	// Datacenter as the failure domain
	// +kubebuilder:validation:Required
	Datacenter string `json:"datacenter"`

	// ComputeCluster as the failure domain
	// +optional
	ComputeCluster *string `json:"computeCluster,omitempty"`

	// Hosts has information required for placement of machines on VSphere hosts.
	// +optional
	Hosts *FailureDomainHosts `json:"hosts,omitempty"`

	// Networks is the list of networks within this failure domain
	// +optional
	Networks []string `json:"networks,omitempty"`

	// Datastore is the name or inventory path of the datastore in which the
	// virtual machine is created/located.
	// +optional
	Datastore string `json:"datastore,omitempty"`
}

type FailureDomainHosts struct {
	// VMGroupName is the name of the VM group
	VMGroupName string `json:"vmGroupName"`

	// HostGroupName is the name of the Host group
	HostGroupName string `json:"hostGroupName"`
}

// +kubebuilder:object:root=true
// +kubebuilder:unservedversion
// +kubebuilder:deprecatedversion
// +kubebuilder:resource:path=vspherefailuredomains,scope=Cluster,categories=cluster-api

// VSphereFailureDomain is the Schema for the vspherefailuredomains API
//
// Deprecated: This type will be removed in one of the next releases.
type VSphereFailureDomain struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec VSphereFailureDomainSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// VSphereFailureDomainList contains a list of VSphereFailureDomain
//
// Deprecated: This type will be removed in one of the next releases.
type VSphereFailureDomainList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereFailureDomain `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &VSphereFailureDomain{}, &VSphereFailureDomainList{})
}
