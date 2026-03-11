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

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	// ClusterFinalizer allows DockerClusterReconciler to clean up resources associated with DockerCluster before
	// removing it from the apiserver.
	ClusterFinalizer = "ibmvpccluster.infrastructure.cluster.x-k8s.io"
)

// IBMVPCClusterSpec defines the desired state of IBMVPCCluster.
type IBMVPCClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The IBM Cloud Region the cluster lives in.
	Region string `json:"region"`

	// The VPC resources should be created under the resource group.
	ResourceGroup string `json:"resourceGroup"`

	// The Name of VPC.
	VPC string `json:"vpc,omitempty"`

	// The Name of availability zone.
	Zone string `json:"zone,omitempty"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint capiv1beta1.APIEndpoint `json:"controlPlaneEndpoint"`

	// ControlPlaneLoadBalancer is optional configuration for customizing control plane behavior.
	// Use this for legacy support, use Network.LoadBalancers for the extended VPC support.
	// +optional
	ControlPlaneLoadBalancer *VPCLoadBalancerSpec `json:"controlPlaneLoadBalancer,omitempty"`

	// image represents the Image details used for the cluster.
	// +optional
	Image *ImageSpec `json:"image,omitempty"`

	// network represents the VPC network to use for the cluster.
	// +optional
	Network *VPCNetworkSpec `json:"network,omitempty"`
}

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
	// ++kubebuilder:validation:UniqueItems=true
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

// ImageSpec defines the desired state of the VPC Custom Image resources for the cluster.
// +kubebuilder:validation:XValidation:rule="(!has(self.cosInstance) && !has(self.cosBucket) && !has(self.cosObject)) || (has(self.cosInstance) && has(self.cosBucket) && has(self.cosObject))",message="if any of cosInstance, cosBucket, or cosObject are specified, all must be specified"
// +kubebuilder:validation:XValidation:rule="has(self.name) || has(self.crn) || (has(self.cosInstance) && has(self.cosBucket) && has(self.cosObject))",message="an existing image name or crn must be provided, or to create a new image the cos resources must be provided, with or without a name"
type ImageSpec struct {
	// name is the name of the desired VPC Custom Image.
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=63
	// +kubebuilder:validation:Pattern=`^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`
	// +optional
	Name *string `json:"name,omitempty"`

	// crn is the IBM Cloud CRN of the existing VPC Custom Image.
	// +optional
	CRN *string `json:"crn,omitempty"`

	// cosInstance is the name of the IBM Cloud COS Instance containing the source of the image, if necessary.
	// +optional
	COSInstance *string `json:"cosInstance,omitempty"`

	// cosBucket is the name of the IBM Cloud COS Bucket containing the source of the image, if necessary.
	// +optional
	COSBucket *string `json:"cosBucket,omitempty"`

	// cosBucketRegion is the COS region the bucket is in.
	// +optional
	COSBucketRegion *string `json:"cosBucketRegion,omitempty"`

	// cosObject is the name of a IBM Cloud COS Object used as the source of the image, if necessary.
	// +optional
	COSObject *string `json:"cosObject,omitempty"`

	// operatingSystem is the Custom Image's Operating System name.
	// +optional
	OperatingSystem *string `json:"operatingSystem,omitempty"`

	// resourceGroup is the Resource Group to create the Custom Image in.
	// +optional
	ResourceGroup *IBMCloudResourceReference `json:"resourceGroup,omitempty"`
}

// VPCNetworkSpec defines the desired state of the network resources for the cluster for extended VPC Infrastructure support.
type VPCNetworkSpec struct {
	// controlPlaneSubnets is a set of Subnet's which define the Control Plane subnets.
	// +optional
	ControlPlaneSubnets []Subnet `json:"controlPlaneSubnets,omitempty"`

	// loadBalancers is a set of VPC Load Balancer definitions to use for the cluster.
	// +optional
	LoadBalancers []VPCLoadBalancerSpec `json:"loadBalancers,omitempty"`

	// resourceGroup is the Resource Group containing all of the newtork resources.
	// This can be different than the Resource Group containing the remaining cluster resources.
	// +optional
	ResourceGroup *IBMCloudResourceReference `json:"resourceGroup,omitempty"`

	// securityGroups is a set of VPCSecurityGroup's which define the VPC Security Groups that manage traffic within and out of the VPC.
	// +optional
	SecurityGroups []VPCSecurityGroup `json:"securityGroups,omitempty"`

	// workerSubnets is a set of Subnet's which define the Worker subnets.
	// +optional
	WorkerSubnets []Subnet `json:"workerSubnets,omitempty"`

	// vpc defines the IBM Cloud VPC for extended VPC Infrastructure support.
	// +optional
	VPC *VPCResource `json:"vpc,omitempty"`
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

// IBMVPCClusterStatus defines the observed state of IBMVPCCluster.
type IBMVPCClusterStatus struct {
	// Important: Run "make" to regenerate code after modifying this file
	// dep: rely on Network instead.
	VPC VPC `json:"vpc,omitempty"`

	// image is the status of the VPC Custom Image.
	// +optional
	Image *ResourceStatus `json:"image,omitempty"`

	// network is the status of the VPC network resources for extended VPC Infrastructure support.
	// +optional
	Network *VPCNetworkStatus `json:"network,omitempty"`

	// Ready is true when the provider resource is ready.
	// +optional
	// +kubebuilder:default=false
	Ready bool `json:"ready"`

	// resourceGroup is the status of the cluster's Resource Group for extended VPC Infrastructure support.
	// +optional
	ResourceGroup *ResourceStatus `json:"resourceGroup,omitempty"`

	Subnet      Subnet      `json:"subnet,omitempty"`
	VPCEndpoint VPCEndpoint `json:"vpcEndpoint,omitempty"`

	// ControlPlaneLoadBalancerState is the status of the load balancer.
	// +optional
	ControlPlaneLoadBalancerState VPCLoadBalancerState `json:"controlPlaneLoadBalancerState,omitempty"`

	// Conditions defines current service state of the load balancer.
	// +optional
	Conditions capiv1beta1.Conditions `json:"conditions,omitempty"`
}

// VPCNetworkStatus provides details on the status of VPC network resources for extended VPC Infrastructure support.
type VPCNetworkStatus struct {
	// controlPlaneSubnets references the VPC Subnets for the cluster's Control Plane.
	// The map simplifies lookups.
	// +optional
	ControlPlaneSubnets map[string]*ResourceStatus `json:"controlPlaneSubnets,omitempty"`

	// loadBalancers references the VPC Load Balancer's for the cluster.
	// The map simplifies lookups.
	// +optional
	LoadBalancers map[string]*VPCLoadBalancerStatus `json:"loadBalancers,omitempty"`

	// publicGateways references the VPC Public Gateways for the cluster.
	// The map simplifies lookups.
	// +optional
	PublicGateways map[string]*ResourceStatus `json:"publicGateways,omitempty"`

	// resourceGroup references the Resource Group for Network resources for the cluster.
	// This can be the same or unique from the cluster's Resource Group.
	// +optional
	ResourceGroup *ResourceStatus `json:"resourceGroup,omitempty"`

	// securityGroups references the VPC Security Groups for the cluster.
	// The map simplifies lookups.
	// +optional
	SecurityGroups map[string]*ResourceStatus `json:"securityGroups,omitempty"`

	// workerSubnets references the VPC Subnets for the cluster's Data Plane.
	// The map simplifies lookups.
	// +optional
	WorkerSubnets map[string]*ResourceStatus `json:"workerSubnets,omitempty"`

	// vpc references the status of the IBM Cloud VPC as part of the extended VPC Infrastructure support.
	// +optional
	VPC *ResourceStatus `json:"vpc,omitempty"`
}

// VPC holds the VPC information.
type VPC struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=ibmvpcclusters,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this IBMVPCCluster belongs"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Cluster infrastructure is ready for IBM VPC instances"

// IBMVPCCluster is the Schema for the ibmvpcclusters API.
type IBMVPCCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IBMVPCClusterSpec   `json:"spec,omitempty"`
	Status IBMVPCClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// IBMVPCClusterList contains a list of IBMVPCCluster.
type IBMVPCClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IBMVPCCluster `json:"items"`
}

// GetConditions returns the observations of the operational state of the IBMVPCCluster resource.
func (r *IBMVPCCluster) GetConditions() capiv1beta1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the IBMVPCCluster to the predescribed clusterv1.Conditions.
func (r *IBMVPCCluster) SetConditions(conditions capiv1beta1.Conditions) {
	r.Status.Conditions = conditions
}

func init() {
	objectTypes = append(objectTypes, &IBMVPCCluster{}, &IBMVPCClusterList{})
}
