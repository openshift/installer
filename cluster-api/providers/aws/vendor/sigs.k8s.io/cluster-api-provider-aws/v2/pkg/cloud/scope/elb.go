/*
Copyright 2018 The Kubernetes Authors.

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
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// ELBScope is a scope for use with the ELB reconciling service.
type ELBScope interface {
	cloud.ClusterScoper

	// Network returns the cluster network object.
	Network() *infrav1.NetworkStatus

	// Subnets returns the cluster subnets.
	Subnets() infrav1.Subnets

	// SecurityGroups returns the cluster security groups as a map, it creates the map if empty.
	SecurityGroups() map[infrav1.SecurityGroupRole]infrav1.SecurityGroup

	// VPC returns the cluster VPC.
	VPC() *infrav1.VPCSpec

	// ControlPlaneLoadBalancer returns the AWSLoadBalancerSpec
	ControlPlaneLoadBalancer() *infrav1.AWSLoadBalancerSpec

	// ControlPlaneLoadBalancerScheme returns the Classic ELB scheme (public or internal facing)
	ControlPlaneLoadBalancerScheme() infrav1.ELBScheme

	// ControlPlaneLoadBalancerName returns the Classic ELB name
	ControlPlaneLoadBalancerName() *string

	// ControlPlaneEndpoint returns AWSCluster control plane endpoint
	ControlPlaneEndpoint() clusterv1.APIEndpoint
}
