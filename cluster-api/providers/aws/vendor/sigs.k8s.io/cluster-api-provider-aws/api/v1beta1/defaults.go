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
package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clusterv1 "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
)

// SetDefaults_Bastion is used by defaulter-gen.
func SetDefaults_Bastion(obj *Bastion) { //nolint:golint,stylecheck
	// Default to allow open access to the bastion host if no CIDR Blocks have been set
	if len(obj.AllowedCIDRBlocks) == 0 && !obj.DisableIngressRules {
		obj.AllowedCIDRBlocks = []string{"0.0.0.0/0"}
	}
}

// SetDefaults_NetworkSpec is used by defaulter-gen.
func SetDefaults_NetworkSpec(obj *NetworkSpec) { //nolint:golint,stylecheck
	// Default to Calico ingress rules if no rules have been set
	if obj.CNI == nil {
		obj.CNI = &CNISpec{
			CNIIngressRules: CNIIngressRules{
				{
					Description: "bgp (calico)",
					Protocol:    SecurityGroupProtocolTCP,
					FromPort:    179,
					ToPort:      179,
				},
				{
					Description: "IP-in-IP (calico)",
					Protocol:    SecurityGroupProtocolIPinIP,
					FromPort:    -1,
					ToPort:      65535,
				},
			},
		}
	}
}

// SetDefaults_AWSClusterSpec is used by defaulter-gen.
func SetDefaults_AWSClusterSpec(s *AWSClusterSpec) { //nolint:golint,stylecheck
	if s.IdentityRef == nil {
		s.IdentityRef = &AWSIdentityReference{
			Kind: ControllerIdentityKind,
			Name: AWSClusterControllerIdentityName,
		}
	}

	// 	If ELB scheme is set to Internet-facing due to an API bug in versions > v0.6.6 and v0.7.0, default it to internet-facing.
	if s.ControlPlaneLoadBalancer == nil {
		s.ControlPlaneLoadBalancer = &AWSLoadBalancerSpec{Scheme: &ClassicELBSchemeInternetFacing}
	} else if s.ControlPlaneLoadBalancer.Scheme != nil && s.ControlPlaneLoadBalancer.Scheme.String() == ClassicELBSchemeIncorrectInternetFacing.String() {
		s.ControlPlaneLoadBalancer.Scheme = &ClassicELBSchemeInternetFacing
	}
}

// SetDefaults_Labels is used to default cluster scope resources for clusterctl move.
func SetDefaults_Labels(obj *metav1.ObjectMeta) { //nolint:golint,stylecheck
	// Defaults to set label if no labels have been set
	if obj.Labels == nil {
		obj.Labels = map[string]string{
			clusterv1.ClusterctlMoveHierarchyLabelName: ""}
	}
}
