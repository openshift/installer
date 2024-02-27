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

package filter

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

const (
	filterNameTagKey        = "tag-key"
	filterNameVpcID         = "vpc-id"
	filterNameState         = "state"
	filterNameVpcAttachment = "attachment.vpc-id"
	filterAvailabilityZone  = "availability-zone"
	filterNameIPAMPoolID    = "ipam-pool-id"
)

// EC2 exposes the ec2 sdk related filters.
var EC2 = new(ec2Filters)

type ec2Filters struct{}

// Cluster returns a filter based on the cluster name.
func (ec2Filters) Cluster(clusterName string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(filterNameTagKey),
		Values: aws.StringSlice([]string{infrav1.ClusterTagKey(clusterName)}),
	}
}

// Name returns a filter based on the resource name.
func (ec2Filters) Name(name string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String("tag:Name"),
		Values: aws.StringSlice([]string{name}),
	}
}

// ClusterOwned returns a filter using the Cluster API per-cluster tag where
// the resource is owned.
func (ec2Filters) ClusterOwned(clusterName string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(fmt.Sprintf("tag:%s", infrav1.ClusterTagKey(clusterName))),
		Values: aws.StringSlice([]string{string(infrav1.ResourceLifecycleOwned)}),
	}
}

// ClusterShared returns a filter using the Cluster API per-cluster tag where
// the resource is shared.
func (ec2Filters) ClusterShared(clusterName string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(fmt.Sprintf("tag:%s", infrav1.ClusterTagKey(clusterName))),
		Values: aws.StringSlice([]string{string(infrav1.ResourceLifecycleShared)}),
	}
}

// ProviderRole returns a filter using cluster-api-provider-aws role tag.
func (ec2Filters) ProviderRole(role string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(fmt.Sprintf("tag:%s", infrav1.NameAWSClusterAPIRole)),
		Values: aws.StringSlice([]string{role}),
	}
}

// ProviderOwned returns a filter using the cloud provider tag where the resource is owned.
func (ec2Filters) ProviderOwned(clusterName string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(fmt.Sprintf("tag:%s", infrav1.ClusterAWSCloudProviderTagKey(clusterName))),
		Values: aws.StringSlice([]string{string(infrav1.ResourceLifecycleOwned)}),
	}
}

// IPAM returns a filter based on the id of the IPAM Pool.
func (ec2Filters) IPAM(ipamPoolID string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(filterNameIPAMPoolID),
		Values: aws.StringSlice([]string{ipamPoolID}),
	}
}

// VPC returns a filter based on the id of the VPC.
func (ec2Filters) VPC(vpcID string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(filterNameVpcID),
		Values: aws.StringSlice([]string{vpcID}),
	}
}

// VPCAttachment returns a filter based on the vpc id attached to the resource.
func (ec2Filters) VPCAttachment(vpcID string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(filterNameVpcAttachment),
		Values: aws.StringSlice([]string{vpcID}),
	}
}

// Available returns a filter based on the state being available.
func (ec2Filters) Available() *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(filterNameState),
		Values: aws.StringSlice([]string{"available"}),
	}
}

// NATGatewayStates returns a filter based on the list of states passed in.
func (ec2Filters) NATGatewayStates(states ...string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String("state"),
		Values: aws.StringSlice(states),
	}
}

// InstanceStates returns a filter based on the list of states passed in.
func (ec2Filters) InstanceStates(states ...string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String("instance-state-name"),
		Values: aws.StringSlice(states),
	}
}

// VPCStates returns a filter based on the list of states passed in.
func (ec2Filters) VPCStates(states ...string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String("state"),
		Values: aws.StringSlice(states),
	}
}

// SubnetStates returns a filter based on the list of states passed in.
func (ec2Filters) SubnetStates(states ...string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String("state"),
		Values: aws.StringSlice(states),
	}
}

func (ec2Filters) AvailabilityZone(zone string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(filterAvailabilityZone),
		Values: aws.StringSlice([]string{zone}),
	}
}

func (ec2Filters) IgnoreLocalZones() *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String("opt-in-status"),
		Values: aws.StringSlice([]string{"opt-in-not-required"}),
	}
}

func (ec2Filters) SecurityGroupName(name string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String("group-name"),
		Values: aws.StringSlice([]string{name}),
	}
}
