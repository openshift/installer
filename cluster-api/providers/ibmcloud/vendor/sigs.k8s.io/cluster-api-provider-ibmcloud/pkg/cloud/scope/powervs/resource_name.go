/*
Copyright 2025 The Kubernetes Authors.

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

package powervs

import (
	"fmt"
	"strings"
)

// ResourceType identifies the kind of IBM Cloud resource whose default name is being generated.
type ResourceType string

const (
	// ResourceTypeWorkspace is a PowerVS service instance (workspace).
	ResourceTypeWorkspace ResourceType = "workspace"
	// ResourceTypeTransitGateway is an IBM Transit Gateway.
	ResourceTypeTransitGateway ResourceType = "transitgateway"
	// ResourceTypeVPC is a Virtual Private Cloud.
	ResourceTypeVPC ResourceType = "vpc"
	// ResourceTypeSubnet is a VPC subnet.
	// Use a VPC zone string (e.g. "us-south-1") or strconv.Itoa(i) as the qualifier.
	ResourceTypeSubnet ResourceType = "subnet"
	// ResourceTypeLBPublic is a public-facing VPC load balancer.
	// Use a non-empty qualifier (e.g. strconv.Itoa(i)) for additional LBs of the same type.
	ResourceTypeLBPublic ResourceType = "loadbalancer-public"
	// ResourceTypeLBPrivate is an internal VPC load balancer.
	// Use a non-empty qualifier (e.g. strconv.Itoa(i)) for additional LBs of the same type.
	ResourceTypeLBPrivate ResourceType = "loadbalancer-private"
	// ResourceTypeSG is a VPC security group.
	ResourceTypeSG ResourceType = "sg"
	// ResourceTypeDHCP is a PowerVS DHCP server.
	// Note: IBM Cloud derives the associated network name via dhcpNetworkName(); this type
	// controls only the DHCP server resource name itself.
	ResourceTypeDHCP ResourceType = "dhcp"
	// ResourceTypeCOS is a Cloud Object Storage service instance.
	ResourceTypeCOS ResourceType = "cos"
	// ResourceTypeCOSBucket is a COS bucket inside a COS instance.
	ResourceTypeCOSBucket ResourceType = "cos-bucket"
	// ResourceTypeCOSHMACKey is a COS HMAC service credential key.
	ResourceTypeCOSHMACKey ResourceType = "cos-hmac"
)

// resourceNameMaxLen is the maximum length allowed for IBM Cloud resource names.
// Most IBM Cloud APIs (VPC, TG, PowerVS) enforce a 63-character limit.
const resourceNameMaxLen = 63

// ResourceName returns the deterministic default name for an IBM Cloud resource owned by the
// named cluster. qualifier is appended when a resource type needs zone or index disambiguation
// (e.g. a VPC zone string for subnets, strconv.Itoa(i) for indexed load balancers).
// Pass an empty string for resource types that do not require a qualifier.
//
// The returned name is always truncated to resourceNameMaxLen characters. Truncation removes
// characters from the cluster-name prefix only, so the resource-type suffix is always preserved
// and names remain distinguishable from one another.
func ResourceName(clusterName string, rt ResourceType, qualifier string) string {
	suffix := string(rt)
	if qualifier != "" {
		suffix = fmt.Sprintf("%s-%s", rt, qualifier)
	}

	// +1 accounts for the "-" separator between base and suffix.
	maxBase := resourceNameMaxLen - len(suffix) - 1
	base := clusterName
	if len(base) > maxBase {
		// Truncate to the budget, then strip any trailing hyphens so the join
		// never produces a double-hyphen (e.g. "my-cluster--subnet-1") and the
		// base does not end with a separator character.
		base = strings.TrimRight(base[:maxBase], "-")
	}

	return fmt.Sprintf("%s-%s", base, suffix)
}

// tgVPCConnectionName returns the default name for the VPC-side TG connection.
func tgVPCConnectionName(tgName string) string { return fmt.Sprintf("%s-vpc-con", tgName) }

// tgPowerVSConnectionName returns the default name for the PowerVS-side TG connection.
func tgPowerVSConnectionName(tgName string) string { return fmt.Sprintf("%s-pvs-con", tgName) }

// dhcpNetworkName returns the network name IBM Cloud assigns to a DHCP server.
func dhcpNetworkName(dhcpServerName string) string {
	return fmt.Sprintf("DHCPSERVER%s_Private", dhcpServerName)
}
