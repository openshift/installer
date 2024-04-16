// Copyright (c) 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VirtualMachineServiceType string describes ingress methods for a service.
type VirtualMachineServiceType string

// These types correspond to a subset of the core Service Types.
const (
	// VirtualMachineServiceTypeClusterIP means a service will only be
	// accessible inside the cluster, via the cluster IP.
	VirtualMachineServiceTypeClusterIP VirtualMachineServiceType = "ClusterIP"

	// VirtualMachineServiceTypeLoadBalancer means a service will be exposed via
	// an external load balancer (if the cloud provider supports it), in
	// addition to 'NodePort' type.
	VirtualMachineServiceTypeLoadBalancer VirtualMachineServiceType = "LoadBalancer"

	// VirtualMachineServiceTypeExternalName means a service consists of only a
	// reference to an external name that kubedns or equivalent will return as a
	// CNAME record, with no exposing or proxying of any VirtualMachines
	// involved.
	VirtualMachineServiceTypeExternalName VirtualMachineServiceType = "ExternalName"
)

// VirtualMachineServicePort describes the specification of a service port to
// be exposed by a VirtualMachineService. This VirtualMachineServicePort
// specification includes attributes that define the external and internal
// representation of the service port.
type VirtualMachineServicePort struct {
	// Name describes the name to be used to identify this
	// VirtualMachineServicePort.
	Name string `json:"name"`

	// Protocol describes the Layer 4 transport protocol for this port.
	// Supports "TCP", "UDP", and "SCTP".
	Protocol string `json:"protocol"`

	// Port describes the external port that will be exposed by the service.
	Port int32 `json:"port"`

	// TargetPort describes the internal port open on a VirtualMachine that
	// should be mapped to the external Port.
	TargetPort int32 `json:"targetPort"`
}

// LoadBalancerStatus represents the status of a load balancer.
type LoadBalancerStatus struct {
	// Ingress is a list containing ingress addresses for the load balancer.
	// Traffic intended for the service should be sent to any of these ingress
	// points.
	// +optional
	Ingress []LoadBalancerIngress `json:"ingress,omitempty"`
}

// LoadBalancerIngress represents the status of a load balancer ingress point:
// traffic intended for the service should be sent to an ingress point.
// IP or Hostname may both be set in this structure. It is up to the consumer to
// determine which field should be used when accessing this LoadBalancer.
type LoadBalancerIngress struct {
	// IP is set for load balancer ingress points that are specified by an IP
	// address.
	// +optional
	IP string `json:"ip,omitempty"`

	// Hostname is set for load balancer ingress points that are specified by a
	// DNS address.
	// +optional
	Hostname string `json:"hostname,omitempty"`
}

// VirtualMachineServiceSpec defines the desired state of VirtualMachineService.
type VirtualMachineServiceSpec struct {
	// Type specifies a desired VirtualMachineServiceType for this
	// VirtualMachineService. Supported types are ClusterIP, LoadBalancer,
	// ExternalName.
	Type VirtualMachineServiceType `json:"type"`

	// Ports specifies a list of VirtualMachineServicePort to expose with this
	// VirtualMachineService. Each of these ports will be an accessible network
	// entry point to access this service by.
	Ports []VirtualMachineServicePort `json:"ports,omitempty"`

	// Selector specifies a map of key-value pairs, also known as a Label
	// Selector, that is used to match this VirtualMachineService with the set
	// of VirtualMachines that should back this VirtualMachineService.
	// +optional
	Selector map[string]string `json:"selector,omitempty"`

	// Only applies to VirtualMachineService Type: LoadBalancer
	// LoadBalancer will get created with the IP specified in this field.
	// This feature depends on whether the underlying load balancer provider
	// supports specifying the loadBalancerIP when a load balancer is created.
	// This field will be ignored if the provider does not support the feature.
	// +optional
	LoadBalancerIP string `json:"loadBalancerIP,omitempty"`

	// LoadBalancerSourceRanges is an array of IP addresses in the format of
	// CIDRs, for example: 103.21.244.0/22 and 10.0.0.0/24.
	// If specified and supported by the load balancer provider, this will
	// restrict ingress traffic to the specified client IPs. This field will be
	// ignored if the provider does not support the feature.
	// +optional
	LoadBalancerSourceRanges []string `json:"loadBalancerSourceRanges,omitempty"`

	// clusterIP is the IP address of the service and is usually assigned
	// randomly by the master. If an address is specified manually and is not in
	// use by others, it will be allocated to the service; otherwise, creation
	// of the service will fail. This field can not be changed through updates.
	// Valid values are "None", empty string (""), or a valid IP address. "None"
	// can be specified for headless services when proxying is not required.
	// Only applies to types ClusterIP and LoadBalancer.
	// Ignored if type is ExternalName.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
	// +optional
	ClusterIP string `json:"clusterIp,omitempty"`

	// externalName is the external reference that kubedns or equivalent will
	// return as a CNAME record for this service. No proxying will be involved.
	// Must be a valid RFC-1123 hostname (https://tools.ietf.org/html/rfc1123)
	// and requires Type to be ExternalName.
	// +optional
	ExternalName string `json:"externalName,omitempty"`
}

// VirtualMachineServiceStatus defines the observed state of
// VirtualMachineService.
type VirtualMachineServiceStatus struct {
	// LoadBalancer contains the current status of the load balancer,
	// if one is present.
	// +optional
	LoadBalancer LoadBalancerStatus `json:"loadBalancer,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=vmservice
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// VirtualMachineService is the Schema for the virtualmachineservices API.
type VirtualMachineService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineServiceSpec   `json:"spec,omitempty"`
	Status VirtualMachineServiceStatus `json:"status,omitempty"`
}

func (s *VirtualMachineService) NamespacedName() string {
	return s.Namespace + "/" + s.Name
}

// +kubebuilder:object:root=true

// VirtualMachineServiceList contains a list of VirtualMachineService.
type VirtualMachineServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineService `json:"items"`
}

func init() {
	SchemeBuilder.Register(
		&VirtualMachineService{},
		&VirtualMachineServiceList{},
	)
}
