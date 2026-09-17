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

package v1beta2

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// VPCLoadBalancerSpec defines the desired state of an VPC load balancer.
type VPCLoadBalancerSpec struct {
	// Name sets the name of the VPC load balancer.
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=63
	// +kubebuilder:validation:Pattern=`^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`
	// +optional
	Name string `json:"name,omitempty"`

	// id of the loadbalancer
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength:=64
	// +kubebuilder:validation:Pattern=`^[-0-9a-z_]+$`
	// +optional
	ID *string `json:"id,omitempty"`

	// public indicates that load balancer is public or private
	// +kubebuilder:default=true
	// +optional
	Public *bool `json:"public,omitempty"`

	// AdditionalListeners sets the additional listeners for the control plane load balancer.
	// +listType=map
	// +listMapKey=port
	// +optional
	AdditionalListeners []AdditionalListenerSpec `json:"additionalListeners,omitempty"`

	// backendPools defines the load balancer's backend pools.
	// +optional
	BackendPools []VPCLoadBalancerBackendPoolSpec `json:"backendPools,omitempty"`

	// securityGroups defines the Security Groups to attach to the load balancer.
	// Security Groups defined here are expected to already exist when the load balancer is reconciled (these do not get created when reconciling the load balancer).
	// +optional
	SecurityGroups []VPCResource `json:"securityGroups,omitempty"`

	// subnets defines the VPC Subnets to attach to the load balancer.
	// Subnets defiens here are expected to already exist when the load balancer is reconciled (these do not get created when reconciling the load balancer).
	// +optional
	Subnets []VPCResource `json:"subnets,omitempty"`
}

// AdditionalListenerSpec defines the desired state of an
// additional listener on an VPC load balancer.
type AdditionalListenerSpec struct {
	// defaultPoolName defines the name of a VPC Load Balancer Backend Pool to use for the VPC Load Balancer Listener.
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=63
	// +kubebuilder:validation:Pattern=`^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`
	// +optional
	DefaultPoolName *string `json:"defaultPoolName,omitempty"`

	// Port sets the port for the additional listener.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int64 `json:"port"`

	// protocol defines the protocol to use for the VPC Load Balancer Listener.
	// Will default to TCP protocol if not specified.
	// +optional
	Protocol *VPCLoadBalancerListenerProtocol `json:"protocol,omitempty"`

	// The selector is used to find IBMPowerVSMachines with matching labels.
	// If the label matches, the machine is then added to the load balancer listener configuration.
	// +kubebuilder:validation:Optional
	Selector metav1.LabelSelector `json:"selector,omitempty"`
}

// VPCLoadBalancerBackendPoolSpec defines the desired configuration of a VPC Load Balancer Backend Pool.
type VPCLoadBalancerBackendPoolSpec struct {
	// name defines the name of the Backend Pool.
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=63
	// +kubebuilder:validation:Pattern=`^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`
	// +optional
	Name *string `json:"name,omitempty"`

	// algorithm defines the load balancing algorithm to use.
	// +required
	Algorithm VPCLoadBalancerBackendPoolAlgorithm `json:"algorithm"`

	// healthMonitor defines the backend pool's health monitor.
	// +required
	HealthMonitor VPCLoadBalancerHealthMonitorSpec `json:"healthMonitor"`

	// protocol defines the protocol to use for the Backend Pool.
	// +required
	Protocol VPCLoadBalancerBackendPoolProtocol `json:"protocol"`
}

// VPCLoadBalancerHealthMonitorSpec defines the desired state of a Health Monitor resource for a VPC Load Balancer Backend Pool.
// kubebuilder:validation:XValidation:rule="self.dely > self.timeout",message="health monitor's delay must be greater than the timeout"
type VPCLoadBalancerHealthMonitorSpec struct {
	// delay defines the seconds to wait between health checks.
	// +kubebuilder:validation:Minimum=2
	// +kubebuilder:validation:Maximum=60
	// +required
	Delay int64 `json:"delay"`

	// retries defines the max retries for health check.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	// +required
	Retries int64 `json:"retries"`

	// port defines the port to perform health monitoring on.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +optional
	Port *int64 `json:"port,omitempty"`

	// timeout defines the seconds to wait for a health check response.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=59
	// +required
	Timeout int64 `json:"timeout"`

	// type defines the protocol used for health checks.
	// +required
	Type VPCLoadBalancerBackendPoolHealthMonitorType `json:"type"`

	// urlPath defines the URL to use for health monitoring.
	// +kubebuilder:validation:Pattern=`^\/(([a-zA-Z0-9-._~!$&'()*+,;=:@]|%[a-fA-F0-9]{2})+(\/([a-zA-Z0-9-._~!$&'()*+,;=:@]|%[a-fA-F0-9]{2})*)*)?(\\?([a-zA-Z0-9-._~!$&'()*+,;=:@\/?]|%[a-fA-F0-9]{2})*)?$`
	// +optional
	URLPath *string `json:"urlPath,omitempty"`
}

// VPCSecurityGroupStatus defines a vpc security group resource status with its id and respective rule's ids.
type VPCSecurityGroupStatus struct {
	// id represents the id of the resource.
	ID *string `json:"id,omitempty"`
	// rules contains the id of rules created under the security group
	RuleIDs []*string `json:"ruleIDs,omitempty"`
	// +kubebuilder:default=false
	// controllerCreated indicates whether the resource is created by the controller.
	ControllerCreated *bool `json:"controllerCreated,omitempty"`
}

// VPCLoadBalancerStatus defines the status VPC load balancer.
type VPCLoadBalancerStatus struct {
	// id of VPC load balancer.
	// +optional
	ID *string `json:"id,omitempty"`
	// State is the status of the load balancer.
	State VPCLoadBalancerState `json:"state,omitempty"`
	// hostname is the hostname of load balancer.
	// +optional
	Hostname *string `json:"hostname,omitempty"`
	// +kubebuilder:default=false
	// controllerCreated indicates whether the resource is created by the controller.
	ControllerCreated *bool `json:"controllerCreated,omitempty"`
}
