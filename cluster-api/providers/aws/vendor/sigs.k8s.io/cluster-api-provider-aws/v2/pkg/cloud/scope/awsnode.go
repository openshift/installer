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

package scope

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
)

// AWSNodeScope is the interface for the scope to be used with the awsnode reconciling service.
type AWSNodeScope interface {
	cloud.ClusterScoper

	// RemoteClient returns the Kubernetes client for connecting to the workload cluster.
	RemoteClient() (client.Client, error)
	// Subnets returns the cluster subnets.
	Subnets() infrav1.Subnets
	// SecondaryCidrBlock returns the optional secondary CIDR block to use for pod IPs
	SecondaryCidrBlock() *string
	// SecurityGroups returns the control plane security groups as a map, it creates the map if empty.
	SecurityGroups() map[infrav1.SecurityGroupRole]infrav1.SecurityGroup
	// DisableVPCCNI returns whether the AWS VPC CNI should be disabled
	DisableVPCCNI() bool
	// VpcCni specifies configuration related to the VPC CNI.
	VpcCni() ekscontrolplanev1.VpcCni
	// VPC returns the given VPC configuration.
	VPC() *infrav1.VPCSpec
}
